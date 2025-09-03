# AI Marketplace Architecture

The AI Marketplace GUI is a thin TypeScript layer that invokes the `synnergy ai_contract deploy` command. Users provide a WebAssembly module, model hash, manifest and gas limit; the GUI executes the CLI and parses the returned contract address. The interface avoids key management and relies on the CLI for virtual machine, consensus and wallet integration, allowing AI-enhanced contracts to be deployed from a web context with minimal surface area.
