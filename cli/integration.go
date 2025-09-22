package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"synnergy/core"
)

var (
	integrationFormat  string
	integrationTimeout time.Duration
)

func newIntegrationStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Run enterprise integration diagnostics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if integrationTimeout <= 0 {
				integrationTimeout = 5 * time.Second
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), integrationTimeout)
			defer cancel()

			integration, err := core.NewPlatformIntegration()
			if err != nil {
				return err
			}
			defer func() {
				if cerr := integration.Close(); cerr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", cerr)
				}
			}()

			status := integration.Diagnostics(ctx)
			switch strings.ToLower(integrationFormat) {
			case "json":
				b, err := json.MarshalIndent(status, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
			case "table":
				fmt.Fprintf(cmd.OutOrStdout(), "Integration Health (%s)\n", status.Timestamp.Format(time.RFC3339))
				fmt.Fprintf(cmd.OutOrStdout(), "  VM: %s concurrency=%d\n", status.VM.Mode, status.VM.Concurrency)
				fmt.Fprintf(cmd.OutOrStdout(), "  Node: %s height=%d pending=%d\n", status.Node.ID, status.Node.BlockHeight, status.Node.PendingTx)
				fmt.Fprintf(cmd.OutOrStdout(), "  Wallet: %s\n", status.Wallet.Address)
				fmt.Fprintf(cmd.OutOrStdout(), "  Consensus networks: %d relayers=%d\n", status.Consensus.Networks, status.Consensus.AuthorizedRelays)
				fmt.Fprintf(cmd.OutOrStdout(), "  Authority nodes: %d\n", status.Authority.Registered)
				fmt.Fprintln(cmd.OutOrStdout(), "  Enterprise readiness:")
				enterprise := map[string]core.DiagnosticCheck{
					"    Security":         status.Enterprise.Security,
					"    Scalability":      status.Enterprise.Scalability,
					"    Privacy":          status.Enterprise.Privacy,
					"    Governance":       status.Enterprise.Governance,
					"    Interoperability": status.Enterprise.Interoperability,
					"    Compliance":       status.Enterprise.Compliance,
				}
				for label, check := range enterprise {
					fmt.Fprintf(cmd.OutOrStdout(), "%s → %s (healthy=%t latency=%s)\n", label, check.Detail, check.Healthy, check.Latency)
				}
				for name, diag := range status.Diagnostics {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s: %s\n", name, diag)
				}
			default:
				return fmt.Errorf("unsupported format: %s", integrationFormat)
			}

			if len(status.Issues) > 0 {
				return fmt.Errorf("integration reported issues: %s", strings.Join(status.Issues, "; "))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&integrationFormat, "format", "json", "Output format: json or table")
	cmd.Flags().DurationVar(&integrationTimeout, "timeout", 5*time.Second, "Diagnostics timeout")
	return cmd
}

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Inspect enterprise subsystem integration",
}

func init() {
	integrationCmd.AddCommand(newIntegrationStatusCommand())
	rootCmd.AddCommand(integrationCmd)
}
