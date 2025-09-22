package cli

import (
	"strings"
	"testing"

	"synnergy/core"
)

func TestRPCWebRTCCommands(t *testing.T) {
	webrtcRPC = core.NewWebRTCRPC()
	peerChans = map[string]<-chan []byte{}
	t.Cleanup(func() {
		webrtcRPC = core.NewWebRTCRPC()
		peerChans = map[string]<-chan []byte{}
	})

	if _, err := executeCLICommand(t, "rpcwebrtc", "connect", "p1"); err != nil {
		t.Fatalf("connect: %v", err)
	}
	if len(webrtcRPC.Peers()) != 1 {
		t.Fatalf("expected 1 peer, got %d", len(webrtcRPC.Peers()))
	}

	out, err := executeCLICommand(t, "rpcwebrtc", "send", "p1", "hello")
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if strings.TrimSpace(out) != "sent" {
		t.Fatalf("unexpected send output: %q", out)
	}

	out, err = executeCLICommand(t, "rpcwebrtc", "recv", "p1")
	if err != nil {
		t.Fatalf("recv: %v", err)
	}
	if strings.TrimSpace(out) != "hello" {
		t.Fatalf("unexpected recv output: %q", out)
	}

	out, err = executeCLICommand(t, "rpcwebrtc", "recv", "p1")
	if err != nil {
		t.Fatalf("recv empty: %v", err)
	}
	if strings.TrimSpace(out) != "no message" {
		t.Fatalf("expected no message, got %q", out)
	}

	if _, err := executeCLICommand(t, "rpcwebrtc", "disconnect", "p1"); err != nil {
		t.Fatalf("disconnect: %v", err)
	}
	if len(webrtcRPC.Peers()) != 0 {
		t.Fatalf("expected no peers after disconnect, got %d", len(webrtcRPC.Peers()))
	}
}
