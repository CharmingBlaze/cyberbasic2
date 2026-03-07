// Package net provides TCP client/server bindings for multiplayer: Connect, Send, Receive, Disconnect, Host, Accept, CloseServer.
package net

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"cyberbasic/compiler/vm"
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
		return x
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

func toFloat(v interface{}) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case int32:
		return float64(x)
	case int64:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}

const maxMessageSize = 256 * 1024 // 256KB max per message (security and resource limit)
const maxSendNumbers = 16         // max numbers in SendNumbers / SendToRoomNumbers

// netEvent is one item in the event queue for ProcessNetworkEvents (connect, disconnect, message).
type netEvent struct {
	typ     string // "connect", "disconnect", "message"
	id      string
	payload string
}

type deadlineListener interface {
	SetDeadline(time.Time) error
}

type serverState struct {
	listener       net.Listener
	deadlineSource deadlineListener
}

var (
	netVM             *vm.VM
	eventQueue        []netEvent
	eventMu           sync.Mutex
	rpcHandlers       = make(map[string]string) // RPC name (lowercase) -> Sub name for InvokeSub
	rpcMu             sync.Mutex
	pingSentAt        = make(map[string]time.Time)
	lastRTTMs         = make(map[string]float64)
	pingMu            sync.Mutex
	remoteEntities    = make(map[string]map[string]interface{}) // entityId -> {x, y, z}
	remoteEntitiesMu  sync.Mutex
	connMessages      = make(map[string][]string) // per-connection message queue (filled by reader goroutine)
	connMessagesMu    sync.Mutex
	conns             = make(map[string]net.Conn)
	readers           = make(map[string]*bufio.Reader)
	servers           = make(map[string]*serverState)
	rooms             = make(map[string]map[string]bool) // roomId -> set of connectionIds
	netMu             sync.Mutex
	connCounter       int
	servCounter       int
	receivedNumbers   = make(map[string][]float64)
	lastNumbersConnID string
	receivedNumbersMu sync.Mutex
)

func pushEvent(typ, id, payload string) {
	eventMu.Lock()
	eventQueue = append(eventQueue, netEvent{typ: typ, id: id, payload: payload})
	eventMu.Unlock()
}

func drainEvents() []netEvent {
	eventMu.Lock()
	out := eventQueue
	eventQueue = nil
	eventMu.Unlock()
	return out
}

func cleanupConnection(cid string, conn net.Conn, sendDisconnectEvent bool) {
	removed := false
	netMu.Lock()
	if existing, ok := conns[cid]; ok && (conn == nil || existing == conn) {
		delete(conns, cid)
		delete(readers, cid)
		for roomID, set := range rooms {
			delete(set, cid)
			if len(set) == 0 {
				delete(rooms, roomID)
			}
		}
		removed = true
	}
	netMu.Unlock()
	if !removed {
		if conn != nil {
			_ = conn.Close()
		}
		return
	}
	connMessagesMu.Lock()
	delete(connMessages, cid)
	connMessagesMu.Unlock()
	receivedNumbersMu.Lock()
	delete(receivedNumbers, cid)
	if lastNumbersConnID == cid {
		lastNumbersConnID = ""
	}
	receivedNumbersMu.Unlock()
	pingMu.Lock()
	delete(pingSentAt, cid)
	delete(lastRTTMs, cid)
	pingMu.Unlock()
	if sendDisconnectEvent {
		pushEvent("disconnect", cid, "")
	}
	if conn != nil {
		_ = conn.Close()
	}
}

func addServer(listener net.Listener, deadlineSource deadlineListener) string {
	netMu.Lock()
	defer netMu.Unlock()
	servCounter++
	id := fmt.Sprintf("server_%d", servCounter)
	servers[id] = &serverState{listener: listener, deadlineSource: deadlineSource}
	return id
}

func acceptServerConnection(state *serverState, timeout time.Duration) (net.Conn, error) {
	if state == nil || state.listener == nil {
		return nil, fmt.Errorf("server not available")
	}
	if timeout > 0 && state.deadlineSource != nil {
		if err := state.deadlineSource.SetDeadline(time.Now().Add(timeout)); err != nil {
			return nil, err
		}
		defer state.deadlineSource.SetDeadline(time.Time{})
	}
	conn, err := state.listener.Accept()
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, nil
		}
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, nil
		}
		if errors.Is(err, net.ErrClosed) {
			return nil, nil
		}
		return nil, err
	}
	return conn, nil
}

// startReader runs in a goroutine; reads lines from conn, appends to connMessages[id], pushes "message" events; on error pushes "disconnect" and removes conn.
func startReader(cid string, conn net.Conn) {
	rd := bufio.NewReader(conn)
	netMu.Lock()
	readers[cid] = rd
	netMu.Unlock()
	for {
		_ = conn.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
		line, err := rd.ReadString('\n')
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			}
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				cleanupConnection(cid, conn, true)
				return
			}
			cleanupConnection(cid, conn, true)
			return
		}
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		if len(line) > maxMessageSize {
			continue
		}
		if line == "" {
			continue
		}
		if line == "P" {
			_, _ = fmt.Fprintln(conn, "O")
			continue
		}
		if line == "O" {
			pingMu.Lock()
			if t, ok := pingSentAt[cid]; ok {
				lastRTTMs[cid] = float64(time.Since(t).Milliseconds())
			}
			pingMu.Unlock()
			continue
		}
		connMessagesMu.Lock()
		connMessages[cid] = append(connMessages[cid], line)
		connMessagesMu.Unlock()
		pushEvent("message", cid, "")
	}
}

// RegisterNet registers TCP multiplayer functions with the VM.
func RegisterNet(v *vm.VM) {
	netVM = v
	// --- Client ---
	v.RegisterForeign("Connect", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Connect(host, port) requires 2 arguments")
		}
		host := toString(args[0])
		port := toInt(args[1])
		addr := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
		if err != nil {
			return nil, nil // return null on failure so BASIC can IsNull check
		}
		netMu.Lock()
		connCounter++
		id := fmt.Sprintf("conn_%d", connCounter)
		conns[id] = conn
		netMu.Unlock()
		go startReader(id, conn)
		return id, nil
	})
	v.RegisterForeign("ConnectTLS", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("ConnectTLS(host, port) requires 2 arguments")
		}
		host := toString(args[0])
		port := toInt(args[1])
		addr := fmt.Sprintf("%s:%d", host, port)
		config := &tls.Config{MinVersion: tls.VersionTLS12, ServerName: host}
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", addr, config)
		if err != nil {
			return nil, nil
		}
		netMu.Lock()
		connCounter++
		id := fmt.Sprintf("conn_%d", connCounter)
		conns[id] = conn
		netMu.Unlock()
		go startReader(id, conn)
		return id, nil
	})
	v.RegisterForeign("ConnectToParent", func(args []interface{}) (interface{}, error) {
		addr := os.Getenv("CYBERBASIC_PARENT")
		if addr == "" {
			return nil, nil
		}
		parts := strings.SplitN(addr, ":", 2)
		if len(parts) != 2 {
			return nil, nil
		}
		host := parts[0]
		port, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || port <= 0 {
			return nil, nil
		}
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
		if err != nil {
			return nil, nil
		}
		netMu.Lock()
		connCounter++
		id := fmt.Sprintf("conn_%d", connCounter)
		conns[id] = conn
		netMu.Unlock()
		go startReader(id, conn)
		return id, nil
	})
	// writeLine sends one line (no embedded newlines); enforces maxMessageSize
	writeLine := func(conn net.Conn, text string) error {
		if len(text) > maxMessageSize || strings.Contains(text, "\n") || strings.Contains(text, "\r") {
			return fmt.Errorf("message too long or contains newline")
		}
		_, err := fmt.Fprintln(conn, text)
		return err
	}
	v.RegisterForeign("Send", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("Send(connectionId, text) requires 2 arguments")
		}
		id := toString(args[0])
		text := toString(args[1])
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		err := writeLine(conn, text)
		return err == nil, err
	})
	v.RegisterForeign("SendJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendJSON(connectionId, jsonText) requires 2 arguments")
		}
		id := toString(args[0])
		text := toString(args[1])
		if len(text) > maxMessageSize {
			return 0, nil
		}
		if !json.Valid([]byte(text)) {
			return 0, nil
		}
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		err := writeLine(conn, text)
		if err != nil {
			return 0, nil
		}
		return 1, nil
	})
	v.RegisterForeign("SendTable", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendTable(connectionId, data) requires 2 arguments")
		}
		id := toString(args[0])
		m, ok := args[1].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("SendTable: data must be a dictionary (CreateDict / table)")
		}
		text, err := json.Marshal(m)
		if err != nil {
			return 0, nil
		}
		if len(text) > maxMessageSize {
			return 0, nil
		}
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		if writeLine(conn, string(text)) != nil {
			return 0, nil
		}
		return 1, nil
	})
	v.RegisterForeign("SendInt", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendInt(connectionId, value) requires 2 arguments")
		}
		id := toString(args[0])
		val := toInt(args[1])
		text := fmt.Sprintf("i %d", val)
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		if writeLine(conn, text) != nil {
			return 0, nil
		}
		return 1, nil
	})
	v.RegisterForeign("SendFloat", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendFloat(connectionId, value) requires 2 arguments")
		}
		id := toString(args[0])
		val := toFloat(args[1])
		text := fmt.Sprintf("f %g", val)
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		if writeLine(conn, text) != nil {
			return 0, nil
		}
		return 1, nil
	})
	v.RegisterForeign("SendNumbers", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendNumbers(connectionId, n1, n2, ...) requires at least 2 arguments")
		}
		id := toString(args[0])
		n := len(args) - 1
		if n > maxSendNumbers {
			n = maxSendNumbers
		}
		parts := make([]string, 0, n+1)
		parts = append(parts, "n")
		for i := 1; i <= n; i++ {
			parts = append(parts, fmt.Sprintf("%g", toFloat(args[i])))
		}
		text := strings.Join(parts, " ")
		if len(text) > maxMessageSize {
			return 0, nil
		}
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		if writeLine(conn, text) != nil {
			return 0, nil
		}
		return 1, nil
	})
	v.RegisterForeign("SendText", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendText(connectionId, text) requires 2 arguments")
		}
		id := toString(args[0])
		text := toString(args[1])
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		ok2 := writeLine(conn, text) == nil
		return ok2, nil
	})
	// popMessage removes and returns the first message in the connection's queue (used by reader goroutine).
	popMessage := func(id string) (string, bool) {
		connMessagesMu.Lock()
		defer connMessagesMu.Unlock()
		list := connMessages[id]
		if len(list) == 0 {
			return "", false
		}
		msg := list[0]
		connMessages[id] = list[1:]
		if len(connMessages[id]) == 0 {
			delete(connMessages, id)
		}
		return msg, true
	}
	v.RegisterForeign("Receive", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Receive(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		if msg, ok := popMessage(id); ok {
			return msg, nil
		}
		return nil, nil
	})
	v.RegisterForeign("ReceiveJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ReceiveJSON(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		msg, ok := popMessage(id)
		if !ok {
			return nil, nil
		}
		if !json.Valid([]byte(msg)) {
			return nil, nil
		}
		return msg, nil
	})
	v.RegisterForeign("ReceiveTable", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ReceiveTable(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		msg, ok := popMessage(id)
		if !ok {
			return nil, nil
		}
		var out map[string]interface{}
		if err := json.Unmarshal([]byte(msg), &out); err != nil {
			return nil, nil
		}
		return out, nil
	})
	v.RegisterForeign("ReceiveNumbers", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ReceiveNumbers(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		line, ok := popMessage(id)
		if !ok {
			return 0, nil
		}
		parts := strings.Split(line, " ")
		var nums []float64
		for i, p := range parts {
			if p == "" {
				continue
			}
			if i == 0 && (p == "n" || p == "i" || p == "f") {
				continue
			}
			f, err := strconv.ParseFloat(p, 64)
			if err != nil {
				receivedNumbersMu.Lock()
				delete(receivedNumbers, id)
				if lastNumbersConnID == id {
					lastNumbersConnID = ""
				}
				receivedNumbersMu.Unlock()
				return 0, nil
			}
			nums = append(nums, f)
		}
		receivedNumbersMu.Lock()
		receivedNumbers[id] = nums
		lastNumbersConnID = id
		receivedNumbersMu.Unlock()
		return len(nums), nil
	})
	v.RegisterForeign("GetReceivedNumber", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetReceivedNumber(index) or GetReceivedNumber(connectionId, index) requires 1 or 2 arguments")
		}
		connID := ""
		idxArg := args[0]
		if len(args) >= 2 {
			connID = toString(args[0])
			idxArg = args[1]
		}
		idx := toInt(idxArg)
		receivedNumbersMu.Lock()
		defer receivedNumbersMu.Unlock()
		if connID == "" {
			connID = lastNumbersConnID
		}
		nums := receivedNumbers[connID]
		if idx < 0 || idx >= len(nums) {
			return 0.0, nil
		}
		return nums[idx], nil
	})
	v.RegisterForeign("Disconnect", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Disconnect(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if ok {
			cleanupConnection(id, conn, false)
		}
		return nil, nil
	})
	v.RegisterForeign("SendPing", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("SendPing(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		if _, err := fmt.Fprintln(conn, "P"); err != nil {
			return false, nil
		}
		pingMu.Lock()
		pingSentAt[id] = time.Now()
		pingMu.Unlock()
		return true, nil
	})
	v.RegisterForeign("GetPing", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetPing(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		pingMu.Lock()
		ms := lastRTTMs[id]
		pingMu.Unlock()
		return ms, nil
	})

	// --- Server ---
	v.RegisterForeign("Host", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Host(port) requires 1 argument")
		}
		port := toInt(args[0])
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return nil, nil
		}
		deadlineSrc, _ := listener.(deadlineListener)
		id := addServer(listener, deadlineSrc)
		return id, nil
	})
	v.RegisterForeign("HostTLS", func(args []interface{}) (interface{}, error) {
		if len(args) < 3 {
			return nil, fmt.Errorf("HostTLS(port, certFile, keyFile) requires 3 arguments")
		}
		port := toInt(args[0])
		certFile := toString(args[1])
		keyFile := toString(args[2])
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, nil
		}
		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		baseListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return nil, nil
		}
		tcpListener, _ := baseListener.(*net.TCPListener)
		listener := tls.NewListener(baseListener, config)
		id := addServer(listener, tcpListener)
		return id, nil
	})
	v.RegisterForeign("Accept", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Accept(serverId) requires 1 argument")
		}
		sid := toString(args[0])
		netMu.Lock()
		state, ok := servers[sid]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown server: %s", sid)
		}
		conn, err := acceptServerConnection(state, 0)
		if err != nil {
			return nil, nil
		}
		if conn == nil {
			return nil, nil
		}
		netMu.Lock()
		connCounter++
		cid := fmt.Sprintf("conn_%d", connCounter)
		conns[cid] = conn
		netMu.Unlock()
		pushEvent("connect", cid, "")
		go startReader(cid, conn)
		return cid, nil
	})
	v.RegisterForeign("CloseServer", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CloseServer(serverId) requires 1 argument")
		}
		sid := toString(args[0])
		netMu.Lock()
		defer netMu.Unlock()
		if state, ok := servers[sid]; ok {
			_ = state.listener.Close()
			delete(servers, sid)
		}
		return nil, nil
	})

	// --- Rooms ---
	v.RegisterForeign("CreateRoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CreateRoom(roomId) requires 1 argument")
		}
		roomId := toString(args[0])
		netMu.Lock()
		if rooms[roomId] == nil {
			rooms[roomId] = make(map[string]bool)
		}
		netMu.Unlock()
		return 1, nil
	})
	v.RegisterForeign("JoinRoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("JoinRoom(roomId, connectionId) requires 2 arguments")
		}
		roomId := toString(args[0])
		cid := toString(args[1])
		netMu.Lock()
		defer netMu.Unlock()
		if _, ok := conns[cid]; !ok {
			return nil, fmt.Errorf("unknown connection: %s", cid)
		}
		if rooms[roomId] == nil {
			rooms[roomId] = make(map[string]bool)
		}
		rooms[roomId][cid] = true
		return 1, nil
	})
	v.RegisterForeign("LeaveRoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("LeaveRoom(connectionId) or LeaveRoom(connectionId, roomId) requires 1 or 2 arguments")
		}
		cid := toString(args[0])
		netMu.Lock()
		defer netMu.Unlock()
		if len(args) == 1 {
			for roomId, set := range rooms {
				delete(set, cid)
				if len(set) == 0 {
					delete(rooms, roomId)
				}
			}
		} else {
			roomId := toString(args[1])
			if set := rooms[roomId]; set != nil {
				delete(set, cid)
				if len(set) == 0 {
					delete(rooms, roomId)
				}
			}
		}
		return nil, nil
	})
	v.RegisterForeign("SendToRoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendToRoom(roomId, text) requires 2 arguments")
		}
		roomId := toString(args[0])
		text := toString(args[1])
		if len(text) > maxMessageSize || strings.Contains(text, "\n") || strings.Contains(text, "\r") {
			return 0, nil
		}
		netMu.Lock()
		set := rooms[roomId]
		if set == nil {
			netMu.Unlock()
			return 0, nil
		}
		cids := make([]string, 0, len(set))
		for cid := range set {
			cids = append(cids, cid)
		}
		netMu.Unlock()
		n := 0
		for _, cid := range cids {
			netMu.Lock()
			conn, ok := conns[cid]
			netMu.Unlock()
			if !ok {
				continue
			}
			if writeLine(conn, text) == nil {
				n++
			}
		}
		return n, nil
	})
	v.RegisterForeign("SendToRoomJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendToRoomJSON(roomId, jsonText) requires 2 arguments")
		}
		roomId := toString(args[0])
		text := toString(args[1])
		if len(text) > maxMessageSize || !json.Valid([]byte(text)) {
			return 0, nil
		}
		netMu.Lock()
		set := rooms[roomId]
		if set == nil {
			netMu.Unlock()
			return 0, nil
		}
		cids := make([]string, 0, len(set))
		for cid := range set {
			cids = append(cids, cid)
		}
		netMu.Unlock()
		n := 0
		for _, cid := range cids {
			netMu.Lock()
			conn, ok := conns[cid]
			netMu.Unlock()
			if !ok {
				continue
			}
			if writeLine(conn, text) == nil {
				n++
			}
		}
		return n, nil
	})
	v.RegisterForeign("SendToRoomInt", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendToRoomInt(roomId, value) requires 2 arguments")
		}
		roomId := toString(args[0])
		text := fmt.Sprintf("i %d", toInt(args[1]))
		netMu.Lock()
		set := rooms[roomId]
		if set == nil {
			netMu.Unlock()
			return 0, nil
		}
		cids := make([]string, 0, len(set))
		for cid := range set {
			cids = append(cids, cid)
		}
		netMu.Unlock()
		n := 0
		for _, cid := range cids {
			netMu.Lock()
			conn, ok := conns[cid]
			netMu.Unlock()
			if !ok {
				continue
			}
			if writeLine(conn, text) == nil {
				n++
			}
		}
		return n, nil
	})
	v.RegisterForeign("SendToRoomFloat", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendToRoomFloat(roomId, value) requires 2 arguments")
		}
		roomId := toString(args[0])
		text := fmt.Sprintf("f %g", toFloat(args[1]))
		netMu.Lock()
		set := rooms[roomId]
		if set == nil {
			netMu.Unlock()
			return 0, nil
		}
		cids := make([]string, 0, len(set))
		for cid := range set {
			cids = append(cids, cid)
		}
		netMu.Unlock()
		n := 0
		for _, cid := range cids {
			netMu.Lock()
			conn, ok := conns[cid]
			netMu.Unlock()
			if !ok {
				continue
			}
			if writeLine(conn, text) == nil {
				n++
			}
		}
		return n, nil
	})
	v.RegisterForeign("SendToRoomNumbers", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendToRoomNumbers(roomId, n1, n2, ...) requires at least 2 arguments")
		}
		roomId := toString(args[0])
		nCount := len(args) - 1
		if nCount > maxSendNumbers {
			nCount = maxSendNumbers
		}
		parts := make([]string, 0, nCount+1)
		parts = append(parts, "n")
		for i := 1; i <= nCount; i++ {
			parts = append(parts, fmt.Sprintf("%g", toFloat(args[i])))
		}
		text := strings.Join(parts, " ")
		if len(text) > maxMessageSize {
			return 0, nil
		}
		netMu.Lock()
		set := rooms[roomId]
		if set == nil {
			netMu.Unlock()
			return 0, nil
		}
		cids := make([]string, 0, len(set))
		for cid := range set {
			cids = append(cids, cid)
		}
		netMu.Unlock()
		n := 0
		for _, cid := range cids {
			netMu.Lock()
			conn, ok := conns[cid]
			netMu.Unlock()
			if !ok {
				continue
			}
			if writeLine(conn, text) == nil {
				n++
			}
		}
		return n, nil
	})
	v.RegisterForeign("GetRoomConnectionCount", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetRoomConnectionCount(roomId) requires 1 argument")
		}
		roomId := toString(args[0])
		netMu.Lock()
		n := len(rooms[roomId])
		netMu.Unlock()
		return n, nil
	})
	v.RegisterForeign("GetRoomConnectionId", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("GetRoomConnectionId(roomId, index) requires 2 arguments")
		}
		roomId := toString(args[0])
		idx := toInt(args[1])
		netMu.Lock()
		set := rooms[roomId]
		if set == nil || idx < 0 {
			netMu.Unlock()
			return "", nil
		}
		cids := make([]string, 0, len(set))
		for cid := range set {
			cids = append(cids, cid)
		}
		sort.Strings(cids)
		netMu.Unlock()
		if idx >= len(cids) {
			return "", nil
		}
		return cids[idx], nil
	})

	// --- Convenience ---
	v.RegisterForeign("IsConnected", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("IsConnected(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		_, ok := conns[id]
		netMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
	// Aliases for simpler names
	v.RegisterForeign("NetHost", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NetHost(port) requires 1 argument")
		}
		port := toInt(args[0])
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return nil, nil
		}
		deadlineSrc, _ := listener.(deadlineListener)
		id := addServer(listener, deadlineSrc)
		return id, nil
	})
	v.RegisterForeign("NetConnect", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("NetConnect(ip, port) requires 2 arguments")
		}
		host := toString(args[0])
		port := toInt(args[1])
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			return nil, nil
		}
		netMu.Lock()
		connCounter++
		cid := fmt.Sprintf("conn_%d", connCounter)
		conns[cid] = conn
		netMu.Unlock()
		go startReader(cid, conn)
		return cid, nil
	})
	v.RegisterForeign("NetSend", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("NetSend(connectionId, data) requires 2 arguments")
		}
		id := toString(args[0])
		text := toString(args[1])
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, nil
		}
		ok2 := writeLine(conn, text) == nil
		return ok2, nil
	})
	v.RegisterForeign("NetReceive", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NetReceive(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		msg, ok := popMessage(id)
		if !ok {
			return nil, nil
		}
		return msg, nil
	})
	v.RegisterForeign("NetIsConnected", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("NetIsConnected(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		_, ok := conns[id]
		netMu.Unlock()
		if ok {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("NetDisconnect", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("Disconnect", args)
	})
	v.RegisterForeign("NetCloseAll", func(args []interface{}) (interface{}, error) {
		netMu.Lock()
		ids := make([]string, 0, len(conns))
		for id := range conns {
			ids = append(ids, id)
		}
		netMu.Unlock()
		for _, id := range ids {
			_, _ = v.CallForeign("Disconnect", []interface{}{id})
		}
		return nil, nil
	})
	v.RegisterForeign("NetIsServer", func(args []interface{}) (interface{}, error) {
		netMu.Lock()
		n := len(servers)
		netMu.Unlock()
		if n > 0 {
			return 1, nil
		}
		return 0, nil
	})
	v.RegisterForeign("NetPlayerID", func(args []interface{}) (interface{}, error) {
		netMu.Lock()
		for cid := range conns {
			netMu.Unlock()
			return cid, nil
		}
		netMu.Unlock()
		return "", nil
	})
	v.RegisterForeign("NetPing", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetPing", args)
	})
	v.RegisterForeign("NetLatency", func(args []interface{}) (interface{}, error) {
		return v.CallForeign("GetPing", args)
	})
	v.RegisterForeign("GetConnectionCount", func(args []interface{}) (interface{}, error) {
		netMu.Lock()
		n := len(conns)
		netMu.Unlock()
		return n, nil
	})
	v.RegisterForeign("GetLocalIP", func(args []interface{}) (interface{}, error) {
		ifaces, err := net.Interfaces()
		if err != nil {
			return "127.0.0.1", nil
		}
		for _, iface := range ifaces {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, a := range addrs {
				s := a.String()
				if idx := strings.Index(s, "/"); idx >= 0 {
					s = s[:idx]
				}
				ip := net.ParseIP(s)
				if ip != nil && ip.To4() != nil {
					return ip.String(), nil
				}
			}
		}
		return "127.0.0.1", nil
	})
	v.RegisterForeign("AcceptTimeout", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("AcceptTimeout(serverId, timeoutMs) requires 2 arguments")
		}
		sid := toString(args[0])
		timeoutMs := toInt(args[1])
		netMu.Lock()
		state, ok := servers[sid]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown server: %s", sid)
		}
		conn, err := acceptServerConnection(state, time.Duration(timeoutMs)*time.Millisecond)
		if err != nil {
			return nil, nil
		}
		if conn == nil {
			return nil, nil
		}
		netMu.Lock()
		connCounter++
		cid := fmt.Sprintf("conn_%d", connCounter)
		conns[cid] = conn
		netMu.Unlock()
		pushEvent("connect", cid, "")
		go startReader(cid, conn)
		return cid, nil
	})

	// --- High-level event-based API ---
	v.RegisterForeign("ProcessNetworkEvents", func(args []interface{}) (interface{}, error) {
		events := drainEvents()
		if netVM == nil || netVM.Chunk() == nil {
			return nil, nil
		}
		for _, ev := range events {
			switch ev.typ {
			case "connect":
				if _, ok := netVM.Chunk().GetFunction("onclientconnect"); ok {
					if err := netVM.InvokeSub("onclientconnect", []interface{}{ev.id}); err != nil {
						return nil, err
					}
				}
			case "disconnect":
				if _, ok := netVM.Chunk().GetFunction("onclientdisconnect"); ok {
					if err := netVM.InvokeSub("onclientdisconnect", []interface{}{ev.id}); err != nil {
						return nil, err
					}
				}
			case "message":
				msg, ok := popMessage(ev.id)
				if !ok {
					continue
				}
				handled := false
				if strings.HasPrefix(msg, "E\t") {
					parts := strings.Split(msg, "\t")
					if len(parts) >= 4 {
						entityId := parts[1]
						x, _ := strconv.ParseFloat(parts[2], 64)
						y, _ := strconv.ParseFloat(parts[3], 64)
						z := 0.0
						if len(parts) >= 5 {
							z, _ = strconv.ParseFloat(parts[4], 64)
						}
						remoteEntitiesMu.Lock()
						if remoteEntities[entityId] == nil {
							remoteEntities[entityId] = make(map[string]interface{})
						}
						remoteEntities[entityId]["x"] = x
						remoteEntities[entityId]["y"] = y
						remoteEntities[entityId]["z"] = z
						remoteEntitiesMu.Unlock()
						if _, ok := netVM.Chunk().GetFunction("onentitysync"); ok {
							if err := netVM.InvokeSub("onentitysync", []interface{}{entityId, x, y, z}); err != nil {
								return nil, err
							}
						}
						handled = true
					}
				}
				if !handled && len(msg) >= 3 && msg[:2] == "R<" {
					if idx := strings.Index(msg, ">"); idx > 2 {
						rpcName := strings.ToLower(msg[2:idx])
						rest := strings.TrimSpace(msg[idx+1:])
						var rpcArgs []interface{}
						if len(rest) > 0 {
							_ = json.Unmarshal([]byte(rest), &rpcArgs)
						}
						rpcMu.Lock()
						subName, hasHandler := rpcHandlers[rpcName]
						rpcMu.Unlock()
						if hasHandler && subName != "" {
							if err := netVM.InvokeSub(subName, rpcArgs); err != nil {
								return nil, err
							}
							handled = true
						}
					}
				}
				if !handled {
					if _, ok := netVM.Chunk().GetFunction("onmessage"); ok {
						if err := netVM.InvokeSub("onmessage", []interface{}{ev.id, msg}); err != nil {
							return nil, err
						}
					}
				}
			}
		}
		return nil, nil
	})
	v.RegisterForeign("RegisterRPC", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("RegisterRPC(name, subName) requires 2 arguments")
		}
		name := strings.ToLower(toString(args[0]))
		subName := toString(args[1])
		rpcMu.Lock()
		rpcHandlers[name] = subName
		rpcMu.Unlock()
		return nil, nil
	})
	v.RegisterForeign("SendRPC", func(args []interface{}) (interface{}, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("SendRPC(connectionId, name, args...) requires at least 2 arguments")
		}
		id := toString(args[0])
		rpcName := toString(args[1])
		rpcArgs := args[2:]
		raw, err := json.Marshal(rpcArgs)
		if err != nil {
			return nil, err
		}
		payload := "R<" + rpcName + ">" + string(raw)
		if len(payload) > maxMessageSize {
			return nil, fmt.Errorf("SendRPC payload too long")
		}
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		if writeLine(conn, payload) != nil {
			return false, nil
		}
		return true, nil
	})
	v.RegisterForeign("StartServer", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("StartServer(port) requires 1 argument")
		}
		port := toInt(args[0])
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return nil, nil
		}
		deadlineSrc, _ := listener.(deadlineListener)
		id := addServer(listener, deadlineSrc)
		return id, nil
	})
	v.RegisterForeign("Broadcast", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Broadcast(text) requires 1 argument")
		}
		text := toString(args[0])
		if len(text) > maxMessageSize || strings.Contains(text, "\n") || strings.Contains(text, "\r") {
			return nil, fmt.Errorf("message too long or contains newline")
		}
		netMu.Lock()
		connList := make([]net.Conn, 0, len(conns))
		for _, c := range conns {
			connList = append(connList, c)
		}
		netMu.Unlock()
		var errs []string
		for _, conn := range connList {
			if _, err := fmt.Fprintln(conn, text); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			return nil, fmt.Errorf("broadcast partial failure: %s", strings.Join(errs, "; "))
		}
		return nil, nil
	})
	v.RegisterForeign("SyncEntity", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SyncEntity(connectionId, entityId, x, y) or SyncEntity(connectionId, entityId, x, y, z) requires 4 or 5 arguments")
		}
		id := toString(args[0])
		entityId := toString(args[1])
		x := toFloat(args[2])
		y := toFloat(args[3])
		z := 0.0
		if len(args) >= 5 {
			z = toFloat(args[4])
		}
		payload := fmt.Sprintf("E\t%s\t%g\t%g\t%g", entityId, x, y, z)
		if len(payload) > maxMessageSize {
			return false, nil
		}
		netMu.Lock()
		conn, ok := conns[id]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown connection: %s", id)
		}
		if writeLine(conn, payload) != nil {
			return false, nil
		}
		return true, nil
	})
	v.RegisterForeign("SyncEntityToRoom", func(args []interface{}) (interface{}, error) {
		if len(args) < 4 {
			return nil, fmt.Errorf("SyncEntityToRoom(roomId, entityId, x, y) or SyncEntityToRoom(roomId, entityId, x, y, z) requires 4 or 5 arguments")
		}
		roomId := toString(args[0])
		entityId := toString(args[1])
		x := toFloat(args[2])
		y := toFloat(args[3])
		z := 0.0
		if len(args) >= 5 {
			z = toFloat(args[4])
		}
		payload := fmt.Sprintf("E\t%s\t%g\t%g\t%g", entityId, x, y, z)
		if len(payload) > maxMessageSize {
			return 0, nil
		}
		netMu.Lock()
		set := rooms[roomId]
		if set == nil {
			netMu.Unlock()
			return 0, nil
		}
		cids := make([]string, 0, len(set))
		for cid := range set {
			cids = append(cids, cid)
		}
		netMu.Unlock()
		n := 0
		for _, cid := range cids {
			netMu.Lock()
			conn, ok := conns[cid]
			netMu.Unlock()
			if ok && writeLine(conn, payload) == nil {
				n++
			}
		}
		return n, nil
	})
	v.RegisterForeign("GetRemoteEntity", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetRemoteEntity(entityId) requires 1 argument")
		}
		entityId := toString(args[0])
		remoteEntitiesMu.Lock()
		m := remoteEntities[entityId]
		if m != nil {
			out := make(map[string]interface{})
			for k, v := range m {
				out[k] = v
			}
			remoteEntitiesMu.Unlock()
			return out, nil
		}
		remoteEntitiesMu.Unlock()
		return nil, nil
	})
}
