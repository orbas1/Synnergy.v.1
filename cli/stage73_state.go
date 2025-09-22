package cli

import (
	"sync"

	"synnergy/core"
)

var (
	stage73StatePath string
	stage73Once      sync.Once
	stage73Store     *core.Stage73Store
	stage73Err       error
)

func stage73State() (*core.Stage73Store, error) {
	stage73Once.Do(func() {
		store := core.NewStage73Store(stage73StatePath)
		stage73Err = store.Load()
		if stage73Err == nil {
			stage73Store = store
		}
	})
	return stage73Store, stage73Err
}

func saveStage73State() error {
	if stage73Store == nil {
		return nil
	}
	return stage73Store.Save()
}

func markStage73Dirty() {
	if stage73Store != nil {
		stage73Store.MarkDirty()
	}
}

func setStage73StatePath(path string) {
	stage73StatePath = path
}

func resetStage73LoadedForTests() {
	stage73Once = sync.Once{}
	stage73Store = nil
	stage73Err = nil
}
