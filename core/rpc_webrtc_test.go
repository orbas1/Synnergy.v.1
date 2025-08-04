package core

import "testing"

func TestWebRTCRPC(t *testing.T) {
	rpc := NewWebRTCRPC()
	recv := rpc.Connect("peer1")

	if ok := rpc.Send("peer1", []byte("hello")); !ok {
		t.Fatalf("send failed")
	}
	msg := <-recv
	if string(msg) != "hello" {
		t.Fatalf("unexpected message: %s", msg)
	}
}
