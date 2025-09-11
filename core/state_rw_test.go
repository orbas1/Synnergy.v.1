package core

import (
	"errors"
	"strings"
	"testing"
)

type mockStateIter struct {
	keys []string
	idx  int
	kv   map[string][]byte
}

func (m *mockStateIter) Next() bool {
	if m.idx >= len(m.keys) {
		return false
	}
	m.idx++
	return true
}

func (m *mockStateIter) Value() []byte {
	return m.kv[m.keys[m.idx-1]]
}

type mockState struct {
	balances map[Address]uint64
	kv       map[string][]byte
}

func newMockState() *mockState {
	return &mockState{balances: make(map[Address]uint64), kv: make(map[string][]byte)}
}

func (m *mockState) Transfer(from, to Address, amount uint64) error {
	if m.balances[from] < amount {
		return errors.New("insufficient balance")
	}
	m.balances[from] -= amount
	m.balances[to] += amount
	return nil
}

func (m *mockState) SetState(key, value []byte) {
	m.kv[string(key)] = value
}

func (m *mockState) GetState(key []byte) ([]byte, error) {
	v, ok := m.kv[string(key)]
	if !ok {
		return nil, errors.New("not found")
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
	return &mockStateIter{keys: keys, kv: m.kv}
}

func (m *mockState) BalanceOf(addr Address) uint64 {
	return m.balances[addr]
}

func TestStateRW(t *testing.T) {
	s := newMockState()
	s.balances["alice"] = 10

	if err := s.Transfer("alice", "bob", 5); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	if s.BalanceOf("alice") != 5 || s.BalanceOf("bob") != 5 {
		t.Fatalf("unexpected balances: %d %d", s.BalanceOf("alice"), s.BalanceOf("bob"))
	}

	s.SetState([]byte("k1"), []byte("v1"))
	s.SetState([]byte("k2"), []byte("v2"))
	if v, _ := s.GetState([]byte("k1")); string(v) != "v1" {
		t.Fatalf("get state failed")
	}
	if ok, _ := s.HasState([]byte("k2")); !ok {
		t.Fatalf("expected key k2 to exist")
	}

	it := s.PrefixIterator([]byte("k"))
	count := 0
	for it.Next() {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 iterated values, got %d", count)
	}
}
