package net

import (
	"bufio"
	stdnet "net"
	"testing"
	"time"

	"cyberbasic/compiler/vm"
)

type fakeTimeoutError struct{}

func (fakeTimeoutError) Error() string   { return "timeout" }
func (fakeTimeoutError) Timeout() bool   { return true }
func (fakeTimeoutError) Temporary() bool { return true }

type fakeAddr string

func (a fakeAddr) Network() string { return "tcp" }
func (a fakeAddr) String() string  { return string(a) }

type fakeListener struct {
	acceptFn  func() (stdnet.Conn, error)
	deadlines []time.Time
}

func (l *fakeListener) Accept() (stdnet.Conn, error) { return l.acceptFn() }
func (l *fakeListener) Close() error                 { return nil }
func (l *fakeListener) Addr() stdnet.Addr            { return fakeAddr("fake") }
func (l *fakeListener) SetDeadline(t time.Time) error {
	l.deadlines = append(l.deadlines, t)
	return nil
}

func resetNetGlobals() {
	netVM = nil
	eventQueue = nil
	rpcHandlers = make(map[string]string)
	pingSentAt = make(map[string]time.Time)
	lastRTTMs = make(map[string]float64)
	remoteEntities = make(map[string]map[string]interface{})
	connMessages = make(map[string][]string)
	conns = make(map[string]stdnet.Conn)
	readers = make(map[string]*bufio.Reader)
	servers = make(map[string]*serverState)
	rooms = make(map[string]map[string]bool)
	connCounter = 0
	servCounter = 0
	receivedNumbers = make(map[string][]float64)
	lastNumbersConnID = ""
}

func TestAcceptServerConnectionUsesDeadlineSource(t *testing.T) {
	listener := &fakeListener{
		acceptFn: func() (stdnet.Conn, error) {
			return nil, fakeTimeoutError{}
		},
	}
	conn, err := acceptServerConnection(&serverState{
		listener:       listener,
		deadlineSource: listener,
	}, 25*time.Millisecond)
	if err != nil {
		t.Fatalf("expected timeout to map to nil error, got %v", err)
	}
	if conn != nil {
		t.Fatalf("expected no connection on timeout")
	}
	if len(listener.deadlines) != 2 {
		t.Fatalf("expected deadline set and reset, got %d calls", len(listener.deadlines))
	}
	if listener.deadlines[0].IsZero() {
		t.Fatalf("expected first deadline to be non-zero")
	}
	if !listener.deadlines[1].IsZero() {
		t.Fatalf("expected deadline reset to zero time")
	}
}

func TestReceiveNumbersStoresPerConnectionBuffers(t *testing.T) {
	resetNetGlobals()
	v := vm.NewVM()
	RegisterNet(v)

	connMessages["c1"] = []string{"n 1 2"}
	connMessages["c2"] = []string{"n 9 8"}

	got, err := v.CallForeign("ReceiveNumbers", []interface{}{"c1"})
	if err != nil || got.(int) != 2 {
		t.Fatalf("expected c1 to parse 2 numbers, got %v, err=%v", got, err)
	}
	got, err = v.CallForeign("ReceiveNumbers", []interface{}{"c2"})
	if err != nil || got.(int) != 2 {
		t.Fatalf("expected c2 to parse 2 numbers, got %v, err=%v", got, err)
	}

	last, err := v.CallForeign("GetReceivedNumber", []interface{}{0})
	if err != nil || last.(float64) != 9 {
		t.Fatalf("expected last receive buffer to be c2, got %v, err=%v", last, err)
	}
	c1First, err := v.CallForeign("GetReceivedNumber", []interface{}{"c1", 0})
	if err != nil || c1First.(float64) != 1 {
		t.Fatalf("expected per-connection read for c1 to stay intact, got %v, err=%v", c1First, err)
	}
	c2Second, err := v.CallForeign("GetReceivedNumber", []interface{}{"c2", 1})
	if err != nil || c2Second.(float64) != 8 {
		t.Fatalf("expected per-connection read for c2 to stay intact, got %v, err=%v", c2Second, err)
	}
}
