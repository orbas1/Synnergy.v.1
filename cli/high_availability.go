package cli

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"synnergy/core"
	security "synnergy/internal/security"
)

var (
	failoverMu        sync.Mutex
	failover          *core.FailoverManager
	initPrimaryRegion string
	initPrimaryRole   string
	initPrimaryPub    string
	addRegion         string
	addRole           string
	addPub            string
	hbLatency         time.Duration
	hbPayload         string
	hbSignature       string
)

func init() {
	haCmd := &cobra.Command{
		Use:   "highavailability",
		Short: "Manage Stage 77 failover orchestration",
		Long:  "Operate the Stage 77 failover manager, including secure node registration, signed heartbeats and resilience diagnostics for CLI and web dashboards.",
	}

	initCmd := &cobra.Command{
		Use:   "init [primary] [timeoutSec]",
		Args:  cobra.ExactArgs(2),
		Short: "Initialise failover manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			timeoutSec, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid timeout: %w", err)
			}
			region := initPrimaryRegion
			if region == "" {
				region = "global"
			}
			role := initPrimaryRole
			if role == "" {
				role = "orchestrator"
			}
			var pub []byte
			if initPrimaryPub != "" {
				decoded, err := base64.StdEncoding.DecodeString(initPrimaryPub)
				if err != nil {
					return fmt.Errorf("invalid primary public key: %w", err)
				}
				pub = decoded
			}

			vm := core.NewSimpleVM(core.VMHeavy)
			_ = vm.Start()
			wallet, err := core.NewWallet()
			if err != nil {
				return fmt.Errorf("wallet init: %w", err)
			}
			manager := core.NewFailoverManager(args[0], time.Duration(timeoutSec)*time.Second,
				core.WithFailoverVirtualMachine(vm),
				core.WithFailoverConsensus(core.NewConsensusNetworkManager()),
				core.WithFailoverWallet(wallet),
				core.WithFailoverRegistry(core.NewAuthorityNodeRegistry()),
				core.WithFailoverLedger(core.NewLedger()),
				core.WithFailoverSigner(security.NewKeyManager()),
			)
			manager.RegisterNode(core.FailoverNode{ID: args[0], Role: role, Region: region, PublicKey: pub})

			failoverMu.Lock()
			failover = manager
			failoverMu.Unlock()

			gasPrint("Stage77FailoverInit")
			printOutput(map[string]any{
				"status":  "initialised",
				"primary": args[0],
				"timeout": timeoutSec,
				"region":  region,
				"role":    role,
			})
			return nil
		},
	}
	initCmd.Flags().StringVar(&initPrimaryRegion, "primary-region", "global", "Region tag for the primary node")
	initCmd.Flags().StringVar(&initPrimaryRole, "primary-role", "orchestrator", "Role metadata for the primary node")
	initCmd.Flags().StringVar(&initPrimaryPub, "primary-pubkey", "", "Base64 encoded secp256r1 public key for signature verification")

	addCmd := &cobra.Command{
		Use:   "add [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Register a backup node",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := ensureFailoverManager(cmd.Context())
			if err != nil {
				return err
			}
			region := addRegion
			if region == "" {
				region = "global"
			}
			role := addRole
			if role == "" {
				role = "validator"
			}
			var pub []byte
			if addPub != "" {
				decoded, err := base64.StdEncoding.DecodeString(addPub)
				if err != nil {
					return fmt.Errorf("invalid public key: %w", err)
				}
				pub = decoded
			}
			manager.RegisterNode(core.FailoverNode{ID: args[0], Role: role, Region: region, PublicKey: pub})
			gasPrint("Stage77FailoverRegister")
			printOutput(map[string]any{
				"status": "registered",
				"id":     args[0],
				"region": region,
				"role":   role,
				"secure": len(pub) > 0,
			})
			return nil
		},
	}
	addCmd.Flags().StringVar(&addRegion, "region", "global", "Region tag for the node")
	addCmd.Flags().StringVar(&addRole, "role", "validator", "Node role for governance analytics")
	addCmd.Flags().StringVar(&addPub, "pubkey", "", "Base64 encoded secp256r1 public key for heartbeats")

	heartbeatCmd := &cobra.Command{
		Use:   "heartbeat [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Record a heartbeat",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := ensureFailoverManager(cmd.Context())
			if err != nil {
				return err
			}
			var sig []byte
			if hbSignature != "" {
				decoded, err := base64.StdEncoding.DecodeString(hbSignature)
				if err != nil {
					return fmt.Errorf("invalid signature: %w", err)
				}
				sig = decoded
			}
			proof := core.HeartbeatProof{
				ID:        args[0],
				Payload:   []byte(hbPayload),
				Signature: sig,
				Latency:   hbLatency,
			}
			if err := manager.RecordHeartbeat(proof); err != nil {
				return err
			}
			gasPrint("Stage77FailoverHeartbeat")
			printOutput(map[string]any{
				"status":   "heartbeat",
				"id":       args[0],
				"verified": len(sig) > 0,
			})
			return nil
		},
	}
	heartbeatCmd.Flags().DurationVar(&hbLatency, "latency", 0, "Observed latency (e.g. 20ms)")
	heartbeatCmd.Flags().StringVar(&hbPayload, "payload", "", "Payload signed by the node wallet")
	heartbeatCmd.Flags().StringVar(&hbSignature, "signature", "", "Base64 encoded signature of the payload")

	activeCmd := &cobra.Command{
		Use:   "active",
		Short: "Show active node",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := ensureFailoverManager(cmd.Context())
			if err != nil {
				return err
			}
			report := manager.Report(cmd.Context())
			gasPrint("Stage77FailoverActive")
			printOutput(map[string]any{
				"active": report.ActiveNode,
				"role":   report.ActiveRole,
				"region": report.ActiveRegion,
			})
			return nil
		},
	}

	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Emit resilience diagnostics",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager, err := ensureFailoverManager(cmd.Context())
			if err != nil {
				return err
			}
			diag := manager.Report(cmd.Context())
			gasPrint("Stage77FailoverReport")
			printOutput(diag)
			return nil
		},
	}

	haCmd.AddCommand(initCmd, addCmd, heartbeatCmd, activeCmd, reportCmd)
	rootCmd.AddCommand(haCmd)
}

func ensureFailoverManager(ctx context.Context) (*core.FailoverManager, error) {
	failoverMu.Lock()
	defer failoverMu.Unlock()
	if failover != nil {
		return failover, nil
	}
	vm := core.NewSimpleVM(core.VMLight)
	_ = vm.Start()
	wallet, err := core.NewWallet()
	if err != nil {
		return nil, fmt.Errorf("wallet init: %w", err)
	}
	manager := core.NewFailoverManager(wallet.Address, 5*time.Second,
		core.WithFailoverVirtualMachine(vm),
		core.WithFailoverConsensus(core.NewConsensusNetworkManager()),
		core.WithFailoverWallet(wallet),
		core.WithFailoverRegistry(core.NewAuthorityNodeRegistry()),
		core.WithFailoverLedger(core.NewLedger()),
		core.WithFailoverSigner(security.NewKeyManager()),
	)
	manager.RegisterNode(core.FailoverNode{
		ID:        wallet.Address,
		Role:      "orchestrator",
		Region:    "global",
		PublicKey: wallet.PublicKeyBytes(),
	})
	failover = manager
	return failover, nil
}
