package version

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

// Version is the semantic version exposed for backwards compatibility with
// existing CLI tooling.  It mirrors the SemVer field of the Info structure so
// that downstream consumers updating gradually do not need to change imports in
// lockstep.
const Version = "1.0.0"

// Info captures build-time and runtime metadata describing the binary.  The
// structure is intentionally compact so that it can be serialised into CLI and
// web manifests without additional allocation churn.
//
// SemVer     - human readable semantic version of the release.
// Commit     - optional git commit identifier embedded at build time.
// BuildDate  - timestamp when the binary was produced.
// Network    - logical network identifier (e.g. mainnet, testnet) for CLI/UI
//
//	switching.
//
// GoVersion  - runtime Go version for debugging mismatch issues.
type Info struct {
	SemVer    string
	Commit    string
	BuildDate time.Time
	Network   string
	GoVersion string
}

var current atomic.Value

func init() {
	current.Store(Info{
		SemVer:    Version,
		Commit:    "development",
		BuildDate: time.Now().UTC().Truncate(time.Second),
		Network:   "synnergy-mainnet",
		GoVersion: runtime.Version(),
	})
}

// Get returns the currently configured metadata snapshot.
func Get() Info {
	v := current.Load()
	if v == nil {
		return Info{SemVer: Version}
	}
	info := v.(Info)
	if info.GoVersion == "" {
		info.GoVersion = runtime.Version()
	}
	return info
}

// Set updates the global metadata snapshot.  Fields left empty inherit from the
// existing snapshot to keep roll-outs safe when only a subset of data is
// available at build time.  Invalid updates return an error so that callers can
// surface misconfigurations before presenting incorrect data to operators.
func Set(update Info) (Info, error) {
	if update.SemVer == "" {
		update.SemVer = Version
	}
	if update.Network == "" {
		update.Network = "synnergy-mainnet"
	}
	if update.BuildDate.IsZero() {
		update.BuildDate = time.Now().UTC().Truncate(time.Second)
	}
	if update.GoVersion == "" {
		update.GoVersion = runtime.Version()
	}
	if update.SemVer[0] == 'v' {
		return Info{}, fmt.Errorf("version: semver %q must not contain leading 'v'", update.SemVer)
	}
	current.Store(update)
	return update, nil
}

// UserAgent constructs a deterministic HTTP user agent string suitable for the
// RPC client, CLI manifest generator and browser integrations.  The helper keeps
// formatting consistent across the codebase and ensures Go version and target
// network are surfaced for support teams.
func UserAgent() string {
	info := Get()
	return fmt.Sprintf("synnergy/%s (%s;%s)", info.SemVer, runtime.GOOS, info.Network)
}
