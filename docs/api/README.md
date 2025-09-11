# API Documentation

This directory hosts reference material for packages that make up the Synnergy blockchain. Each markdown file summarises the exported types and functions of a package and links back to the source for deeper exploration.

## Generating Documentation

Several options exist for browsing the Go API locally:

### `go doc`

```bash
go doc ./...
```

Prints package level comments to the terminal.

### `godoc` server

```bash
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
```

Navigate to <http://localhost:6060/pkg/> to interact with full HTML documentation.

### pkgsite container

A containerised viewer can be launched without installing Go tooling:

```bash
docker run --rm -p 8080:8080 golang/pkgsite -http=:8080
```

Browse to <http://localhost:8080> and point the interface at the repository path.

## Package Index

- [Core](core.md) – fundamental types, consensus, and blockchain logic
- Additional packages are documented inline and may be added here as the project grows.

## Contribution Guidelines

- Exported identifiers **must** have GoDoc comments.
- Keep examples in sync with code; run `go test` to verify snippets.
- Regenerate documentation whenever public APIs change.

