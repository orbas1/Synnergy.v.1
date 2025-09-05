# How to use the CLI

The `synnergy` binary exposes network and contract functionality.
Stage 39 extends the toolset with authority governance and banking commands alongside liquidity pool utilities used by the DEX screener.
Stage 40 adds monetary policy utilities; `synnergy coin` validates inputs and can emit JSON for supply and reward queries.

```bash
# Manage authority applications and node membership
synnergy authority_apply submit node1 validator "candidate node"
synnergy authority_apply list --json
synnergy authority register node1 validator
synnergy authority is node1

# Register participating institutions and inspect supported node types
synnergy bankinst register MyBank
synnergy bankinst list
synnergy banknodes types

# Create a new liquidity pool with default fee
synnergy liquidity_pools create TOKENA TOKENB

# Add liquidity and then inspect all pools
synnergy liquidity_pools add TOKENA-TOKENB provider 100 100
synnergy liquidity_views list

# Query a specific pool
synnergy liquidity_views info TOKENA-TOKENB

# Inspect coin parameters
synnergy coin --json info
```

Most commands accept the `--json` flag to produce machine readable output for GUIs.

For comprehensive regression coverage the Stage 46 network harness exercises
CLI commands against live wallet services and in-memory nodes, ensuring end-to-end
flows operate consistently across releases.
## Infrastructure Automation
Stage 50 introduces Terraform and Ansible support for deploying nodes:
```bash
cd deploy/terraform
terraform init
terraform apply -var 'ami_id=ami-123456'

ansible-playbook -i inventory ../ansible/playbook.yml
```
