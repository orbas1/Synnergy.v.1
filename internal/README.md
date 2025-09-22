# Internal Packages

The `internal/` tree contains the Synnergy platform's hardened runtime libraries. These packages are the glue between the CLI, node processes, consensus layers, wallet services and the browser UX. They are not intended for consumption by external modules; instead they expose well-tested, enterprise-grade primitives that the rest of the codebase builds upon.

Key highlights of the internal services upgraded in Stage 85 include:

- **API** – Production-ready authentication middleware with signed token validation, rate limiting and RBAC integration. The HTTP gateway now exposes fault-tolerant lifecycle management with graceful shutdown, contextual claim propagation and consistent error handling for CLI and web clients.
- **Auth** – The RBAC engine supports conditional permissions, in-memory and stream audit logging, and explicit role lifecycle management to meet governance and compliance requirements.
- **Registry & Identity services** – Identity verification and wallet registration components now emit structured events for the CLI and web UI, enforce cryptographic integrity and support secure metadata distribution across consensus and VM layers.

Only binaries within this repository should import packages from `internal/`. Application teams can rely on the interfaces here for auditable, replay-safe behaviours while keeping the surface area tightly controlled.
