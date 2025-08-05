package core

// Pause suspends Plasma bridge operations.
func (b *PlasmaBridge) Pause() {
	b.mu.Lock()
	b.paused = true
	b.mu.Unlock()
}

// Resume resumes Plasma bridge operations.
func (b *PlasmaBridge) Resume() {
	b.mu.Lock()
	b.paused = false
	b.mu.Unlock()
}

// Status reports whether the Plasma bridge is paused.
func (b *PlasmaBridge) Status() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.paused
}
