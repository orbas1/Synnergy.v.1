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

func run(ctx context.Context, node *synnergy.WatchtowerNode, logger *log.Logger, signals <-chan os.Signal, shutdownTimeout time.Duration) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-signals:
			if !ok {
				cancel()
				return
			}
			if logger != nil {
				logger.Println("shutdown signal received")
			}
			cancel()
		}
	}()

	if err := node.Start(ctx); err != nil {
		return err
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	done := make(chan struct{})
	go func() {
		if err := node.Stop(); err != nil && logger != nil {
			logger.Printf("error stopping watchtower node: %v", err)
		}
		close(done)
	}()

	select {
	case <-done:
		if logger != nil {
			logger.Println("watchtower node stopped")
		}
		return nil
	case <-shutdownCtx.Done():
		if logger != nil {
			logger.Println("timeout waiting for watchtower node to stop")
		}
		return shutdownCtx.Err()
	}
}

// main launches a standalone watchtower node. The node periodically
// collects system metrics and logs them. Shutdown is handled gracefully
// on SIGINT/SIGTERM signals.
func main() {
	logger := log.New(os.Stdout, "watchtower: ", log.LstdFlags|log.Lmicroseconds)
	node := synnergy.NewWatchtowerNode("watchtower-1", logger)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	if err := run(context.Background(), node, logger, sigCh, 5*time.Second); err != nil {
		logger.Printf("watchtower shutdown error: %v", err)
	}
}
