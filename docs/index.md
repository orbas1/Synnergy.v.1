# Synnergy Documentation

Synnergy is a modular research blockchain. The `docs/` tree is the primary
reference for contributors and operators and is published as a static web site
with [MkDocs](https://www.mkdocs.org/). The documentation covers design
decisions, end‑user guides and benchmarking data for the core platform.

## Structure

The documentation is split into several key sections:

| Path | Description |
|------|-------------|
| [`guides/`](guides/) | Quickstart and how‑to articles for developers, node operators, and GUI users. |
| [`api/`](api/) | Package‑level references and instructions for generating API documentation. |
| [`adr/`](adr/) | Architecture Decision Records describing significant technical choices made during development. |
| [`reference/`](reference/) | Opcode listings, error catalogues, and other low‑level references. |
| [`ux/`](ux/) | User‑experience guidelines on accessibility, theming, and localization. |
| [`Whitepaper_detailed/`](Whitepaper_detailed/) | Expanded whitepaper chapters covering background and rationale. |
| [`financial_models.md`](financial_models.md) | Economic modelling and incentive design notes. |
| [`performance_benchmarks.md`](performance_benchmarks.md) | Baseline performance figures and instructions for running benchmarks. |

Links to additional material in the repository root are provided throughout the
guides where relevant.

## Building the Site

Configuration for the site lives in [`mkdocs.yml`](../mkdocs.yml). Install
MkDocs and the material theme:

```bash
pip install mkdocs mkdocs-material
```

Serve the documentation locally while editing:

```bash
mkdocs serve
```

or use the Makefile target:

```bash
make docs-serve
```

or generate a static site under `site/`:

```bash
mkdocs build
```

which is also available via:

```bash
make docs
```

## Contributing

Documentation lives alongside the code and should be updated whenever behaviour
changes. Follow the general contribution guidelines in
[`docs/guides/developer_guide.md`](guides/developer_guide.md):

1. Make concise commits and run the project's test suite.
2. Cross‑link related guides to keep information discoverable.
3. Run `mkdocs serve` to preview your changes before submitting a pull request.

Issues and improvement ideas are tracked in the main repository issue tracker.
