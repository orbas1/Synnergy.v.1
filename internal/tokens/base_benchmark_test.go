package tokens

import "testing"

func BenchmarkBaseTokenTransfer(b *testing.B) {
	tok := NewBaseToken(1, "Bench", "BN", 0)
	_ = tok.Mint("alice", 1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tok.Transfer("alice", "bob", 1)
	}
}

func BenchmarkRegistryInfo(b *testing.B) {
	reg := NewRegistry()
	for i := 0; i < 1000; i++ {
		tkn := NewBaseToken(reg.NextID(), "T", "T", 0)
		reg.Register(tkn)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reg.Info(1)
	}
}
