package cli

import (
	"context"
	"encoding/json"
	"fmt"
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
	orchestratorMu   sync.RWMutex
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
			fmt.Fprintf(cmd.OutOrStdout(), "  Consensus networks: %d (relayers=%d)\n", diag.ConsensusNetworks, diag.ConsensusRelayers)
			fmt.Fprintf(cmd.OutOrStdout(), "  Authority nodes: %d\n", diag.AuthorityNodes)
			if len(diag.AuthorityRoles) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "  Authority roles: %v\n", diag.AuthorityRoles)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "  Wallet address: %s (sealed=%t)\n", diag.WalletAddress, diag.WalletSealed)
			if !diag.GasLastSyncedAt.IsZero() {
				fmt.Fprintf(cmd.OutOrStdout(), "  Gas schedule synced: %s\n", diag.GasLastSyncedAt.Format(time.RFC3339))
			}
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

	orchestratorCmd.AddCommand(statusCmd)
	orchestratorCmd.AddCommand(syncCmd)
	orchestratorCmd.AddCommand(&cobra.Command{
		Use:   "bootstrap",
		Short: "Run EnterpriseBootstrap to verify readiness",
		RunE: func(cmd *cobra.Command, args []string) error {
			orch, err := getEnterpriseOrchestrator(cmd.Context())
			if err != nil {
				return err
			}
			diag, err := orch.EnterpriseBootstrap(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Bootstrap complete: wallet=%s sealed=%t consensus=%d authority=%d\n", diag.WalletAddress, diag.WalletSealed, diag.ConsensusNetworks, diag.AuthorityNodes)
			return nil
		},
	})
	rootCmd.AddCommand(orchestratorCmd)
}

func getEnterpriseOrchestrator(ctx context.Context) (*core.EnterpriseOrchestrator, error) {
	orchestratorMu.RLock()
	inst := orchestratorInst
	err := orchestratorErr
	orchestratorMu.RUnlock()
	if inst != nil || err != nil {
		return inst, err
	}
	orchestratorOnce.Do(func() {
		useCtx := ctx
		if useCtx == nil {
			useCtx = context.Background()
		}
		var cancel context.CancelFunc
		useCtx, cancel = context.WithTimeout(useCtx, 5*time.Second)
		defer cancel()
		inst, err := core.NewEnterpriseOrchestrator(useCtx)
		orchestratorMu.Lock()
		orchestratorInst = inst
		orchestratorErr = err
		orchestratorMu.Unlock()
	})
	orchestratorMu.RLock()
	defer orchestratorMu.RUnlock()
	return orchestratorInst, orchestratorErr
}

// InjectEnterpriseOrchestrator allows the main package to provide a pre-warmed
// orchestrator instance so CLI and web integrations reuse the same diagnostics
// pool.
func InjectEnterpriseOrchestrator(inst *core.EnterpriseOrchestrator) {
	orchestratorMu.Lock()
	orchestratorInst = inst
	orchestratorErr = nil
	orchestratorMu.Unlock()
}
