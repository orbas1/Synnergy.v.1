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
	if interval <= 0 {
		interval = time.Second
	}
	go func() {
		ctx, span := telemetry.Tracer("core.consensus").Start(ctx, "ConsensusService.Start")
		defer span.End()
		timer := time.NewTimer(0)
		defer func() {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		}()
		base := interval
		for {
			select {
			case <-timer.C:
				var (
					mined   *Block
					pending int
					maxTx   int
				)
				if s.node != nil {
					mined = s.node.MineBlock()
					pending = s.node.PendingTransactionCount()
					maxTx = s.node.MaxTxPerBlock
				}
				next := nextConsensusInterval(base, pending, mined, maxTx)
				timer.Reset(next)
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

func nextConsensusInterval(base time.Duration, pending int, mined *Block, maxTxPerBlock int) time.Duration {
	const (
		minInterval       = 5 * time.Millisecond
		idleBackoffFactor = 4
		busySpeedupFactor = 2
		maxIdleInterval   = 5 * time.Second
	)

	next := base
	if next <= 0 {
		next = time.Second
	}

	if mined == nil && pending == 0 {
		idle := base * idleBackoffFactor
		if idle > maxIdleInterval {
			idle = maxIdleInterval
		}
		next = idle
	} else if maxTxPerBlock > 0 && pending > maxTxPerBlock {
		accelerated := base / busySpeedupFactor
		if accelerated < minInterval {
			accelerated = minInterval
		}
		if accelerated > 0 {
			next = accelerated
		}
	}

	if next < minInterval {
		next = minInterval
	}
	return next
}
