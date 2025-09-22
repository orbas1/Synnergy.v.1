package security

import (
	"math"
	"sort"
	"sync"
	"time"
)

// MitigationConfig configures the adaptive DDoS protection subsystem. Window
// defines the observation window for rate calculations, while BlockDuration
// determines how long an offending IP remains quarantined.
type MitigationConfig struct {
	Window         time.Duration
	MaxRequests    int
	BurstAllowance int
	BlockDuration  time.Duration
	ScoreWeight    float64
}

// clientState stores per-IP metrics that feed the adaptive scoring model.
type clientState struct {
	timestamps []time.Time
	score      float64
	blocked    time.Time
	lastSeen   time.Time
}

// DDoSMitigator performs in-memory statistical analysis of inbound requests to
// determine if an IP should be rate limited or blocked. The implementation is
// entirely deterministic, making it suitable for the CLI status dashboard and
// the JS administration UI.
type DDoSMitigator struct {
	mu      sync.RWMutex
	cfg     MitigationConfig
	clients map[string]*clientState
}

// NewDDoSMitigator creates a mitigator with sane defaults. The defaults are
// tuned for validator RPC endpoints and can be adjusted at runtime.
func NewDDoSMitigator(cfg MitigationConfig) *DDoSMitigator {
	if cfg.Window == 0 {
		cfg.Window = time.Minute
	}
	if cfg.MaxRequests == 0 {
		cfg.MaxRequests = 100
	}
	if cfg.BurstAllowance == 0 {
		cfg.BurstAllowance = 25
	}
	if cfg.BlockDuration == 0 {
		cfg.BlockDuration = time.Minute * 5
	}
	if cfg.ScoreWeight == 0 {
		cfg.ScoreWeight = 0.85
	}
	return &DDoSMitigator{
		cfg:     cfg,
		clients: make(map[string]*clientState),
	}
}

// Allow reports whether the request identified by ip should be processed. The
// function updates the adaptive score and purges stale samples.
func (d *DDoSMitigator) Allow(ip string, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	state := d.ensureClient(ip)
	if now.Before(state.blocked) {
		return false
	}
	d.trim(state, now)
	state.timestamps = append(state.timestamps, now)
	state.lastSeen = now
	over := float64(len(state.timestamps) - d.cfg.MaxRequests)
	if over < 0 {
		over = 0
	}
	state.score = d.cfg.ScoreWeight*state.score + (1-d.cfg.ScoreWeight)*over
	if len(state.timestamps) > d.cfg.MaxRequests+d.cfg.BurstAllowance || state.score > float64(d.cfg.BurstAllowance) {
		state.blocked = now.Add(d.cfg.BlockDuration)
		return false
	}
	return true
}

// Block manually quarantines an IP address.
func (d *DDoSMitigator) Block(ip string, until time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	state := d.ensureClient(ip)
	state.blocked = until
	state.lastSeen = time.Now().UTC()
}

// IsBlocked reports whether the IP is currently quarantined.
func (d *DDoSMitigator) IsBlocked(ip string, now time.Time) bool {
	d.mu.RLock()
	state, ok := d.clients[ip]
	d.mu.RUnlock()
	if !ok {
		return false
	}
	return now.Before(state.blocked)
}

// Score returns the current adaptive score for the IP. The score decays over
// time as requests slow down.
func (d *DDoSMitigator) Score(ip string) float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if state, ok := d.clients[ip]; ok {
		return state.score
	}
	return 0
}

// Snapshot returns a deterministic list of clients ordered by score. It is used
// by the CLI, monitoring stack and governance dashboards.
func (d *DDoSMitigator) Snapshot(now time.Time) []ClientSnapshot {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]ClientSnapshot, 0, len(d.clients))
	for ip, state := range d.clients {
		remaining := time.Duration(0)
		if now.Before(state.blocked) {
			remaining = state.blocked.Sub(now)
		}
		out = append(out, ClientSnapshot{
			IP:           ip,
			Score:        state.score,
			BlockedUntil: remaining,
			Requests:     len(state.timestamps),
			LastSeen:     state.lastSeen,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if math.Abs(out[i].Score-out[j].Score) > 0.001 {
			return out[i].Score > out[j].Score
		}
		return out[i].IP < out[j].IP
	})
	return out
}

// ClientSnapshot provides aggregated metrics for a single IP.
type ClientSnapshot struct {
	IP           string
	Score        float64
	BlockedUntil time.Duration
	Requests     int
	LastSeen     time.Time
}

func (d *DDoSMitigator) ensureClient(ip string) *clientState {
	state, ok := d.clients[ip]
	if !ok {
		state = &clientState{}
		d.clients[ip] = state
	}
	return state
}

func (d *DDoSMitigator) trim(state *clientState, now time.Time) {
	if len(state.timestamps) == 0 {
		return
	}
	cutoff := now.Add(-d.cfg.Window)
	idx := sort.Search(len(state.timestamps), func(i int) bool {
		return state.timestamps[i].After(cutoff)
	})
	if idx > 0 {
		state.timestamps = append([]time.Time(nil), state.timestamps[idx:]...)
	}
	if now.After(state.blocked) {
		state.blocked = time.Time{}
	}
	// Score decays with inactivity
	if len(state.timestamps) == 0 {
		state.score *= d.cfg.ScoreWeight
	}
}
