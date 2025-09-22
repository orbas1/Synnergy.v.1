package synnergy

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type staticIDGenerator struct {
	id string
}

func (g staticIDGenerator) NewID(ConnectionSpec) (string, error) {
	if g.id == "" {
		return "", errors.New("id required")
	}
	return g.id, nil
}

type sequenceIDGenerator struct {
	next int
}

func (g *sequenceIDGenerator) NewID(ConnectionSpec) (string, error) {
	g.next++
	return fmt.Sprintf("conn-%d", g.next), nil
}

func TestOpenConnectionSuccess(t *testing.T) {
	ctx := context.Background()
	payload := []byte("payload")
	spec := ConnectionSpec{
		LocalChain:        "syn-mainnet",
		RemoteChain:       "eth-mainnet",
		LocalEndpoint:     "wss://syn",
		RemoteEndpoint:    "wss://eth",
		Signer:            "authority-1",
		HandshakePayload:  payload,
		HandshakeProof:    []byte("signature"),
		Metadata:          map[string]string{"vm": "synvm"},
		HeartbeatInterval: 5 * time.Second,
	}

	var verified bool
	verifier := SignatureVerifierFunc(func(ctx context.Context, receivedPayload, signature []byte, signer string) error {
		verified = true
		if string(receivedPayload) != string(payload) {
			t.Fatalf("unexpected payload: %q", receivedPayload)
		}
		if string(signature) != "signature" {
			t.Fatalf("unexpected signature: %q", signature)
		}
		if signer != spec.Signer {
			t.Fatalf("unexpected signer: %q", signer)
		}
		return nil
	})

	mgr := NewConnectionManager(WithSignatureVerifier(verifier), WithIDGenerator(staticIDGenerator{id: "conn-1"}))
	events, cancel := mgr.Subscribe(4)
	defer cancel()

	conn, err := mgr.OpenConnection(ctx, spec)
	if err != nil {
		t.Fatalf("open connection failed: %v", err)
	}
	if !verified {
		t.Fatalf("signature verifier was not invoked")
	}
	if conn.ID != "conn-1" {
		t.Fatalf("unexpected id: %s", conn.ID)
	}
	if conn.Status != ConnectionStatusActive {
		t.Fatalf("unexpected status: %s", conn.Status)
	}
	if conn.Spec.Metadata["vm"] != "synvm" {
		t.Fatalf("metadata not cloned")
	}
	if conn.HeartbeatInterval != spec.HeartbeatInterval {
		t.Fatalf("unexpected heartbeat: %s", conn.HeartbeatInterval)
	}

	select {
	case evt := <-events:
		if evt.Type != ConnectionEventOpened {
			t.Fatalf("unexpected event type: %s", evt.Type)
		}
		if evt.Connection == nil || evt.Connection.ID != conn.ID {
			t.Fatalf("expected event snapshot for %s", conn.ID)
		}
	case <-time.After(time.Second):
		t.Fatalf("expected open event")
	}
}

func TestOpenConnectionValidation(t *testing.T) {
	mgr := NewConnectionManager()
	_, err := mgr.OpenConnection(context.Background(), ConnectionSpec{
		RemoteChain:      "remote",
		LocalEndpoint:    "local",
		RemoteEndpoint:   "remote",
		Signer:           "signer",
		HandshakeProof:   []byte("sig"),
		HandshakePayload: []byte("payload"),
	})
	if err == nil {
		t.Fatal("expected validation error for missing local chain")
	}
}

func TestCloseConnectionTransitions(t *testing.T) {
	ctx := context.Background()
	mgr := NewConnectionManager(WithIDGenerator(staticIDGenerator{id: "conn-1"}))
	events, cancel := mgr.Subscribe(8)
	defer cancel()

	spec := ConnectionSpec{
		LocalChain:       "syn",
		RemoteChain:      "ally",
		LocalEndpoint:    "local",
		RemoteEndpoint:   "remote",
		Signer:           "authority",
		HandshakeProof:   []byte("sig"),
		HandshakePayload: []byte("payload"),
	}
	_, err := mgr.OpenConnection(ctx, spec)
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	<-events // consume opened event

	closed, err := mgr.CloseConnection(ctx, "conn-1", "rotation")
	if err != nil {
		t.Fatalf("close failed: %v", err)
	}
	if closed.Status != ConnectionStatusClosed {
		t.Fatalf("expected closed status, got %s", closed.Status)
	}
	if closed.ClosingReason != "rotation" {
		t.Fatalf("unexpected closing reason: %s", closed.ClosingReason)
	}

	var closingSeen, closedSeen bool
	timeout := time.After(time.Second)
	for !closingSeen || !closedSeen {
		select {
		case evt := <-events:
			switch evt.Type {
			case ConnectionEventClosing:
				closingSeen = true
			case ConnectionEventClosed:
				closedSeen = true
			}
		case <-timeout:
			t.Fatalf("timed out waiting for lifecycle events")
		}
	}

	if _, err := mgr.CloseConnection(ctx, "conn-1", "again"); !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("expected ErrConnectionClosed, got %v", err)
	}
}

func TestFailConnectionRecordsFault(t *testing.T) {
	mgr := NewConnectionManager(WithIDGenerator(staticIDGenerator{id: "faulty"}))
	_, err := mgr.OpenConnection(context.Background(), ConnectionSpec{
		LocalChain:       "syn",
		RemoteChain:      "ally",
		LocalEndpoint:    "local",
		RemoteEndpoint:   "remote",
		Signer:           "authority",
		HandshakeProof:   []byte("sig"),
		HandshakePayload: []byte("payload"),
	})
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	fault := ConnectionFault{Code: "timeout", Detail: "no heartbeat", Severity: "critical"}
	conn, err := mgr.FailConnection("faulty", fault)
	if err != nil {
		t.Fatalf("fail connection: %v", err)
	}
	if conn.Status != ConnectionStatusFailed {
		t.Fatalf("expected failed status, got %s", conn.Status)
	}
	if len(conn.Faults) != 1 || conn.Faults[0].Code != "timeout" {
		t.Fatalf("fault not recorded: %+v", conn.Faults)
	}
}

func TestMarkHeartbeat(t *testing.T) {
	mgr := NewConnectionManager(WithIDGenerator(staticIDGenerator{id: "hb"}))
	conn, err := mgr.OpenConnection(context.Background(), ConnectionSpec{
		LocalChain:       "syn",
		RemoteChain:      "ally",
		LocalEndpoint:    "local",
		RemoteEndpoint:   "remote",
		Signer:           "authority",
		HandshakeProof:   []byte("sig"),
		HandshakePayload: []byte("payload"),
	})
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	nextBeat := conn.LastHeartbeat.Add(10 * time.Second)
	updated, err := mgr.MarkHeartbeat("hb", nextBeat)
	if err != nil {
		t.Fatalf("mark heartbeat: %v", err)
	}
	if !updated.LastHeartbeat.Equal(nextBeat) {
		t.Fatalf("heartbeat not updated: %v", updated.LastHeartbeat)
	}

	if _, err := mgr.CloseConnection(context.Background(), "hb", "shutdown"); err != nil {
		t.Fatalf("close failed: %v", err)
	}
	if _, err := mgr.MarkHeartbeat("hb", time.Now()); !errors.Is(err, ErrConnectionClosed) {
		t.Fatalf("expected ErrConnectionClosed, got %v", err)
	}
}

func TestListConnectionsFilter(t *testing.T) {
	seqGen := &sequenceIDGenerator{}
	mgr := NewConnectionManager(WithIDGenerator(seqGen))

	_, err := mgr.OpenConnection(context.Background(), ConnectionSpec{
		LocalChain:       "syn",
		RemoteChain:      "ally",
		LocalEndpoint:    "local",
		RemoteEndpoint:   "remote",
		Signer:           "authority",
		HandshakeProof:   []byte("sig"),
		HandshakePayload: []byte("payload"),
	})
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}

	_, err = mgr.OpenConnection(context.Background(), ConnectionSpec{
		LocalChain:       "syn",
		RemoteChain:      "ally",
		LocalEndpoint:    "local",
		RemoteEndpoint:   "remote",
		Signer:           "authority",
		HandshakeProof:   []byte("sig"),
		HandshakePayload: []byte("payload"),
	})
	if err != nil {
		t.Fatalf("second open failed: %v", err)
	}

	_, err = mgr.FailConnection("conn-2", ConnectionFault{Code: "fault"})
	if err != nil {
		t.Fatalf("fail connection: %v", err)
	}

	active := mgr.ListConnections(ConnectionFilter{Statuses: []ConnectionStatus{ConnectionStatusActive}})
	if len(active) != 1 {
		t.Fatalf("expected one active connection, got %d", len(active))
	}

	ended := mgr.ListConnections(ConnectionFilter{IncludeEnded: true, Statuses: []ConnectionStatus{ConnectionStatusFailed}})
	if len(ended) != 1 || ended[0].Status != ConnectionStatusFailed {
		t.Fatalf("expected failed connection in results")
	}
}

func TestSubscribeCancellation(t *testing.T) {
	mgr := NewConnectionManager()
	events, cancel := mgr.Subscribe(1)
	cancel()
	if _, ok := <-events; ok {
		t.Fatalf("expected channel to be closed after cancel")
	}
}
