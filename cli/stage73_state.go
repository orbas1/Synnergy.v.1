package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"synnergy/core"
)

// stage73StateFile encapsulates persisted module state used across Stage 73 CLI components.
type stage73StateFile struct {
	Modules map[string]json.RawMessage `json:"modules"`
	Meta    map[string]string          `json:"meta,omitempty"`
}

var (
	stage73PathMu   sync.RWMutex
	stage73StateMu  sync.RWMutex
	stage73State    = stage73StateFile{Modules: make(map[string]json.RawMessage)}
	stage73StateErr error
	stage73Once     sync.Once
	stage73Path     string
)

func init() {
	setStage73StatePath(filepath.Join(os.TempDir(), "synnergy_stage73_state.json"))
}

// setStage73StatePath overrides the persistence path used for Stage 73 module state.
// Tests rely on this to isolate their state files from developer machines. The path can
// reference either an existing file or a new location; missing parent directories will
// be created on demand when the state is persisted.
func setStage73StatePath(path string) {
	stage73PathMu.Lock()
	defer stage73PathMu.Unlock()
	stage73Path = filepath.Clean(path)
	stage73Once = sync.Once{}
	stage73StateErr = nil
	stage73StateMu.Lock()
	stage73State = stage73StateFile{Modules: make(map[string]json.RawMessage)}
	stage73StateMu.Unlock()
	grantRegistryOnce = sync.Once{}
	grantRegistryErr = nil
	grantRegistry = core.NewGrantRegistry()
	grantEngineOnce = sync.Once{}
	grantEngineErr = nil
	grantOrchestrator = nil
}

// resetStage73LoadedForTests clears the cached loader flag so the next call to
// ensureStage73StateLoaded re-reads any on-disk state. This mirrors behaviour that CLI
// commands rely on when the binary is executed repeatedly but the same process handles
// multiple invocations during tests.
func resetStage73LoadedForTests() {
	stage73PathMu.Lock()
	stage73Once = sync.Once{}
	stage73StateErr = nil
	stage73PathMu.Unlock()
	stage73StateMu.Lock()
	stage73State = stage73StateFile{Modules: make(map[string]json.RawMessage)}
	stage73StateMu.Unlock()
	grantRegistryOnce = sync.Once{}
	grantRegistryErr = nil
	grantRegistry = core.NewGrantRegistry()
	grantEngineOnce = sync.Once{}
	grantEngineErr = nil
	grantOrchestrator = nil
}

// ensureStage73StateLoaded lazily initialises the shared state store from disk.
func ensureStage73StateLoaded() error {
	stage73PathMu.RLock()
	path := stage73Path
	stage73PathMu.RUnlock()
	stage73Once.Do(func() {
		if path == "" {
			stage73StateErr = errors.New("stage73 state path not configured")
			return
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			stage73StateErr = fmt.Errorf("create stage73 state dir: %w", err)
			return
		}
		data, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			stage73StateMu.Lock()
			if stage73State.Modules == nil {
				stage73State.Modules = make(map[string]json.RawMessage)
			}
			stage73StateMu.Unlock()
			return
		}
		if err != nil {
			stage73StateErr = fmt.Errorf("read stage73 state: %w", err)
			return
		}
		if len(data) == 0 {
			stage73StateMu.Lock()
			if stage73State.Modules == nil {
				stage73State.Modules = make(map[string]json.RawMessage)
			}
			stage73StateMu.Unlock()
			return
		}
		var decoded stage73StateFile
		if err := json.Unmarshal(data, &decoded); err != nil {
			stage73StateErr = fmt.Errorf("decode stage73 state: %w", err)
			return
		}
		stage73StateMu.Lock()
		stage73State = decoded
		if stage73State.Modules == nil {
			stage73State.Modules = make(map[string]json.RawMessage)
		}
		stage73StateMu.Unlock()
	})
	return stage73StateErr
}

// stage73ReadModule unmarshals the stored JSON blob for the requested module into out.
// If out is nil the method simply reports whether module data exists.
func stage73ReadModule(module string, out interface{}) (bool, error) {
	if err := ensureStage73StateLoaded(); err != nil {
		return false, err
	}
	stage73StateMu.RLock()
	raw, ok := stage73State.Modules[module]
	stage73StateMu.RUnlock()
	if !ok {
		return false, nil
	}
	if out == nil {
		return true, nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return false, fmt.Errorf("decode module %s: %w", module, err)
	}
	return true, nil
}

// stage73WriteModule persists a module snapshot back to disk. When value is nil the
// module entry is removed which mirrors the behaviour expected by CLI reset flows.
func stage73WriteModule(module string, value interface{}) error {
	if err := ensureStage73StateLoaded(); err != nil {
		return err
	}
	stage73StateMu.Lock()
	defer stage73StateMu.Unlock()
	if stage73State.Modules == nil {
		stage73State.Modules = make(map[string]json.RawMessage)
	}
	if value == nil {
		delete(stage73State.Modules, module)
	} else {
		raw, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("encode module %s: %w", module, err)
		}
		stage73State.Modules[module] = raw
	}
	return persistStage73StateLocked()
}

func persistStage73StateLocked() error {
	snapshot := stage73State
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("encode stage73 state: %w", err)
	}
	stage73PathMu.RLock()
	path := stage73Path
	stage73PathMu.RUnlock()
	if path == "" {
		return errors.New("stage73 state path not configured")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create stage73 state dir: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// stage73Clear removes all persisted module data. Intended for cleanup commands and tests.
func stage73Clear() error {
	stage73StateMu.Lock()
	stage73State = stage73StateFile{Modules: make(map[string]json.RawMessage)}
	stage73StateMu.Unlock()
	if err := persistStage73StateLocked(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}
