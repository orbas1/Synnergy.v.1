package fuzz

import "testing"

func FuzzNetwork(f *testing.F) {
	f.Add([]byte("seed"))
	f.Fuzz(func(t *testing.T, data []byte) {})
}
