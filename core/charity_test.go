package core

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

type testElectorate struct{}

func (testElectorate) IsIDTokenHolder(addr Address) bool { return addr != "bad" }

type testState struct {
	balances map[Address]uint64
	store    map[string][]byte
}

func newTestState() *testState {
	return &testState{balances: map[Address]uint64{}, store: map[string][]byte{}}
}

func (s *testState) Transfer(from, to Address, amount uint64) error {
	s.balances[from] -= amount
	s.balances[to] += amount
	return nil
}
func (s *testState) SetState(k, v []byte)                       { s.store[string(k)] = v }
func (s *testState) GetState(k []byte) ([]byte, error)          { return s.store[string(k)], nil }
func (s *testState) HasState(k []byte) (bool, error)            { _, ok := s.store[string(k)]; return ok, nil }
func (s *testState) PrefixIterator(prefix []byte) StateIterator { return nil }
func (s *testState) BalanceOf(a Address) uint64                 { return s.balances[a] }

func TestCharityPool(t *testing.T) {
	st := newTestState()
	st.balances["donor"] = 100
	cp := NewCharityPool(logrus.New(), st, testElectorate{}, time.Now())

	if err := cp.Deposit("donor", 50); err != nil {
		t.Fatalf("deposit: %v", err)
	}
	if bal := st.BalanceOf(CharityPoolAccount); bal != 50 {
		t.Fatalf("pool balance %d", bal)
	}
	if err := cp.Register("char1", "Charity One", HungerRelief); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := cp.Vote("voter", "char1"); err != nil {
		t.Fatalf("vote: %v", err)
	}
	if err := cp.Vote("bad", "char1"); err == nil {
		t.Fatalf("expected ineligible voter error")
	}
}
