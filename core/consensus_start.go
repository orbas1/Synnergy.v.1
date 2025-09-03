package core

import (
	"context"
	"sync/atomic"
	"time"

	"synnergy/internal/telemetry"
)

// ConsensusService runs the consensus engine for a node in the background.
type ConsensusService struct {
	node    *Node
	running int32
	quit    chan struct{}
}

// NewConsensusService creates a new service for the given node.
func NewConsensusService(n *Node) *ConsensusService {
	return &ConsensusService{node: n, quit: make(chan struct{})}
}

// Start begins the mining loop at the specified interval. The loop stops when
// Stop is called or the provided context is cancelled.
func (s *ConsensusService) Start(ctx context.Context, interval time.Duration) {
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return
	}
	go func() {
		ctx, span := telemetry.Tracer().Start(ctx, "ConsensusService.Start")
		defer span.End()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.node.MineBlock()
			case <-s.quit:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop halts the running consensus loop.
func (s *ConsensusService) Stop() {
	if atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		close(s.quit)
		s.quit = make(chan struct{})
	}
}

// Info returns the current blockchain height and whether the service is running.
func (s *ConsensusService) Info() (height int, running bool) {
	if s.node != nil {
		height = len(s.node.Blockchain)
	}
	running = atomic.LoadInt32(&s.running) == 1
	return
}
