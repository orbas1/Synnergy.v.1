package cli

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	// Ensure any network services are stopped so tests exit promptly.
	network.Stop()
	os.Exit(code)
}
