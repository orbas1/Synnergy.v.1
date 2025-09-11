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

