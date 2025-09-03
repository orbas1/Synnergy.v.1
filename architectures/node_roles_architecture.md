# Node Roles Architecture

Various modules define specialized node behaviors and responsibilities across the network.

**Key Modules**
- authority_nodes.go
- authority_node_index.go
- authority_apply.go
- bank_institutional_node.go
- bank_nodes_index.go
- mining_node.go
- mobile_mining_node.go
- regulatory_node.go
- watchtower_node.go
- warfare_node.go
- environmental_monitoring_node.go
- energy_efficient_node.go
- geospatial_node.go
- indexing_node.go
- central_banking_node.go

**Related CLI Files**
- cli/authority_nodes.go
- cli/authority_node_index.go
- cli/authority_apply.go
- cli/bank_institutional_node.go
- cli/bank_nodes_index.go
- cli/mining_node.go
- cli/mobile_mining_node.go
- cli/regulatory_node.go
- cli/watchtower_node.go
- cli/geospatial.go
- cli/centralbank.go
- cli/warfare_node.go

Together they orchestrate different node types to handle governance, institutional banking, mining, regulation, monitoring, and geographic services. JSON output flags on the CLI facilitate integration with web interfaces and automation.
