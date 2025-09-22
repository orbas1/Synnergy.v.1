# Reusable Packages

The `pkg` tree contains public libraries intended for consumption by external
applications. Code placed here must remain backward compatible and adhere to the
same error-handling guidelines enforced across the CLI, consensus engine and web
integration layers.

## Version Metadata (`pkg/version`)

The version package exposes a concurrency-safe metadata registry describing the
running build. It powers the CLI manifest, browser-based command palette and
external monitoring agents.

Key capabilities:

* `version.Version` – semantic version string preserved for backwards
  compatibility with existing tooling.
* `version.Get()` – returns a snapshot of the current build metadata including
  semantic version, commit hash, build timestamp, network identifier and Go
  runtime version.
* `version.Set(info)` – allows build pipelines to inject authoritative metadata
  at compile time. Invalid payloads are rejected to prevent corrupted release
  artefacts from reaching operators.
* `version.UserAgent()` – emits a deterministic user agent string consumed by
  RPC, CLI and web clients when negotiating capabilities with remote services.

### Embedding in binaries

```
go build -ldflags "-X 'synnergy/pkg/version.Version=1.0.0'" ./cmd/synnergy
```

Runtime systems with additional context (for example, the governance authority
deployments) should call `version.Set` during initialisation to record commit
hashes, build time and target network before serving requests.

### Manifest integration

`cli/gui_manifest.go` consumes `version.Get()` when generating the machine
readable manifest consumed by the web UI. Any future packages added here should
follow the same pattern: isolate shared functionality, document it in this file
and preserve backwards compatibility for external automation.
