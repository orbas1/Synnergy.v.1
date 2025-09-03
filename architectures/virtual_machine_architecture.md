# Virtual Machine Architecture

Stage 11 formalises the Synnergy VM (SVM) as a modular execution engine. Each VM instance enforces resource limits through an internal sandbox manager and supports context-aware timeouts for contract calls. Sandboxes can be deleted once processing completes, ensuring that long-running contracts do not leak memory or gas accounting state.
