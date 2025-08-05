package core

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// mockState is an in-memory implementation of StateRW used for testing.
type mockState struct {
	balances map[Address]uint64
	kv       map[string][]byte
}

type mockIter struct {
	keys []string
	idx  int
	kv   map[string][]byte
}

func newMockState() *mockState {
	return &mockState{balances: make(map[Address]uint64), kv: make(map[string][]byte)}
}

func (m *mockState) Transfer(from, to Address, amount uint64) error {
	if m.balances[from] < amount {
		return fmt.Errorf("insufficient balance")
	}
	m.balances[from] -= amount
	m.balances[to] += amount
	return nil
}

func (m *mockState) SetState(key, value []byte) { m.kv[string(key)] = value }

func (m *mockState) GetState(key []byte) ([]byte, error) {
	v, ok := m.kv[string(key)]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return v, nil
}

func (m *mockState) HasState(key []byte) (bool, error) {
	_, ok := m.kv[string(key)]
	return ok, nil
}

func (m *mockState) PrefixIterator(prefix []byte) StateIterator {
	keys := make([]string, 0)
	for k := range m.kv {
		if strings.HasPrefix(k, string(prefix)) {
			keys = append(keys, k)
		}
	}
	return &mockIter{keys: keys, kv: m.kv}
}

func (m *mockState) BalanceOf(addr Address) uint64 { return m.balances[addr] }

func (it *mockIter) Next() bool {
	if it.idx >= len(it.keys) {
		return false
	}
	it.idx++
	return true
}

func (it *mockIter) Value() []byte { return it.kv[it.keys[it.idx-1]] }

// mockElectorate is a trivial electorate that always returns true.
type mockElectorate struct{}

func (mockElectorate) IsIDTokenHolder(addr Address) bool { return true }

func TestCharityCategoryString(t *testing.T) {
	cases := []struct {
		cat  CharityCategory
		want string
	}{
		{HungerRelief, "HungerRelief"},
		{ChildrenHelp, "ChildrenHelp"},
		{WildlifeHelp, "WildlifeHelp"},
		{SeaSupport, "SeaSupport"},
		{DisasterSupport, "DisasterSupport"},
		{WarSupport, "WarSupport"},
		{CharityCategory(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.cat.String(); got != c.want {
			t.Fatalf("expected %s got %s", c.want, got)
		}
	}
}

func TestCharityPoolFlow(t *testing.T) {
	lg := logrus.New()
	st := newMockState()
	st.balances[Address("donor")] = 100
	cp := NewCharityPool(lg, st, mockElectorate{}, time.Now())

	// Deposit funds
	if err := cp.Deposit(Address("donor"), 40); err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	if bal := st.BalanceOf(CharityPoolAccount); bal != 40 {
		t.Fatalf("expected pool balance 40 got %d", bal)
	}
	if bal := st.BalanceOf(Address("donor")); bal != 60 {
		t.Fatalf("expected donor balance 60 got %d", bal)
	}

	// Insufficient funds should error
	if err := cp.Deposit(Address("donor"), 1000); err == nil {
		t.Fatalf("expected insufficient balance error")
	}

	// Register charity and retrieve
	charity := Address("0x01")
	if err := cp.Register(charity, "Helping Hands", SeaSupport); err != nil {
		t.Fatalf("register: %v", err)
	}
	reg, ok, err := cp.GetRegistration(0, charity)
	if err != nil || !ok {
		t.Fatalf("get registration failed: %v", err)
	}
	if reg.Name != "Helping Hands" || reg.Category != SeaSupport {
		t.Fatalf("unexpected registration %+v", reg)
	}
	if reg.Cycle != 0 {
		t.Fatalf("cycle not persisted correctly")
	}

	// Vote for charity and ensure persistence
	voter := Address("0x02")
	if err := cp.Vote(voter, charity); err != nil {
		t.Fatalf("vote: %v", err)
	}
	key := []byte(fmt.Sprintf("charity:vote:%s:%s", voter.Hex(), charity.Hex()))
	if ok, _ := st.HasState(key); !ok {
		t.Fatalf("vote not recorded")
	}

	// Winners currently returns an empty slice
	winners, err := cp.Winners(1)
	if err != nil {
		t.Fatalf("winners: %v", err)
	}
	if len(winners) != 0 {
		t.Fatalf("expected no winners")
	}

	// Non-existent registration should indicate missing
	if _, ok, _ := cp.GetRegistration(0, Address("0xdead")); ok {
		t.Fatalf("expected missing registration")
	}

	// Tick is a no-op but should be callable
	cp.Tick(time.Now())
}
