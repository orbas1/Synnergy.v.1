# API Documentation

This directory contains reference material for Go packages within Synnergy.

To generate HTML documentation locally, run:
```bash
go doc ./...
```
or use `godoc`:
```bash
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:6060
```

## Packages
- [Core](core.md) – fundamental types and blockchain logic
