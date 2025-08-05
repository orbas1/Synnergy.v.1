package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	synnergy "synnergy"
)

// main launches a standalone watchtower node. The node periodically
// collects system metrics and logs them. Shutdown is handled gracefully
// on SIGINT/SIGTERM signals.
func main() {
	logger := log.New(os.Stdout, "watchtower: ", log.LstdFlags|log.Lmicroseconds)
	node := synnergy.NewWatchtowerNode("watchtower-1", logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// capture OS signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Println("shutdown signal received")
		cancel()
	}()

	if err := node.Start(ctx); err != nil {
		logger.Fatalf("failed to start watchtower node: %v", err)
	}

	// block until context is cancelled
	<-ctx.Done()

	// allow some time for clean shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	done := make(chan struct{})
	go func() {
		if err := node.Stop(); err != nil {
			logger.Printf("error stopping watchtower node: %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
		logger.Println("watchtower node stopped")
	case <-shutdownCtx.Done():
		logger.Println("timeout waiting for watchtower node to stop")
	}
}
