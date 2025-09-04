package security

// PatchManager tracks applied patches.
type PatchManager struct {
	applied []string
}

// NewPatchManager creates a PatchManager.
func NewPatchManager() *PatchManager { return &PatchManager{} }

// Apply records a patch ID.
func (p *PatchManager) Apply(id string) { p.applied = append(p.applied, id) }

// Applied returns all applied patch IDs.
func (p *PatchManager) Applied() []string {
	out := make([]string, len(p.applied))
	copy(out, p.applied)
	return out
}
