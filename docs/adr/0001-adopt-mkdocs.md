# 0001: Adopt MkDocs for Project Documentation

Date: 2025-08-05

## Status
Accepted

## Context
Documentation was previously scattered across the repository and difficult to discover. A central documentation site simplifies maintenance and onboarding.

## Decision
Use MkDocs with a top-level `docs/` directory to host all guides and to maintain architecture decision records.

## Consequences
- Contributors author docs in Markdown and preview them locally with `mkdocs serve`.
- Guides are centralized under `docs/guides` for easy navigation.
- ADRs live under `docs/adr` to track significant technical decisions over time.
