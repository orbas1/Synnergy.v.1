# Production Readiness Plan

This document outlines a 20-stage roadmap for reorganizing the repository and preparing it for production-grade deployment. The stages can be executed in parallel by separate teams but are ordered to show overall dependencies.

1. **Establish Module Boundaries** ✅
   - Audit existing packages and identify domain groupings (core, nodes, tokens, cross-chain, security, etc.). *(completed)*
   - Create `internal/` and `pkg/` directories to separate internal code from reusable libraries. *(completed)*

2. **Relocate Domain Packages**  
   - Move node-related code into `internal/nodes/` with subpackages per node type.  
   - Group token implementations under `internal/tokens/`.  
   - Update import paths accordingly.

3. **Create `cmd/` Entrypoints**
   - Maintain a dedicated `cmd/` directory for all binaries (e.g., node, tooling).
   - Each binary receives its own subdirectory with a `main.go` entrypoint.

4. **Introduce Dependency Management**
   - Use Go modules with explicit versioning and keep `go.mod`/`go.sum` under version control.
   - Enforce clean dependency graphs with `go mod tidy` and integrity checks via `go mod verify`.
   - Provide `Makefile` targets (`tidy`, `verify`, `update`) to standardize dependency maintenance.
   - Configure Dependabot in `.github/dependabot.yml` to automatically open pull requests for Go modules and GitHub Actions.
   - Establish a review policy for dependency changes, including security scanning before merge.

5. **Implement Configuration Management**  
   - Centralize configuration logic under `internal/config/`.  
   - Support environment variables and configuration files (YAML/JSON) with validation.

6. **Logging and Instrumentation**
   - Standardize structured logging across the codebase using a high-performance library such as `zap` or `logrus`.
   - Create an `internal/log` package that wraps the chosen logger and exposes helper functions for common patterns.
   - Ensure every package obtains loggers via dependency injection or context and avoid global state.
   - Define log levels (debug, info, warn, error, fatal) and allow them to be configured through configuration files and environment variables.
   - Emit logs in JSON format with timestamps, component names, and key–value pairs to ease parsing by log aggregators.
   - Include correlation IDs and request context in log entries to enable end-to-end tracing of operations.
   - Support multiple log sinks (stdout, rotating files, and remote endpoints like ELK/Loki/Splunk) with pluggable backends.
   - Provide log rotation, retention policies, and size limits for long-running services.
   - Document usage guidelines so new modules follow consistent logging practices.
   - Instrument application metrics using Prometheus, exposing a `/metrics` HTTP endpoint.
   - Capture counters, gauges, and histograms for key operations (transaction throughput, block processing latency, resource utilization, error rates).
   - Export standard Go runtime metrics and add custom collectors where necessary.
   - Supply starter Grafana dashboards and alerting rules for common metrics.
   - Add health and readiness probes backed by metrics to integrate with orchestration platforms like Kubernetes.
   - Write unit tests verifying that the logging facade initializes correctly and that metrics are registered without conflict.

7. **Error Handling and Observability**
   - Create an `internal/errors` package that defines typed errors, error
     codes, and helper functions for wrapping with `%w`.
   - Replace panics with structured error propagation and enforce
     consistent handling using `errors.Is`/`errors.As`.
   - Attach contextual metadata (component, operation, severity) to
     errors so they can be correlated in logs and metrics.
   - Integrate OpenTelemetry across all modules to emit traces, metrics,
     and logs; propagate `context.Context` to carry trace identifiers.
   - Provide reference deployments for OTLP collectors (e.g., Jaeger) and
     dashboard templates in Prometheus/Grafana with alerting on error
     rates and latency.
   - Document error-handling conventions and observability setup for
     contributors and operators.

8. **Security Hardening**  
   - Add static analysis tools (`gosec`, `staticcheck`).  
   - Enforce code scanning in CI and address findings.  
   - Perform dependency vulnerability scans.

9. **Testing Framework**
   - Standardize test layout using table-driven tests and shared `testdata/` fixtures.
   - Achieve at least 80% unit-test coverage across all packages with coverage reports gated in CI.
   - Utilize mocks and fakes (e.g., `testify`, `gomock`) to isolate external dependencies.
   - Add integration tests for cross-chain, token, security, and node workflows.
   - Provide end-to-end tests orchestrated with `docker-compose` to simulate multi-node networks.
   - Incorporate fuzz and property-based testing for critical components.
   - Run race detector (`go test -race`) and `go vet` as part of the test pipeline.
   - Publish test and coverage reports to services like Codecov for visibility.
   - Document testing guidelines and assign ownership for maintaining test quality.
   - Use GitHub Actions to execute the full test suite on each push and pull request.

10. **CI/CD Pipeline (Completed)**
    - Implemented a GitHub Actions workflow that builds, tests, lints, and packages binaries.
    - Enabled caching for modules and test results to speed up builds.

11. **Documentation Standardization**
    - Move guides into a `docs/` directory.
    - Use a documentation generator (e.g., MkDocs) to produce a docs site.
    - Maintain ADRs for architectural decisions.

12. **API and RPC Layer**
    - Define versioned gRPC and REST interfaces for all node and service interactions, ensuring backward compatibility as the protocol evolves.
    - Organize protobuf definitions under `api/` with clear package boundaries and generate client/server stubs using a reproducible toolchain (e.g., `buf` or `protoc`).
    - Generate OpenAPI/Swagger specifications for the REST layer and publish language-specific client SDKs.
    - Expose streaming endpoints (gRPC streaming or WebSockets) for real-time events and data feeds.
    - Authenticate and authorize requests using JWT/OAuth2 and mutual TLS for gRPC; support API keys for service-to-service calls.
    - Enforce request validation, rate limiting, and quotas via an API gateway layer to protect nodes from abuse and DoS attacks.
    - Standardize error codes and response envelopes to provide machine-parseable failures and rich debugging information.
    - Instrument RPC handlers with structured logs, Prometheus metrics, and OpenTelemetry traces for observability.
    - Provide health-check and reflection endpoints for service discovery and automated tooling.
    - Include integration tests validating RPC contracts and backward compatibility guarantees.

13. **Configuration of Build Tags and Environments**
    - Use build tags for optional features (e.g., experimental nodes).
    - Provide separate configs for dev/test/production environments.

14. **Containerization**  
    - Create Dockerfiles for each binary.  
    - Use multi-stage builds for minimal runtime images.  
    - Provide a `docker-compose.yml` for local orchestration.

15. **Release Management**  
    - Adopt semantic versioning.  
    - Automate release notes and changelog generation.  
    - Sign releases and provide checksums.

16. **Performance Benchmarking**  
    - Add Go benchmarks for critical paths.  
    - Establish performance baselines and set budgets.  
    - Monitor regressions in CI.

17. **Persistence and State Management**
    - Introduce an `internal/storage` layer that exposes clean interfaces for CRUD operations and transactions, allowing services to remain database agnostic.
    - Provide adapters for multiple backends (e.g., Postgres via `database/sql`, LevelDB, in-memory) with pluggable drivers and configurable connection pooling, TLS, and authentication.
    - Support ACID transactions and both optimistic and pessimistic locking while honoring context deadlines for timeouts and cancellation.
    - Add migration tooling (e.g., `golang-migrate` or `atlas`) with versioned DDL files, rollback capability, and CI hooks to verify schema consistency.
    - Include schema compatibility tests for each backend and enforce migrations through automated checks before deployment.
    - Implement a caching layer (such as Redis) behind an interface, featuring TTL controls and cache invalidation strategies to keep reads fast and consistent.
    - Document backup/restore, replication, and snapshot procedures, and emit metrics for query latency, cache hit rates, and failure conditions.
    - Provide configuration and operational runbooks for development, testing, and production environments.

18. **Networking and P2P Layer**  
    - Encapsulate networking code under `internal/p2p/`.  
    - Support secure transports (TLS, Noise protocol) and peer discovery.

19. **Governance and Access Control**  
    - Centralize RBAC logic in `internal/auth/`.  
    - Implement policy enforcement and audit logging.

20. **Packaging and Distribution**
    - Offer install scripts and shell-based installers for quick setup on new systems.
    - Maintain package manager integrations including Homebrew (macOS), APT/YUM (Linux), and Scoop/Chocolatey (Windows).
    - Produce signed binaries for Linux, macOS, and Windows across `amd64` and `arm64` architectures with checksum files and GPG signatures.
    - Publish Docker images for every release to public and private registries with SBOM metadata.
    - Generate versioned binary tarballs and OS-specific packages (`.deb`, `.rpm`, `.msi`) via GitHub Releases.
    - Automate packaging and artifact publishing in CI/CD with distinct stable and pre-release channels.
    - Document installation, upgrade, and rollback procedures and verify packages through integration tests.

These stages, once completed, will transition the repository from a prototype into a maintainable, production-grade codebase suitable for enterprise deployment.

