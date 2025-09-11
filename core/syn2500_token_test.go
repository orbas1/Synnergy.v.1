package core

import (
	"sync"
	"testing"
)

func TestSyn2500Registry(t *testing.T) {
	reg := NewSyn2500Registry()
	m := NewSyn2500Member("1", "addr1", 1, nil)
	if err := reg.AddMember(m); err != nil {
		t.Fatalf("add member: %v", err)
	}
	if err := reg.AddMember(m); err != ErrMemberExists {
		t.Fatalf("expected ErrMemberExists, got %v", err)
	}
	if _, ok := reg.GetMember("1"); !ok {
		t.Fatalf("member not found")
	}
	if err := reg.RemoveMember("1"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if err := reg.RemoveMember("1"); err != ErrMemberNotFound {
		t.Fatalf("expected ErrMemberNotFound, got %v", err)
	}
}

func TestSyn2500RegistryConcurrentAdd(t *testing.T) {
	reg := NewSyn2500Registry()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			_ = reg.AddMember(NewSyn2500Member(id, id, 1, nil))
		}(string('a' + rune(i)))
	}
	wg.Wait()
	if len(reg.ListMembers()) == 0 {
		t.Fatalf("expected members to be added")
	}
}
