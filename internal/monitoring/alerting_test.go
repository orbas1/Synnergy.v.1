package monitoring

import "testing"

func TestAlerterAlert(t *testing.T) {
	a := NewAlerter()
	ch := make(chan string, 1)
	a.Subscribe(ch)
	a.Alert("hello")
	select {
	case msg := <-ch:
		if msg != "hello" {
			t.Fatalf("expected 'hello', got %q", msg)
		}
	default:
		t.Fatal("no message received")
	}
}
