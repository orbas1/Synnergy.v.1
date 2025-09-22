package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	synn "synnergy"
	"synnergy/core"
)

const grantRegistryModule = "syn3800"

var (
	grantRegistry     = core.NewGrantRegistry()
	grantRegistryOnce sync.Once
	grantRegistryErr  error
	grantEngineOnce   sync.Once
	grantEngineErr    error
	grantOrchestrator *core.GrantOrchestrator
)

func ensureGrantRegistryLoaded() error {
	grantRegistryOnce.Do(func() {
		var snapshot core.GrantRegistrySnapshot
		ok, err := stage73ReadModule(grantRegistryModule, &snapshot)
		if err != nil {
			grantRegistryErr = err
			return
		}
		if ok {
			grantRegistry.Restore(snapshot)
			syncGrantAuthorities()
		}
	})
	return grantRegistryErr
}

func ensureGrantOrchestrator() error {
	if err := ensureGrantRegistryLoaded(); err != nil {
		return err
	}
	grantEngineOnce.Do(func() {
		vm := synn.NewSimpleVM()
		consensus := core.NewSynnergyConsensus()
		authority := core.NewAuthorityNodeRegistry()
		orchestrator, err := core.NewGrantOrchestrator(grantRegistry, vm, consensus, authority)
		if err != nil {
			grantEngineErr = err
			return
		}
		grantOrchestrator = orchestrator
		syncGrantAuthorities()
	})
	return grantEngineErr
}

func syncGrantAuthorities() {
	if grantOrchestrator == nil {
		return
	}
	grants := grantRegistry.ListGrants()
	for _, record := range grants {
		if record == nil {
			continue
		}
		for addr := range record.Authorizers {
			grantOrchestrator.EnsureAuthority(addr, "grant-authoriser")
		}
		if record.RevokedBy != "" {
			grantOrchestrator.EnsureAuthority(record.RevokedBy, "grant-reviewer")
		}
	}
}

func persistGrantRegistry() error {
	if err := ensureGrantRegistryLoaded(); err != nil {
		return err
	}
	snapshot := grantRegistry.Snapshot()
	return stage73WriteModule(grantRegistryModule, snapshot)
}

func parseAuthorizerFlag(values []string) ([]string, error) {
	var result []string
	for _, val := range values {
		if val == "" {
			continue
		}
		parts := strings.SplitN(val, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("authorizer must be path:password pair")
		}
		wallet, err := loadWallet(parts[0], parts[1])
		if err != nil {
			return nil, fmt.Errorf("load authorizer wallet: %w", err)
		}
		result = append(result, wallet.Address)
	}
	return result, nil
}

func init() {
	cmd := &cobra.Command{
		Use:   "syn3800",
		Short: "Manage SYN3800 grant records",
	}

	createCmd := &cobra.Command{
		Use:   "create <beneficiary> <name> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureGrantOrchestrator(); err != nil {
				return err
			}
			gasPrint("GrantCreate")
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			var authorizers []string
			if cmd.Flags().Changed("authorizer") {
				authorizerFlags, _ := cmd.Flags().GetStringSlice("authorizer")
				authorizers, err = parseAuthorizerFlag(authorizerFlags)
				if err != nil {
					return err
				}
			}
			proof, err := core.NewWalletProof(wallet, core.GrantCreateMessage(args[0], args[1], amt))
			if err != nil {
				return err
			}
			record, err := grantOrchestrator.CreateGrant(cmd.Context(), core.GrantCreationRequest{
				Beneficiary: args[0],
				Name:        args[1],
				Amount:      amt,
				Authorizers: authorizers,
				Creator:     proof,
			})
			if err != nil {
				return err
			}
			if err := persistGrantRegistry(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), record.ID)
			return nil
		},
	}
	createCmd.Flags().StringSlice("authorizer", nil, "Wallet path:password pairs to authorise")
	createCmd.Flags().String("wallet", "", "Wallet file for creator authentication")
	createCmd.Flags().String("password", "", "Wallet password")

	releaseCmd := &cobra.Command{
		Use:   "release <id> <amount> [note]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Release funds for a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureGrantOrchestrator(); err != nil {
				return err
			}
			gasPrint("GrantRelease")
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			proof, err := core.NewWalletProof(wallet, core.GrantReleaseMessage(id, amt, note))
			if err != nil {
				return err
			}
			if _, err := grantOrchestrator.ReleaseGrant(cmd.Context(), core.GrantReleaseRequest{
				GrantID: id,
				Amount:  amt,
				Note:    note,
				Proof:   proof,
			}); err != nil {
				return err
			}
			if err := persistGrantRegistry(); err != nil {
				return err
			}
			cmd.Println("released")
			return nil
		},
	}
	releaseCmd.Flags().String("wallet", "", "Wallet file for authorisation")
	releaseCmd.Flags().String("password", "", "Wallet password")

	authorizeCmd := &cobra.Command{
		Use:   "authorize <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Authorise a wallet to manage a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureGrantOrchestrator(); err != nil {
				return err
			}
			gasPrint("GrantAuthorize")
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			proof, err := core.NewWalletProof(wallet, core.GrantAuthorizeMessage(id))
			if err != nil {
				return err
			}
			if _, err := grantOrchestrator.AuthorizeGrant(cmd.Context(), core.GrantAuthorizationRequest{GrantID: id, Proof: proof}); err != nil {
				return err
			}
			if err := persistGrantRegistry(); err != nil {
				return err
			}
			cmd.Println("authorized")
			return nil
		},
	}
	authorizeCmd.Flags().String("wallet", "", "Wallet file to register")
	authorizeCmd.Flags().String("password", "", "Wallet password")

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant details",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureGrantOrchestrator(); err != nil {
				return err
			}
			gasPrint("GetGrant")
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			g, ok := grantOrchestrator.GetGrant(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			b, _ := json.MarshalIndent(g, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureGrantOrchestrator(); err != nil {
				return err
			}
			gasPrint("ListGrants")
			gs := grantOrchestrator.ListGrants()
			b, _ := json.MarshalIndent(gs, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	auditCmd := &cobra.Command{
		Use:   "audit <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant audit trail",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureGrantOrchestrator(); err != nil {
				return err
			}
			gasPrint("GrantAudit")
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			events, err := grantOrchestrator.Audit(id)
			if err != nil {
				return err
			}
			b, _ := json.MarshalIndent(events, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show grant status totals",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureGrantOrchestrator(); err != nil {
				return err
			}
			gasPrint("GrantStatus")
			summary := grantOrchestrator.StatusSummary()
			b, _ := json.MarshalIndent(summary, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	cmd.AddCommand(createCmd, releaseCmd, authorizeCmd, getCmd, listCmd, auditCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
