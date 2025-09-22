package version

import (
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestGetReturnsDefaults(t *testing.T) {
	info := Get()
	if info.SemVer != Version {
		t.Fatalf("expected semver %q, got %q", Version, info.SemVer)
	}
	if info.Network == "" {
		t.Fatalf("expected default network to be populated")
	}
	if info.GoVersion == "" {
		t.Fatalf("expected go version to be populated")
	}
	if info.BuildDate.IsZero() {
		t.Fatalf("expected build date to be set")
	}
}

func TestSetOverridesFields(t *testing.T) {
	ts := time.Date(2024, time.April, 2, 12, 30, 0, 0, time.UTC)
	updated, err := Set(Info{
		SemVer:    "1.2.3",
		Commit:    "abcdef1",
		BuildDate: ts,
		Network:   "qa-net",
	})
	if err != nil {
		t.Fatalf("unexpected error from set: %v", err)
	}
	if updated.SemVer != "1.2.3" {
		t.Fatalf("semver mismatch: %v", updated.SemVer)
	}
	if updated.BuildDate != ts {
		t.Fatalf("build date mismatch: %v", updated.BuildDate)
	}
	if updated.Network != "qa-net" {
		t.Fatalf("network mismatch: %v", updated.Network)
	}
	if updated.GoVersion != runtime.Version() {
		t.Fatalf("go version mismatch: %v", updated.GoVersion)
	}

	if got := Get(); got.SemVer != "1.2.3" {
		t.Fatalf("expected global semver to update, got %v", got.SemVer)
	}
}

func TestSetRejectsLeadingV(t *testing.T) {
	if _, err := Set(Info{SemVer: "v1.0.0"}); err == nil {
		t.Fatalf("expected validation error for leading v")
	}
}

func TestUserAgentFormatting(t *testing.T) {
	if _, err := Set(Info{SemVer: "2.0.0", Network: "staging"}); err != nil {
		t.Fatalf("unexpected error updating metadata: %v", err)
	}
	ua := UserAgent()
	if want := "synnergy/2.0.0"; len(ua) < len(want) || ua[:len(want)] != want {
		t.Fatalf("unexpected user agent prefix: %q", ua)
	}
	if !strings.Contains(ua, "staging") {
		t.Fatalf("expected network name in user agent: %q", ua)
	}
}
