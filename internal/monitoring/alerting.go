package monitoring

// Alerter provides simple alert broadcasting.
type Alerter struct {
	listeners []chan string
}

// NewAlerter creates an Alerter.
func NewAlerter() *Alerter { return &Alerter{} }

// Subscribe registers a listener channel.
func (a *Alerter) Subscribe(ch chan string) {
	a.listeners = append(a.listeners, ch)
}

// Alert sends a message to all listeners.
func (a *Alerter) Alert(msg string) {
	for _, ch := range a.listeners {
		ch <- msg
	}
}
