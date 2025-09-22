package cli

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"synnergy/core"
)

type stage73Persistent struct {
	Grants   []*core.GrantRecord   `json:"grants,omitempty"`
	Benefits []*core.BenefitRecord `json:"benefits,omitempty"`
	Syn500   *core.SYN500Snapshot  `json:"syn500,omitempty"`
	Syn3700  *core.SYN3700Snapshot `json:"syn3700,omitempty"`
}

var (
	stage73StateMu sync.Mutex
)

func stage73StateFile() string {
	stage73Mu.Lock()
	path := stage73Path
	stage73Mu.Unlock()
	if path != "" {
		return path
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "synnergy", "stage73_state.json")
}

func ensureStage73Loaded() error {
	stage73StateMu.Lock()
	defer stage73StateMu.Unlock()
	stage73Mu.Lock()
	if stage73Loaded {
		stage73Mu.Unlock()
		return nil
	}
	stage73Mu.Unlock()

	path := stage73StateFile()
	data, err := os.ReadFile(path)
	if err == nil {
		var state stage73Persistent
		if err := json.Unmarshal(data, &state); err == nil {
			if len(state.Grants) > 0 {
				grantRegistry = core.NewGrantRegistryFromRecords(state.Grants)
			} else {
				grantRegistry = core.NewGrantRegistry()
			}
			if len(state.Benefits) > 0 {
				benefitRegistry = core.NewBenefitRegistryFromRecords(state.Benefits)
			} else {
				benefitRegistry = core.NewBenefitRegistry()
			}
			if state.Syn500 != nil {
				syn500Token = core.RestoreSYN500Token(state.Syn500)
			} else {
				syn500Token = nil
			}
			if state.Syn3700 != nil {
				syn3700 = core.NewSYN3700FromSnapshot(state.Syn3700)
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if grantRegistry == nil {
		grantRegistry = core.NewGrantRegistry()
	}
	if benefitRegistry == nil {
		benefitRegistry = core.NewBenefitRegistry()
	}
	stage73Mu.Lock()
	stage73Loaded = true
	stage73Mu.Unlock()
	return nil
}

func persistStage73() error {
	if err := ensureStage73Loaded(); err != nil {
		return err
	}
	stage73StateMu.Lock()
	defer stage73StateMu.Unlock()
	state := stage73Persistent{
		Grants:   grantRegistry.Snapshot(),
		Benefits: benefitRegistry.Snapshot(),
	}
	if syn500Token != nil {
		state.Syn500 = syn500Token.Snapshot()
	}
	if syn3700 != nil {
		state.Syn3700 = syn3700.Snapshot()
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	path := stage73StateFile()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func loadWalletSpec(spec string) (*core.Wallet, error) {
	if spec == "" {
		return nil, errors.New("wallet specification required")
	}
	path := spec
	password := ""
	if idx := strings.LastIndex(spec, ":"); idx != -1 {
		path = spec[:idx]
		password = spec[idx+1:]
	}
	if password == "" {
		return nil, errors.New("password required")
	}
	return loadWallet(path, password)
}
