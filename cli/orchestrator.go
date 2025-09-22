package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"synnergy/core"
)

var (
	orchestratorOnce sync.Once
	orchestratorInst *core.EnterpriseOrchestrator
	orchestratorErr  error
	orchestratorJSON bool
)

func init() {
	orchestratorCmd := &cobra.Command{
		Use:   "orchestrator",
		Short: "Enterprise orchestrator utilities",
		Long: "Coordinate Stage 78 subsystems including the VM, consensus network, " +
			"wallet, node registry and gas schedule so operators can verify readiness from the CLI.",
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Report orchestrator diagnostics",
		RunE: func(cmd *cobra.Command, args []string) error {
			orch, err := getEnterpriseOrchestrator(cmd.Context())
			if err != nil {
				return err
			}
			diag := orch.Diagnostics(cmd.Context())
			if orchestratorJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(diag)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Orchestrator status at %s\n", diag.Timestamp.Format(time.RFC3339))
			fmt.Fprintf(cmd.OutOrStdout(), "  VM: %s (running=%t, concurrency=%d)\n", diag.VMMode, diag.VMRunning, diag.VMConcurrency)
			fmt.Fprintf(cmd.OutOrStdout(), "  Consensus networks: %d\n", diag.ConsensusNetworks)
			fmt.Fprintf(cmd.OutOrStdout(), "  Authority nodes: %d\n", diag.AuthorityNodes)
			fmt.Fprintf(cmd.OutOrStdout(), "  Bootstrap nodes: %d\n", diag.BootstrapNodes)
			fmt.Fprintf(cmd.OutOrStdout(), "  Replication active: %t\n", diag.ReplicationActive)
			fmt.Fprintf(cmd.OutOrStdout(), "  Wallet address: %s\n", diag.WalletAddress)
			if len(diag.InsertedOpcodes) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  Newly documented opcodes: %v\n", diag.InsertedOpcodes)
			}
			if len(diag.MissingOpcodes) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  Missing opcode documentation: %v\n", diag.MissingOpcodes)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "  Opcode documentation: complete")
			}
			return nil
		},
	}
	statusCmd.Flags().BoolVar(&orchestratorJSON, "json", false, "Emit diagnostics as JSON")

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Refresh gas schedule and emit diagnostics",
		RunE: func(cmd *cobra.Command, args []string) error {
			orch, err := getEnterpriseOrchestrator(cmd.Context())
			if err != nil {
				return err
			}
			diag, err := orch.SyncGasSchedule(cmd.Context(), nil)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Synced %d opcodes; %d registered authority nodes\n", len(diag.GasCoverage), diag.AuthorityNodes)
			if len(diag.MissingOpcodes) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Missing opcode documentation: %v\n", diag.MissingOpcodes)
			}
			return nil
		},
	}

	var (
		bootstrapNodeID      string
		bootstrapAddress     string
		bootstrapConsensus   string
		bootstrapGovernance  string
		bootstrapReplicate   bool
		bootstrapRegulator   bool
		bootstrapAuthorities []string
	)

	bootstrapCmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap an enterprise node",
		Long: "Provision a ledger-backed node, register authorities, attach regulatory checks " +
			"and enable replication so enterprise networks can be stood up deterministically.",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("EnterpriseNetworkBootstrap")
			orch, err := getEnterpriseOrchestrator(cmd.Context())
			if err != nil {
				return err
			}
			cfg := core.EnterpriseBootstrapConfig{
				NodeID:            bootstrapNodeID,
				Address:           bootstrapAddress,
				ConsensusProfile:  bootstrapConsensus,
				GovernanceProfile: bootstrapGovernance,
				EnableReplication: bootstrapReplicate,
				EnableRegulator:   bootstrapRegulator,
				Authorities:       make(map[string]string),
			}
			for _, entry := range bootstrapAuthorities {
				parts := strings.SplitN(entry, "=", 2)
				addr := strings.TrimSpace(parts[0])
				if addr == "" {
					continue
				}
				role := "authority"
				if len(parts) == 2 {
					if trimmed := strings.TrimSpace(parts[1]); trimmed != "" {
						role = trimmed
					}
				}
				cfg.Authorities[addr] = role
			}
			if cfg.NodeID == "" {
				cfg.NodeID = "enterprise-node"
			}
			res, err := orch.BootstrapNetwork(cmd.Context(), cfg)
			if err != nil {
				return err
			}
			printOutput(res)
			return nil
		},
	}
	bootstrapCmd.Flags().StringVar(&bootstrapNodeID, "node-id", "enterprise-node", "Identifier for the bootstrap node")
	bootstrapCmd.Flags().StringVar(&bootstrapAddress, "address", "", "Network address for the node")
	bootstrapCmd.Flags().StringVar(&bootstrapConsensus, "consensus", "Synnergy-PBFT", "Consensus profile to join")
	bootstrapCmd.Flags().StringVar(&bootstrapGovernance, "governance", "SYN-Gov", "Governance profile to align with")
	bootstrapCmd.Flags().BoolVar(&bootstrapReplicate, "replicate", true, "Enable ledger replication services")
	bootstrapCmd.Flags().BoolVar(&bootstrapRegulator, "regulator", true, "Attach regulatory checks to consensus")
	bootstrapCmd.Flags().StringSliceVar(&bootstrapAuthorities, "authority", nil, "Additional authority entries as address=role")

	orchestratorCmd.AddCommand(statusCmd)
	orchestratorCmd.AddCommand(syncCmd)
	orchestratorCmd.AddCommand(bootstrapCmd)
	rootCmd.AddCommand(orchestratorCmd)
}

func getEnterpriseOrchestrator(ctx context.Context) (*core.EnterpriseOrchestrator, error) {
	orchestratorOnce.Do(func() {
		useCtx := ctx
		if useCtx == nil {
			useCtx = context.Background()
		}
		var cancel context.CancelFunc
		useCtx, cancel = context.WithTimeout(useCtx, 5*time.Second)
		defer cancel()
		orchestratorInst, orchestratorErr = core.NewEnterpriseOrchestrator(useCtx)
	})
	return orchestratorInst, orchestratorErr
}
