# Opcodes and gas

Stage 39 introduces opcodes for querying liquidity pools. `Liquidity_Pool`
returns a single pool view while `Liquidity_Pools` lists all pools. Each opcode
has a default gas cost of 1 as recorded in `gas_table_list.md`.

To guard against regressions, the Stage 46 network harness captures a JSON
snapshot of the gas table during end-to-end tests so that opcode pricing can be
validated by external dashboards and wallets.
