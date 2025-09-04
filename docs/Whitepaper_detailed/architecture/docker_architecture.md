# Deployment and Container Architecture

The Synnergy network can be packaged into Docker containers for reproducible deployments. Multi-stage builds compile Go binaries for the node and wallet server and produce minimal Alpine-based runtime images.

## Components
- **synnergy**: runs the core blockchain node.
- **walletserver**: provides REST endpoints for GUI and CLI wallets.

Compose files orchestrate the services, mount configuration files and define port mappings. The architecture supports scaling each service independently and serves as the foundation for higher level orchestration platforms.
