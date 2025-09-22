package militarynodes

import "testing"

func TestSecureWarfareNode(t *testing.T) {
	node := NewSecureWarfareNode("n1", WithLogisticsLimit(2), WithCommandLimit(2))
	if err := node.SecureCommand("operator", "cmd1"); err != nil {
		t.Fatalf("secure command: %v", err)
	}
	if err := node.SecureCommand("operator", "cmd2"); err != nil {
		t.Fatalf("secure command: %v", err)
	}
	if err := node.SecureCommand("operator", "cmd3"); err != nil {
		t.Fatalf("secure command: %v", err)
	}
	cmds := node.CommandHistory()
	if len(cmds) != 2 || cmds[0].Command != "cmd2" {
		t.Fatalf("unexpected command history: %#v", cmds)
	}

	node.TrackLogistics("asset", "loc1", "ok")
	node.TrackLogistics("asset", "loc2", "ok")
	node.TrackLogistics("asset", "loc3", "ok")
	logs := node.Logistics()
	if len(logs) != 2 || logs[0].Location != "loc2" {
		t.Fatalf("unexpected logistics: %#v", logs)
	}
}

func TestSecureCommandValidation(t *testing.T) {
	node := NewSecureWarfareNode("n2")
	if err := node.SecureCommand("", "cmd"); err == nil {
		t.Fatalf("expected operator validation error")
	}
	if err := node.SecureCommand("operator", ""); err == nil {
		t.Fatalf("expected command validation error")
	}
}
