package core

import (
	"errors"
	"sync"
)

// PlasmaExit represents a pending withdrawal from a Plasma bridge.
type PlasmaExit struct {
	Nonce     uint64
	Owner     string
	Token     string
	Amount    uint64
	Finalized bool
}

// PlasmaBridge tracks deposits and exits for a Plasma chain.
type PlasmaBridge struct {
	mu     sync.RWMutex
	seq    uint64
	exits  map[uint64]*PlasmaExit
	paused bool
}

// NewPlasmaBridge creates a new Plasma bridge instance.
func NewPlasmaBridge() *PlasmaBridge {
	return &PlasmaBridge{exits: make(map[uint64]*PlasmaExit)}
}

// ErrBridgePaused is returned when operations are attempted while the bridge is
// paused.
var ErrBridgePaused = errors.New("plasma bridge paused")
