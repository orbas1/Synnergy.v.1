package tokens

import (
	"testing"
	"time"
)

func TestSYN2700Distribute(t *testing.T) {
	tok := NewSYN2700Token()
	tok.AddHolder("A", 10)
	tok.AddHolder("B", 30)
	dist := tok.Distribute(40)
	if dist["A"] != 10 || dist["B"] != 30 {
		t.Fatalf("unexpected distribution: %v", dist)
	}
}

func TestSYN3200Convert(t *testing.T) {
	tok := NewSYN3200Token(2.0)
	if v := tok.Convert(5); v != 10 {
		t.Fatalf("convert result %d", v)
	}
}

func TestSYN3600Weight(t *testing.T) {
	tok := NewSYN3600Token()
	tok.SetWeight("A", 5)
	if tok.Weight("A") != 5 {
		t.Fatal("weight not set")
	}
}

func TestSYN3800Cap(t *testing.T) {
	tok := NewSYN3800Token(100)
	if err := tok.Mint(80); err != nil {
		t.Fatal(err)
	}
	if err := tok.Mint(30); err == nil {
		t.Fatal("expected cap exceeded")
	}
}

func TestSYN3900Vesting(t *testing.T) {
	tok := NewSYN3900Token()
	now := time.Now()
	tok.Grant("A", 10, now.Add(-time.Hour))
	if r := tok.Release("A", now); r != 10 {
		t.Fatalf("release %d", r)
	}
}

func TestSYN500Points(t *testing.T) {
	tok := NewSYN500Token()
	now := time.Now()
	tok.Mint("A", 20, now.Add(time.Hour))
	if r := tok.Redeem("A", now); r != 20 {
		t.Fatalf("redeem %d", r)
	}
}

func TestSYN5000Transfer(t *testing.T) {
	tok := NewSYN5000Token()
	tok.Mint("chain1", "alice", 50)
	if err := tok.Transfer("chain1", "alice", "chain2", "bob", 30); err != nil {
		t.Fatal(err)
	}
	if tok.Balance("chain2", "bob") != 30 {
		t.Fatalf("unexpected balance")
	}
}
