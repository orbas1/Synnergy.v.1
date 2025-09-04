# Synnergy Documentation

Synnergy is a modular research blockchain. The `docs/` tree is the primary
reference for contributors and operators and is published as a static web site
with [MkDocs](https://www.mkdocs.org/). The documentation covers design
decisions, end‑user guides and benchmarking data for the core platform.

## Structure

The documentation is split into three main sections:

| Path | Description |
|------|-------------|
| [`guides/`](guides/) | Quickstart and how‑to articles for developers, node operators, and GUI users. |
| [`api/`](api/) | Package‑level references and instructions for generating API documentation. |
| [`adr/`](adr/) | Architecture Decision Records describing significant technical choices made during development. |
| [`performance_benchmarks.md`](performance_benchmarks.md) | Baseline performance figures and instructions for running benchmarks. |

Additional reference material such as opcode listings and error catalogues can
be found in the repository root. Links are provided throughout the guides where
relevant.

## Building the Site

Install MkDocs and the material theme:

```bash
pip install mkdocs mkdocs-material
```

Serve the documentation locally while editing:

```bash
mkdocs serve
```

or generate a static site under `site/`:

```bash
mkdocs build
```

## Contributing

Documentation lives alongside the code and should be updated whenever behaviour
changes. Follow the general contribution guidelines in
[`docs/guides/developer_guide.md`](guides/developer_guide.md):

1. Make concise commits and run the project's test suite.
2. Cross‑link related guides to keep information discoverable.
3. Run `mkdocs serve` to preview your changes before submitting a pull request.

Issues and improvement ideas are tracked in the main repository issue tracker.
