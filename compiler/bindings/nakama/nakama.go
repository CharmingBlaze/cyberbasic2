// Package nakama provides Nakama cloud backend bindings for multiplayer: accounts, matchmaking, realtime matches.
package nakama

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/vm"
	"github.com/ascii8/nakama-go"
)

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func toInt(v interface{}) int {
	switch x := v.(type) {
	case int:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}

func toBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch x := v.(type) {
	case int:
		return x != 0
	case int32:
		return x != 0
	case int64:
		return x != 0
	case float64:
		return x != 0
	case bool:
		return x
	default:
		return false
	}
}

type nakamaEvent struct {
	typ      string
	matchId  string
	opCode   int64
	data     string
	sender   string
	presences string // JSON array for join/leave
	token    string
}

var (
	nakamaVM     *vm.VM
	nakamaClient *nakama.Client
	nakamaConn   *nakama.Conn
	nakamaCtx    context.Context
	nakamaCancel context.CancelFunc
	eventQueue   []nakamaEvent
	eventMu      sync.Mutex
)

func pushEvent(ev nakamaEvent) {
	eventMu.Lock()
	eventQueue = append(eventQueue, ev)
	eventMu.Unlock()
}

func drainEvents() []nakamaEvent {
	eventMu.Lock()
	out := eventQueue
	eventQueue = nil
	eventMu.Unlock()
	return out
}

// RegisterNakama registers Nakama multiplayer functions with the VM.
func RegisterNakama(v *vm.VM) {
	nakamaVM = v

	// NakamaConnect(host, port, serverKey [, useSSL])
	v.RegisterForeign("NakamaConnect", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("NakamaConnect(host, port, serverKey [, useSSL]) requires at least 3 arguments")
		}
		host := toString(args[0])
		port := toInt(args[1])
		serverKey := toString(args[2])
		useSSL := false
		if len(args) >= 4 {
			useSSL = toBool(args[3])
		}
		scheme := "http"
		if useSSL {
			scheme = "https"
		}
		url := fmt.Sprintf("%s://%s:%d", scheme, host, port)
		nakamaClient = nakama.New(nakama.WithURL(url), nakama.WithServerKey(serverKey))
		return 1, nil
	})

	// NakamaAuthenticateDevice(deviceId [, create, username])
	v.RegisterForeign("NakamaAuthenticateDevice", func(args []interface{}) (interface{}, error) {
		if nakamaClient == nil {
			return nil, fmt.Errorf("call NakamaConnect first")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("NakamaAuthenticateDevice(deviceId [, create, username]) requires at least 1 argument")
		}
		deviceId := toString(args[0])
		create := true
		username := ""
		if len(args) >= 2 {
			create = toBool(args[1])
		}
		if len(args) >= 3 {
			username = toString(args[2])
		}
		ctx := context.Background()
		if err := nakamaClient.AuthenticateDevice(ctx, deviceId, create, username); err != nil {
			return nil, err
		}
		return 1, nil
	})

	// NakamaAuthenticateCustom(customId [, create, username])
	v.RegisterForeign("NakamaAuthenticateCustom", func(args []interface{}) (interface{}, error) {
		if nakamaClient == nil {
			return nil, fmt.Errorf("call NakamaConnect first")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("NakamaAuthenticateCustom(customId [, create, username]) requires at least 1 argument")
		}
		customId := toString(args[0])
		create := true
		username := ""
		if len(args) >= 2 {
			create = toBool(args[1])
		}
		if len(args) >= 3 {
			username = toString(args[2])
		}
		ctx := context.Background()
		if err := nakamaClient.AuthenticateCustom(ctx, customId, create, username); err != nil {
			return nil, err
		}
		return 1, nil
	})

	// NakamaAuthenticateEmail(email, password [, create, username])
	v.RegisterForeign("NakamaAuthenticateEmail", func(args []interface{}) (interface{}, error) {
		if nakamaClient == nil {
			return nil, fmt.Errorf("call NakamaConnect first")
		}
		if len(args) < 2 {
			return nil, fmt.Errorf("NakamaAuthenticateEmail(email, password [, create, username]) requires at least 2 arguments")
		}
		email := toString(args[0])
		password := toString(args[1])
		create := true
		username := ""
		if len(args) >= 3 {
			create = toBool(args[2])
		}
		if len(args) >= 4 {
			username = toString(args[3])
		}
		ctx := context.Background()
		if err := nakamaClient.AuthenticateEmail(ctx, email, password, create, username); err != nil {
			return nil, err
		}
		return 1, nil
	})

	// NakamaCreateSocket() - creates realtime socket (Conn)
	v.RegisterForeign("NakamaCreateSocket", func(args []interface{}) (interface{}, error) {
		if nakamaClient == nil {
			return nil, fmt.Errorf("call NakamaConnect and authenticate first")
		}
		if nakamaConn != nil {
			nakamaCancel()
			nakamaConn = nil
		}
		nakamaCtx, nakamaCancel = context.WithCancel(context.Background())
		conn, err := nakamaClient.NewConn(nakamaCtx,
			nakama.WithConnPersist(true),
			nakama.WithConnHandler(&nakamaHandler{vm: nakamaVM}))
		if err != nil {
			return nil, err
		}
		nakamaConn = conn
		return 1, nil
	})

	// NakamaSocketConnect() - opens the socket (already created by NakamaCreateSocket)
	v.RegisterForeign("NakamaSocketConnect", func(args []interface{}) (interface{}, error) {
		if nakamaConn == nil {
			return nil, fmt.Errorf("call NakamaCreateSocket first")
		}
		if err := nakamaConn.Open(nakamaCtx); err != nil {
			return nil, err
		}
		return 1, nil
	})

	// NakamaCreateMatch([name])
	v.RegisterForeign("NakamaCreateMatch", func(args []interface{}) (interface{}, error) {
		if nakamaConn == nil {
			return nil, fmt.Errorf("call NakamaCreateSocket and NakamaSocketConnect first")
		}
		name := ""
		if len(args) >= 1 {
			name = toString(args[0])
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		match, err := nakamaConn.MatchCreate(ctx, name)
		if err != nil {
			return nil, err
		}
		if match == nil {
			return nil, nil
		}
		return match.MatchId, nil
	})

	// NakamaJoinMatch(matchId [, token])
	v.RegisterForeign("NakamaJoinMatch", func(args []interface{}) (interface{}, error) {
		if nakamaConn == nil {
			return nil, fmt.Errorf("call NakamaCreateSocket and NakamaSocketConnect first")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("NakamaJoinMatch(matchId [, token]) requires at least 1 argument")
		}
		matchId := toString(args[0])
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var match *nakama.MatchMsg
		var err error
		if len(args) >= 2 && toString(args[1]) != "" {
			token := toString(args[1])
			match, err = nakamaConn.MatchJoinToken(ctx, token, nil)
		} else {
			match, err = nakamaConn.MatchJoin(ctx, matchId, nil)
		}
		if err != nil {
			return nil, err
		}
		if match == nil {
			return nil, nil
		}
		return match.MatchId, nil
	})

	// NakamaLeaveMatch(matchId)
	v.RegisterForeign("NakamaLeaveMatch", func(args []interface{}) (interface{}, error) {
		if nakamaConn == nil {
			return nil, fmt.Errorf("call NakamaCreateSocket first")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("NakamaLeaveMatch(matchId) requires 1 argument")
		}
		matchId := toString(args[0])
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := nakamaConn.MatchLeave(ctx, matchId); err != nil {
			return nil, err
		}
		return nil, nil
	})

	// NakamaSendMatchState(matchId, opCode, data [, reliable])
	v.RegisterForeign("NakamaSendMatchState", func(args []interface{}) (interface{}, error) {
		if nakamaConn == nil {
			return nil, fmt.Errorf("call NakamaCreateSocket first")
		}
		if len(args) < 3 {
			return nil, fmt.Errorf("NakamaSendMatchState(matchId, opCode, data [, reliable]) requires at least 3 arguments")
		}
		matchId := toString(args[0])
		opCode := int64(toInt(args[1]))
		data := []byte(toString(args[2]))
		reliable := true
		if len(args) >= 4 {
			reliable = toBool(args[3])
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := nakamaConn.MatchDataSend(ctx, matchId, opCode, data, reliable); err != nil {
			return nil, err
		}
		return 1, nil
	})

	// NakamaAddMatchmaker([minPlayers, maxPlayers, query])
	v.RegisterForeign("NakamaAddMatchmaker", func(args []interface{}) (interface{}, error) {
		if nakamaConn == nil {
			return nil, fmt.Errorf("call NakamaCreateSocket first")
		}
		minPlayers := 2
		maxPlayers := 8
		query := ""
		if len(args) >= 1 {
			minPlayers = toInt(args[0])
		}
		if len(args) >= 2 {
			maxPlayers = toInt(args[1])
		}
		if len(args) >= 3 {
			query = toString(args[2])
		}
		if minPlayers <= 0 {
			minPlayers = 2
		}
		if maxPlayers < minPlayers {
			maxPlayers = minPlayers
		}
		msg := nakama.MatchmakerAdd(query, minPlayers, maxPlayers)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ticket, err := nakamaConn.MatchmakerAdd(ctx, msg)
		if err != nil {
			return nil, err
		}
		if ticket == nil {
			return nil, nil
		}
		return ticket.Ticket, nil
	})

	// NakamaRemoveMatchmaker(ticket)
	v.RegisterForeign("NakamaRemoveMatchmaker", func(args []interface{}) (interface{}, error) {
		if nakamaConn == nil {
			return nil, fmt.Errorf("call NakamaCreateSocket first")
		}
		if len(args) < 1 {
			return nil, fmt.Errorf("NakamaRemoveMatchmaker(ticket) requires 1 argument")
		}
		ticket := toString(args[0])
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := nakamaConn.MatchmakerRemove(ctx, ticket); err != nil {
			return nil, err
		}
		return nil, nil
	})

	// NakamaGetAccount() - returns JSON string of account or empty on error
	v.RegisterForeign("NakamaGetAccount", func(args []interface{}) (interface{}, error) {
		if nakamaClient == nil {
			return nil, fmt.Errorf("call NakamaConnect and authenticate first")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		account, err := nakamaClient.Account(ctx)
		if err != nil {
			return "", nil
		}
		if account == nil {
			return "", nil
		}
		// Return simple JSON-like string: {"userId":"...","username":"...","devices":N}
		userId, username := "", ""
		if account.User != nil {
			userId = account.User.Id
			username = account.User.Username
		}
		return fmt.Sprintf(`{"userId":"%s","username":"%s","devices":%d}`,
			userId, username, len(account.Devices)), nil
	})

	// NakamaRPC(id, input) - HTTP RPC, returns response as string (JSON)
	v.RegisterForeign("NakamaRPC", func(args []interface{}) (interface{}, error) {
		if nakamaClient == nil {
			return nil, fmt.Errorf("call NakamaConnect and authenticate first")
		}
		if len(args) < 2 {
			return nil, fmt.Errorf("NakamaRPC(id, input) requires 2 arguments")
		}
		id := toString(args[0])
		input := toString(args[1])
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var result interface{}
		if err := nakamaClient.Rpc(ctx, id, input, &result); err != nil {
			return nil, err
		}
		if s, ok := result.(string); ok {
			return s, nil
		}
		b, _ := json.Marshal(result)
		return string(b), nil
	})

	// NakamaProcessEvents() - drain event queue and invoke BASIC callbacks
	v.RegisterForeign("NakamaProcessEvents", func(args []interface{}) (interface{}, error) {
		events := drainEvents()
		if nakamaVM == nil || nakamaVM.Chunk() == nil {
			return nil, nil
		}
		for _, ev := range events {
			switch ev.typ {
			case "match_data":
				if _, ok := nakamaVM.Chunk().GetFunction("onnakamamatchdata"); ok {
					_ = nakamaVM.InvokeSub("onnakamamatchdata", []interface{}{ev.matchId, ev.opCode, ev.data, ev.sender})
				}
			case "match_join":
				if _, ok := nakamaVM.Chunk().GetFunction("onnakamamatchjoin"); ok {
					_ = nakamaVM.InvokeSub("onnakamamatchjoin", []interface{}{ev.matchId, ev.presences})
				}
			case "match_leave":
				if _, ok := nakamaVM.Chunk().GetFunction("onnakamamatchleave"); ok {
					_ = nakamaVM.InvokeSub("onnakamamatchleave", []interface{}{ev.matchId, ev.presences})
				}
			case "matchmaker_matched":
				if _, ok := nakamaVM.Chunk().GetFunction("onnakamamatchmakermatched"); ok {
					_ = nakamaVM.InvokeSub("onnakamamatchmakermatched", []interface{}{ev.matchId, ev.token})
				}
			}
		}
		return nil, nil
	})

	v.SetGlobal("nakama", modfacade.New(v, nakamaV2))
}

// nakamaHandler implements ConnHandler for pushing events to the queue.
type nakamaHandler struct {
	vm *vm.VM
}

func (h *nakamaHandler) MatchDataHandler(ctx context.Context, msg *nakama.MatchDataMsg) {
	sender := ""
	if msg.Presence != nil {
		sender = msg.Presence.UserId
	}
	pushEvent(nakamaEvent{
		typ:     "match_data",
		matchId: msg.MatchId,
		opCode:  msg.OpCode,
		data:    string(msg.Data),
		sender:  sender,
	})
}

func (h *nakamaHandler) MatchPresenceEventHandler(ctx context.Context, msg *nakama.MatchPresenceEventMsg) {
	// MatchPresenceEvent has Joins and Leaves
	presences := make([]string, 0)
	for _, p := range msg.Joins {
		if p != nil {
			presences = append(presences, p.UserId)
		}
	}
	for _, p := range msg.Leaves {
		if p != nil {
			presences = append(presences, "leave:"+p.UserId)
		}
	}
	presStr := "[" + strings.Join(presences, ",") + "]"
	if len(msg.Joins) > 0 {
		pushEvent(nakamaEvent{typ: "match_join", matchId: msg.MatchId, presences: presStr})
	}
	if len(msg.Leaves) > 0 {
		pushEvent(nakamaEvent{typ: "match_leave", matchId: msg.MatchId, presences: presStr})
	}
}

func (h *nakamaHandler) MatchmakerMatchedHandler(ctx context.Context, msg *nakama.MatchmakerMatchedMsg) {
	matchId := msg.GetMatchId()
	token := msg.GetToken()
	pushEvent(nakamaEvent{typ: "matchmaker_matched", matchId: matchId, token: token})
}
