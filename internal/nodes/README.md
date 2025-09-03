# Nodes

This package collects lightweight node implementations used across the project.
It exposes a small `Node` interface and a reference `BasicNode` implementation
that offers a concurrency‑safe lifecycle and peer management layer.

Subpackages such as `watchtower`, `military_nodes` and `optimization_nodes`
demonstrate how specialised roles can build upon `BasicNode` without depending
on the larger `core` packages. Additional node types can embed `BasicNode` and
extend it with domain specific behaviour.
