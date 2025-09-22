package cli

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestModuleCatalogueIncludesCoreDomains(t *testing.T) {
	modules := ModuleCatalogue()
	if len(modules) == 0 {
		t.Fatalf("expected module catalogue to be populated")
	}
	found := map[string]bool{}
	for _, m := range modules {
		found[m.Command] = true
	}
	for _, cmd := range []string{"consensus", "simplevm", "wallet", "node", "authority", "modules"} {
		if !found[cmd] {
			t.Fatalf("expected command %s to be registered", cmd)
		}
	}
}

func TestModuleStatusesReportGasCoverage(t *testing.T) {
	statuses := ModuleStatuses()
	if len(statuses) == 0 {
		t.Fatalf("expected module statuses")
	}
	for _, status := range statuses {
		if status.Name == "Consensus Control" {
			if len(status.Opcodes) == 0 {
				t.Fatalf("consensus module missing opcodes")
			}
			if len(status.MissingOpcodes) > 0 {
				t.Fatalf("consensus module has undocumented opcodes: %v", status.MissingOpcodes)
			}
		}
	}
}

func TestModulesCommandJSONOutput(t *testing.T) {
	cmd := newModulesCommand()
	cmd.SetArgs([]string{"list"})
	old := jsonOutput
	jsonOutput = true
	defer func() { jsonOutput = old }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	original := os.Stdout
	os.Stdout = w
	err = cmd.Execute()
	w.Close()
	os.Stdout = original
	if err != nil {
		t.Fatalf("modules command failed: %v", err)
	}
	data, readErr := io.ReadAll(r)
	r.Close()
	if readErr != nil {
		t.Fatalf("read pipe: %v", readErr)
	}
	var parsed []ModuleStatus
	if decodeErr := json.Unmarshal(data, &parsed); decodeErr != nil {
		t.Fatalf("unexpected json output: %v\n%s", decodeErr, string(data))
	}
	if len(parsed) == 0 {
		t.Fatalf("expected parsed module statuses")
	}
}
