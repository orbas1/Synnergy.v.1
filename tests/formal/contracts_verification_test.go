package formal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"testing/quick"

	synnergy "synnergy"
)

type noopVM struct{}

func (noopVM) Execute([]byte, string, []byte, uint64) ([]byte, uint64, error) { return nil, 0, nil }
func (noopVM) Start() error                                                   { return nil }
func (noopVM) Stop() error                                                    { return nil }
func (noopVM) Status() bool                                                   { return true }

func TestContractsFormalVerification(t *testing.T) {
	vm := noopVM{}
	cfg := &quick.Config{MaxCount: 50}
	property := func(data []byte) bool {
		if len(data) == 0 {
			return true
		}
		reg := synnergy.NewContractRegistry(vm, nil)
		addr, err := reg.Deploy(data, `{"name":"contract"}`, 1, "owner")
		if err != nil {
			return false
		}
		sum := sha256.Sum256(data)
		if addr != hex.EncodeToString(sum[:]) {
			return false
		}
		c, ok := reg.Get(addr)
		if !ok {
			return false
		}
		return bytes.Equal(c.WASM, data)
	}
	if err := quick.Check(property, cfg); err != nil {
		t.Fatalf("contract registry invariant violated: %v", err)
	}
}
