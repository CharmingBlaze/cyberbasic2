// Package net provides TCP client/server bindings for multiplayer: Connect, Send, Receive, Disconnect, Host, Accept, CloseServer.
package net

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

var (
	conns             = make(map[string]net.Conn)
	readers           = make(map[string]*bufio.Reader)
	servers           = make(map[string]net.Listener)
	rooms             = make(map[string]map[string]bool) // roomId -> set of connectionIds
	netMu             sync.Mutex
	connCounter       int
	servCounter       int
	receivedNumbers   []float64
	receivedNumbersMu sync.Mutex
)

// RegisterNet registers TCP multiplayer functions with the VM.
func RegisterNet(v *vm.VM) {
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
		readers[id] = bufio.NewReader(conn)
		netMu.Unlock()
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
		readers[id] = bufio.NewReader(conn)
		netMu.Unlock()
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
		readers[id] = bufio.NewReader(conn)
		netMu.Unlock()
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
	v.RegisterForeign("Receive", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Receive(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		rd, ok := readers[id]
		netMu.Unlock()
		if !ok {
			return nil, nil
		}
		// Non-blocking: short deadline
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		}
		line, err := rd.ReadString('\n')
		if err != nil {
			if conn, ok := conns[id]; ok {
				_ = conn.SetReadDeadline(time.Time{})
			}
			return nil, nil // no data or closed
		}
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Time{})
		}
		// trim newline
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		return line, nil
	})
	v.RegisterForeign("ReceiveJSON", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ReceiveJSON(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		rd, ok := readers[id]
		netMu.Unlock()
		if !ok {
			return nil, nil
		}
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		}
		line, err := rd.ReadString('\n')
		if err != nil {
			if conn, ok := conns[id]; ok {
				_ = conn.SetReadDeadline(time.Time{})
			}
			return nil, nil
		}
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Time{})
		}
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		if !json.Valid([]byte(line)) {
			return nil, nil // not valid JSON, discard
		}
		return line, nil
	})
	v.RegisterForeign("ReceiveNumbers", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("ReceiveNumbers(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		rd, ok := readers[id]
		netMu.Unlock()
		if !ok {
			return 0, nil
		}
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		}
		line, err := rd.ReadString('\n')
		if err != nil {
			if conn, ok := conns[id]; ok {
				_ = conn.SetReadDeadline(time.Time{})
			}
			return 0, nil
		}
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Time{})
		}
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
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
				receivedNumbers = nil
				receivedNumbersMu.Unlock()
				return 0, nil
			}
			nums = append(nums, f)
		}
		receivedNumbersMu.Lock()
		receivedNumbers = nums
		receivedNumbersMu.Unlock()
		return len(nums), nil
	})
	v.RegisterForeign("GetReceivedNumber", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("GetReceivedNumber(index) requires 1 argument")
		}
		idx := toInt(args[0])
		receivedNumbersMu.Lock()
		defer receivedNumbersMu.Unlock()
		if idx < 0 || idx >= len(receivedNumbers) {
			return 0.0, nil
		}
		return receivedNumbers[idx], nil
	})
	v.RegisterForeign("Disconnect", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Disconnect(connectionId) requires 1 argument")
		}
		id := toString(args[0])
		netMu.Lock()
		defer netMu.Unlock()
		if conn, ok := conns[id]; ok {
			_ = conn.Close()
			delete(conns, id)
			delete(readers, id)
		}
		// Remove connection from all rooms
		for roomId, set := range rooms {
			delete(set, id)
			if len(set) == 0 {
				delete(rooms, roomId)
			}
		}
		return nil, nil
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
		netMu.Lock()
		servCounter++
		id := fmt.Sprintf("server_%d", servCounter)
		servers[id] = listener
		netMu.Unlock()
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
		listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), config)
		if err != nil {
			return nil, nil
		}
		netMu.Lock()
		servCounter++
		id := fmt.Sprintf("server_%d", servCounter)
		servers[id] = listener
		netMu.Unlock()
		return id, nil
	})
	v.RegisterForeign("Accept", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("Accept(serverId) requires 1 argument")
		}
		sid := toString(args[0])
		netMu.Lock()
		listener, ok := servers[sid]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown server: %s", sid)
		}
		conn, err := listener.Accept()
		if err != nil {
			return nil, nil
		}
		netMu.Lock()
		connCounter++
		cid := fmt.Sprintf("conn_%d", connCounter)
		conns[cid] = conn
		readers[cid] = bufio.NewReader(conn)
		netMu.Unlock()
		return cid, nil
	})
	v.RegisterForeign("CloseServer", func(args []interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("CloseServer(serverId) requires 1 argument")
		}
		sid := toString(args[0])
		netMu.Lock()
		defer netMu.Unlock()
		if listener, ok := servers[sid]; ok {
			_ = listener.Close()
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
		netMu.Lock()
		servCounter++
		id := fmt.Sprintf("server_%d", servCounter)
		servers[id] = listener
		netMu.Unlock()
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
		readers[cid] = bufio.NewReader(conn)
		netMu.Unlock()
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
		netMu.Lock()
		rd, ok := readers[id]
		netMu.Unlock()
		if !ok {
			return nil, nil
		}
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
		}
		line, err := rd.ReadString('\n')
		if err != nil {
			if conn, ok := conns[id]; ok {
				_ = conn.SetReadDeadline(time.Time{})
			}
			return nil, nil
		}
		if conn, ok := conns[id]; ok {
			_ = conn.SetReadDeadline(time.Time{})
		}
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		return line, nil
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
		listener, ok := servers[sid]
		netMu.Unlock()
		if !ok {
			return nil, fmt.Errorf("unknown server: %s", sid)
		}
		// TCPListener has SetDeadline; net.Listener interface does not
		if tcpListener, ok := listener.(*net.TCPListener); ok {
			tcpListener.SetDeadline(time.Now().Add(time.Duration(timeoutMs) * time.Millisecond))
			defer tcpListener.SetDeadline(time.Time{})
		}
		conn, err := listener.Accept()
		if err != nil {
			return nil, nil
		}
		netMu.Lock()
		connCounter++
		cid := fmt.Sprintf("conn_%d", connCounter)
		conns[cid] = conn
		readers[cid] = bufio.NewReader(conn)
		netMu.Unlock()
		return cid, nil
	})
}
