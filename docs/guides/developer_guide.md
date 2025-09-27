# Developer Guide

This guide consolidates best practices for contributing to Synnergy.

## Environment Setup

- Install Go 1.21+ and Node.js 18+.
- Clone the repository and download dependencies with `go mod tidy`.
- Optional: run `./setup_synn.sh` to install toolchain components.

## Workflow

1. Create a feature branch from `main`.
2. Make small, focused commits touching no more than a few files.
3. Run formatting and tests:
   ```bash
   go fmt ./...
   go vet ./...
   go test ./...
   ```
4. Update or add documentation for any behaviour changes.
5. Open a pull request referencing related issues.

## Coding Standards

- Follow idiomatic Go style and run `go fmt` before committing.
- All exported symbols require GoDoc comments.
- Prefer small interfaces and package‑level examples.

## Commit Messages

Use the conventional form:
```
module: short description

Longer explanation if necessary.
```

## Code Review

- Ensure the CI pipeline passes before requesting review.
- Address reviewer comments with follow‑up commits rather than force pushes.

## Security

- Never commit secrets; `.env` files are ignored by `.gitignore`.
- Run `gosec ./...` when modifying crypto or networking code.

## Documentation

- Preview docs locally with `mkdocs serve`.
- Cross‑link related guides to keep information discoverable.

For architectural context, see the detailed whitepaper under `Whitepaper_detailed/`.

## Operational Scripts

Synnergy ships an extensive library of operational bash scripts under `scripts/` that are orchestrated through `script_launcher.sh` and the shared workflow engine. Key behaviours:

- Each `*.sh` entry point delegates to the launcher, which generates a manifest, runs prerequisite checks and can execute commands with timeouts.
- Real commands only run when the `SYN_WORKFLOW_EXECUTE_CMDS=1` environment variable is set. Without it, workflows operate in plan mode and record the intended operations. Synnergy CLI invocations additionally require `SYN_WORKFLOW_EXECUTE_CLI=1`.
- Common tasks have dedicated helpers:
  - `./scripts/install_dependencies.sh --set section=web` installs Node dependencies while recording Go module downloads in the manifest.
  - `./scripts/format_code.sh` runs `go fmt` when execution is enabled and logs formatted packages.
  - `./scripts/contract_coverage_report.sh --output build/coverage.out` ensures the output directory exists, runs the Go coverage suite and captures the manifest in `scripts/state/testing/`.
  - `./scripts/artifact_checksum.sh --set path=dist` computes SHA-256 checksums for release artifacts and persists the summary JSON.
- All workflows accept `--plan` and `--output <path>` to preview actions or override manifest destinations. Notes attached with `--note` become part of the recorded metadata.

