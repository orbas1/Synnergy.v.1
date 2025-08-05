package synnergy

import "testing"

func BenchmarkTransactionManagerLockAndMint(b *testing.B) {
    m := NewTransactionManager()
    for i := 0; i < b.N; i++ {
        m.LockAndMint("bridge", "asset", 100, "proof")
    }
}

func BenchmarkTransactionManagerBurnAndRelease(b *testing.B) {
    m := NewTransactionManager()
    for i := 0; i < b.N; i++ {
        m.BurnAndRelease("bridge", "to", "asset", 100)
    }
}

func BenchmarkTransactionManagerListTransactions(b *testing.B) {
    m := NewTransactionManager()
    for i := 0; i < b.N; i++ {
        m.LockAndMint("bridge", "asset", 100, "proof")
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m.ListTransactions()
    }
}

func BenchmarkTransactionManagerGetTransaction(b *testing.B) {
    m := NewTransactionManager()
    ids := make([]string, b.N)
    for i := 0; i < b.N; i++ {
        ids[i] = m.LockAndMint("bridge", "asset", 100, "proof")
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = m.GetTransaction(ids[i])
    }
}

