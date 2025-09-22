package core

import "time"

// Stage73State aggregates the enterprise token registries that require
// persistence across CLI, web and automation flows.
type Stage73State struct {
	GeneratedAt time.Time                `json:"generated_at"`
	SYN3700     *SYN3700Snapshot         `json:"syn3700,omitempty"`
	Grants      *GrantRegistrySnapshot   `json:"syn3800,omitempty"`
	Benefits    *BenefitRegistrySnapshot `json:"syn3900,omitempty"`
}

// CaptureStage73State creates a snapshot of the Stage 73 modules.
func CaptureStage73State(tok *SYN3700Token, grants *GrantRegistry, benefits *BenefitRegistry) Stage73State {
	state := Stage73State{GeneratedAt: time.Now().UTC()}
	if tok != nil {
		snap := tok.Snapshot()
		state.SYN3700 = &snap
	}
	if grants != nil {
		snap := grants.Snapshot()
		state.Grants = &snap
	}
	if benefits != nil {
		snap := benefits.Snapshot()
		state.Benefits = &snap
	}
	return state
}

// ApplyStage73State restores module state from a captured snapshot.
func ApplyStage73State(state Stage73State, tok *SYN3700Token, grants *GrantRegistry, benefits *BenefitRegistry) {
	if state.SYN3700 != nil && tok != nil {
		tok.Restore(*state.SYN3700)
	}
	if state.Grants != nil && grants != nil {
		grants.Restore(*state.Grants)
	}
	if state.Benefits != nil && benefits != nil {
		benefits.Restore(*state.Benefits)
	}
}
