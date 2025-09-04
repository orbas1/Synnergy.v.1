package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestMainOutput(t *testing.T) {
	got := captureOutput(main)
	expected := "secrets manager CLI placeholder\n"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
