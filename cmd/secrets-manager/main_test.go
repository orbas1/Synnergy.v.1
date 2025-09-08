package main

import (
	"bytes"
	"io"
	"os"
	"strings"
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

func run(args ...string) string {
	old := os.Args
	os.Args = append([]string{"secrets-manager"}, args...)
	out := captureOutput(main)
	os.Args = old
	return out
}

func TestSetGet(t *testing.T) {
	if out := run("set", "k", "v"); out != "ok\n" {
		t.Fatalf("unexpected set output: %q", out)
	}
	if out := run("get", "k"); out != "v\n" {
		t.Fatalf("unexpected get output: %q", out)
	}
	if out := run("get", "missing"); !strings.Contains(out, "error") {
		t.Fatalf("expected error, got %q", out)
	}
}
