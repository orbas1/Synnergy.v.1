package core

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"sync"
	"testing"
	"time"
)

func TestWarfareNodeExecuteSecureCommand(t *testing.T) {
	base := NewNode("n1", "addr", NewLedger())
	wn := NewWarfareNode(base)

	cred, err := wn.IssueCommander("alpha")
	if err != nil {
		t.Fatalf("issue commander: %v", err)
	}
	privBytes, err := hex.DecodeString(cred.PrivateKey)
	if err != nil {
		t.Fatalf("decode private key: %v", err)
	}
	req := CommandRequest{
		Commander: cred.ID,
		Command:   "deploy shield",
		Timestamp: time.Now().UTC(),
		Metadata:  map[string]string{"scope": "integration"},
	}
	req.Signature = ed25519.Sign(ed25519.PrivateKey(privBytes), req.CanonicalPayload())
	record, err := wn.ExecuteSecureCommand(context.Background(), req)
	if err != nil {
		t.Fatalf("execute command: %v", err)
	}
	if !record.Accepted {
		t.Fatalf("expected command accepted: %+v", record)
	}
	if got := record.Metadata["scope"]; got != "integration" {
		t.Fatalf("metadata lost: %+v", record.Metadata)
	}
	if len(wn.CommandLog()) == 0 {
		t.Fatalf("expected audit log entry")
	}
	events := wn.Events()
	if len(events) == 0 || events[len(events)-1].Type != WarfareEventCommandAccepted {
		t.Fatalf("expected command accepted event, got %+v", events)
	}
	metrics := wn.MetricsSnapshot()
	if metrics.Commands == 0 {
		t.Fatalf("expected command metric increment, got %+v", metrics)
	}
}

func TestWarfareNodeRejectsInvalidSignature(t *testing.T) {
	wn := NewWarfareNode(NewNode("n2", "addr", NewLedger()))
	req := CommandRequest{
		Commander: "root",
		Command:   "override",
		Timestamp: time.Now().UTC(),
		Nonce:     10,
	}
	req.Signature = []byte("bad")
	if _, err := wn.ExecuteSecureCommand(context.Background(), req); err == nil {
		t.Fatalf("expected signature validation error")
	}
	events := wn.Events()
	if len(events) == 0 || events[len(events)-1].Type != WarfareEventCommandRejected {
		t.Fatalf("expected rejection event, got %+v", events)
	}
}

func TestWarfareNodeRecordLogisticsAndEvents(t *testing.T) {
	wn := NewWarfareNode(NewNode("n3", "addr", NewLedger()))
	rec, err := wn.RecordLogistics(LogisticsUpdate{AssetID: "asset1", Location: "L1", Status: "ready", Reporter: "ops"})
	if err != nil {
		t.Fatalf("record logistics: %v", err)
	}
	if rec.AssetID != "asset1" {
		t.Fatalf("unexpected record %+v", rec)
	}
	// invalid update should return error and emit rejection event
	if _, err := wn.RecordLogistics(LogisticsUpdate{}); err == nil {
		t.Fatalf("expected validation error")
	}
	logs := wn.LogisticsByAsset("asset1")
	if len(logs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(logs))
	}
	events := wn.Events()
	var seenRecorded bool
	var seenRejected bool
	for _, ev := range events {
		if ev.Type == WarfareEventLogisticsRecorded {
			seenRecorded = true
		}
		if ev.Type == WarfareEventLogisticsRejected {
			seenRejected = true
		}
	}
	if !seenRecorded || !seenRejected {
		t.Fatalf("expected recorded and rejected events, got %+v", events)
	}
}

func TestWarfareNodeSubscribeEventsReceivesBacklog(t *testing.T) {
	wn := NewWarfareNode(NewNode("n4", "addr", NewLedger()))
	if err := wn.BroadcastTactical("status green", map[string]string{"unit": "alpha"}); err != nil {
		t.Fatalf("broadcast: %v", err)
	}
	ch, cancel := wn.SubscribeEvents(2)
	defer cancel()
	select {
	case ev := <-ch:
		if ev.Type != WarfareEventTacticalBroadcast {
			t.Fatalf("unexpected event %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for backlog event")
	}

	go func() {
		_ = wn.BroadcastTactical("status yellow", nil)
	}()
	select {
	case ev := <-ch:
		if ev.Type != WarfareEventTacticalBroadcast {
			t.Fatalf("unexpected live event %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for live event")
	}
}

func TestWarfareNodeStressConcurrentLogistics(t *testing.T) {
	wn := NewWarfareNode(NewNode("n5", "addr", NewLedger()))
	const workers = 8
	const iterations = 25
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				asset := "asset" + string('A'+rune(id))
				_, err := wn.RecordLogistics(LogisticsUpdate{
					AssetID:  asset,
					Location: "L",
					Status:   "ok",
				})
				if err != nil {
					t.Errorf("record logistics: %v", err)
				}
			}
		}(i)
	}
	wg.Wait()
	metrics := wn.MetricsSnapshot()
	if metrics.Logistics != workers*iterations {
		t.Fatalf("unexpected logistics metric %+v", metrics)
	}
}
