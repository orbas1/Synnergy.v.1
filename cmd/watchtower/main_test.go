package main

import (
	"context"
	"io"
	"log"
	"os"
	"syscall"
	"testing"
	"time"

	synnergy "synnergy"
)

func TestRunHandlesSignal(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	node := synnergy.NewWatchtowerNode("test", logger)
	signals := make(chan os.Signal, 1)

	result := make(chan error, 1)
	go func() {
		result <- run(context.Background(), node, logger, signals, 100*time.Millisecond)
	}()

	signals <- syscall.SIGINT

	if err := <-result; err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	if err := node.Start(context.Background()); err != nil {
		t.Fatalf("node should restart after run: %v", err)
	}
	if err := node.Stop(); err != nil {
		t.Fatalf("stop after restart: %v", err)
	}
}
