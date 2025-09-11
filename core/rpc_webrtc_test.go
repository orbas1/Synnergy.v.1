package core

import "testing"

func TestWebRTCRPC(t *testing.T) {
	rpc := NewWebRTCRPC()
	recv1 := rpc.Connect("peer1")
	recv2 := rpc.Connect("peer2")

	if ok := rpc.Send("peer1", []byte("hello")); !ok {
		t.Fatalf("send failed")
	}
	msg := <-recv1
	if string(msg) != "hello" {
		t.Fatalf("unexpected message: %s", msg)
	}

	if n := rpc.Broadcast([]byte("all")); n != 2 {
		t.Fatalf("broadcast count %d", n)
	}
	<-recv1
	<-recv2

	peers := rpc.Peers()
	if len(peers) != 2 {
		t.Fatalf("expected two peers")
	}
}
