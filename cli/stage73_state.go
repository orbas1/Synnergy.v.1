package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"synnergy/core"
)

var (
	stage73Path            string
	stage73State           core.Stage73State
	stage73Loaded          bool
	stage73Syn3700Applied  bool
	stage73GrantsApplied   bool
	stage73BenefitsApplied bool
	stage73Mu              sync.Mutex
)

func init() {
	rootCmd.PersistentFlags().StringVar(&stage73Path, "stage73-state", defaultStage73Path(), "Path to Stage 73 state file")
}

func defaultStage73Path() string {
	if env := os.Getenv("SYN_STAGE73_STATE"); env != "" {
		return env
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".synnergy", "stage73_state.json")
	}
	return "stage73_state.json"
}

func stage73StatePathValue() string {
	if stage73Path == "" {
		stage73Path = defaultStage73Path()
	}
	return stage73Path
}

func loadStage73StateLocked() error {
	if stage73Loaded {
		return nil
	}
	path := stage73StatePathValue()
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		stage73State = core.Stage73State{}
		stage73Loaded = true
		return nil
	}
	if err != nil {
		return fmt.Errorf("load stage73 state: %w", err)
	}
	var state core.Stage73State
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("decode stage73 state: %w", err)
	}
	stage73State = state
	stage73Loaded = true
	return nil
}

func persistStage73StateLocked() error {
	state := core.CaptureStage73State(syn3700, grantRegistry, benefitRegistry)
	state.GeneratedAt = time.Now().UTC()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal stage73 state: %w", err)
	}
	path := stage73StatePathValue()
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("ensure stage73 dir: %w", err)
		}
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write stage73 state: %w", err)
	}
	stage73State = state
	return nil
}

func ensureSyn3700Loaded() error {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	if err := loadStage73StateLocked(); err != nil {
		return err
	}
	if stage73Syn3700Applied {
		return nil
	}
	if stage73State.SYN3700 != nil {
		if syn3700 == nil {
			syn3700 = core.NewSYN3700Token(stage73State.SYN3700.Name, stage73State.SYN3700.Symbol)
		}
		syn3700.Restore(*stage73State.SYN3700)
	}
	stage73Syn3700Applied = true
	return nil
}

func persistSyn3700() error {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	stage73Syn3700Applied = true
	return persistStage73StateLocked()
}

func ensureGrantRegistryLoaded() error {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	if err := loadStage73StateLocked(); err != nil {
		return err
	}
	if stage73GrantsApplied {
		return nil
	}
	if stage73State.Grants != nil {
		grantRegistry.Restore(*stage73State.Grants)
	}
	stage73GrantsApplied = true
	return nil
}

func persistGrantRegistry() error {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	stage73GrantsApplied = true
	return persistStage73StateLocked()
}

func ensureBenefitRegistryLoaded() error {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	if err := loadStage73StateLocked(); err != nil {
		return err
	}
	if stage73BenefitsApplied {
		return nil
	}
	if stage73State.Benefits != nil {
		benefitRegistry.Restore(*stage73State.Benefits)
	}
	stage73BenefitsApplied = true
	return nil
}

func persistBenefitRegistry() error {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	stage73BenefitsApplied = true
	return persistStage73StateLocked()
}

func setStage73StatePath(path string) {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	stage73Path = path
	stage73Loaded = false
	stage73Syn3700Applied = false
	stage73GrantsApplied = false
	stage73BenefitsApplied = false
}

func resetStage73LoadedForTests() {
	stage73Mu.Lock()
	defer stage73Mu.Unlock()
	stage73Loaded = false
	stage73Syn3700Applied = false
	stage73GrantsApplied = false
	stage73BenefitsApplied = false
	stage73State = core.Stage73State{}
	syn3700 = nil
	grantRegistry = core.NewGrantRegistry()
	benefitRegistry = core.NewBenefitRegistry()
}
