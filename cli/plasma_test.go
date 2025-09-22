package cli

import (
	"encoding/json"
	"testing"

	"synnergy/core"
)

func TestPlasmaStatusJSON(t *testing.T) {
	plasmaBridge = core.NewPlasmaBridge()
	plasmaJSON = false
	jsonOutput = false
	t.Cleanup(func() {
		plasmaBridge = core.NewPlasmaBridge()
		plasmaJSON = false
		jsonOutput = false
	})

	out, err := executeCLICommand(t, "plasma", "--json", "plasma-mgmt", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	var payload map[string]bool
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v, ok := payload["paused"]; !ok || v {
		t.Fatalf("expected paused=false, got %v (map: %v)", payload["paused"], payload)
	}
}
