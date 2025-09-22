package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"synnergy/core"
	"synnergy/internal/nodes"
)

var (
	enterpriseSpecialNode = core.NewEnterpriseSpecialNode(nodes.Address("enterprise-special"))
)

func init() {
	_ = enterpriseSpecialNode.Start()

	enterpriseCmd := &cobra.Command{
		Use:   "enterprise-special",
		Short: "Manage the enterprise special combined node",
	}

	enterpriseCmd.AddCommand(newEnterpriseAttachCmd())
	enterpriseCmd.AddCommand(newEnterpriseDetachCmd())
	enterpriseCmd.AddCommand(newEnterpriseStatusCmd())
	enterpriseCmd.AddCommand(newEnterpriseEventsCmd())
	enterpriseCmd.AddCommand(newEnterpriseBroadcastCmd())
	enterpriseCmd.AddCommand(newEnterpriseLabelsCmd())

	rootCmd.AddCommand(enterpriseCmd)
}

func newEnterpriseAttachCmd() *cobra.Command {
	var (
		id       string
		roleText string
		labels   map[string]string
		seed     uint64
		external bool
	)

	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach a plugin to the enterprise special node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if id == "" {
				return fmt.Errorf("id required")
			}
			role, err := core.ParseCombinedNodeRole(roleText)
			if err != nil {
				return err
			}
			var plugin core.EnterpriseNodePlugin
			if external {
				plugin = core.EnterpriseNodePlugin{ID: id, Role: role, Labels: labels}
			} else {
				ledger := core.NewLedger()
				if seed > 0 {
					ledger.Mint(id, seed)
				}
				node := core.NewNode(id, id, ledger)
				plugin = core.EnterpriseNodePluginFromNode(id, role, node, labels)
			}
			if err := enterpriseSpecialNode.AttachPlugin(plugin); err != nil {
				return err
			}
			cmd.Printf("attached %s as %s\n", id, role)
			return nil
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "Identifier for the plugin")
	cmd.Flags().StringVar(&roleText, "role", string(core.CombinedRoleConsensus), "Plugin role: consensus|execution|analytics|archive")
	cmd.Flags().StringToStringVar(&labels, "label", nil, "Metadata label in key=value form (repeatable)")
	cmd.Flags().Uint64Var(&seed, "seed-balance", 0, "Optional balance minted to the plugin ledger")
	cmd.Flags().BoolVar(&external, "external", false, "Register metadata without creating a synthetic node")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func newEnterpriseDetachCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "detach [id]",
		Short: "Detach a plugin from the enterprise special node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !enterpriseSpecialNode.DetachPlugin(args[0]) {
				return fmt.Errorf("plugin %s not found", args[0])
			}
			cmd.Printf("detached %s\n", args[0])
			return nil
		},
	}
}

func newEnterpriseStatusCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display the combined node snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			snap := enterpriseSpecialNode.Snapshot()
			if asJSON {
				data, err := json.MarshalIndent(snap, "", "  ")
				if err != nil {
					return err
				}
				cmd.Println(string(data))
				return nil
			}
			cmd.Printf("Timestamp: %s\n", snap.Timestamp.Format(time.RFC3339))
			cmd.Printf("Plugins: %d\n", snap.NodeCount)
			cmd.Printf("Total mempool: %d\n", snap.TotalMempool)
			cmd.Printf("Total validators: %d\n", snap.TotalValidators)
			cmd.Printf("Highest block height: %d\n", snap.HighestBlockHeight)
			if len(snap.Roles) > 0 {
				var keys []string
				for role, count := range snap.Roles {
					keys = append(keys, fmt.Sprintf("%s=%d", role, count))
				}
				sort.Strings(keys)
				cmd.Printf("Roles: %s\n", strings.Join(keys, ", "))
			}
			if len(snap.Plugins) > 0 {
				cmd.Println("Plugins:")
				for _, plugin := range snap.Plugins {
					labelPairs := make([]string, 0, len(plugin.Labels))
					for k, v := range plugin.Labels {
						labelPairs = append(labelPairs, fmt.Sprintf("%s=%s", k, v))
					}
					sort.Strings(labelPairs)
					cmd.Printf("  - %s (%s) mempool=%d validators=%d height=%d",
						plugin.ID,
						plugin.Role,
						plugin.Metrics.MempoolSize,
						plugin.Metrics.ValidatorCount,
						plugin.Metrics.BlockHeight,
					)
					if len(labelPairs) > 0 {
						cmd.Printf(" labels=%s", strings.Join(labelPairs, ","))
					}
					cmd.Println()
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Emit JSON output")
	return cmd
}

func newEnterpriseEventsCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Print enterprise special node events",
		Run: func(cmd *cobra.Command, args []string) {
			events := enterpriseSpecialNode.Events()
			if limit > 0 && len(events) > limit {
				events = events[len(events)-limit:]
			}
			for _, evt := range events {
				cmd.Printf("%06d %s %s %v\n", evt.Sequence, evt.Timestamp.Format(time.RFC3339), evt.PluginID, evt.Details)
			}
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 10, "Number of events to display")
	return cmd
}

func newEnterpriseBroadcastCmd() *cobra.Command {
	var (
		from   string
		to     string
		amount uint64
		fee    uint64
		nonce  uint64
	)
	cmd := &cobra.Command{
		Use:   "broadcast",
		Short: "Broadcast a transaction to all attached plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			if from == "" || to == "" {
				return fmt.Errorf("from and to addresses are required")
			}
			tx := core.NewTransaction(from, to, amount, fee, nonce)
			results, err := enterpriseSpecialNode.BroadcastTransaction(tx)
			if err != nil {
				return err
			}
			if len(results) == 0 {
				cmd.Println("no plugins registered with broadcast capability")
				return nil
			}
			for id, resErr := range results {
				if resErr != nil {
					cmd.Printf("%s: %v\n", id, resErr)
				} else {
					cmd.Printf("%s: queued\n", id)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "Sender address")
	cmd.Flags().StringVar(&to, "to", "", "Recipient address")
	cmd.Flags().Uint64Var(&amount, "amount", 0, "Transaction amount")
	cmd.Flags().Uint64Var(&fee, "fee", 0, "Transaction fee")
	cmd.Flags().Uint64Var(&nonce, "nonce", uint64(time.Now().UnixNano()), "Transaction nonce")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

func newEnterpriseLabelsCmd() *cobra.Command {
	var labels map[string]string
	cmd := &cobra.Command{
		Use:   "labels [id]",
		Short: "Update plugin labels",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return enterpriseSpecialNode.UpdatePluginLabels(args[0], labels)
		},
	}
	cmd.Flags().StringToStringVar(&labels, "set", nil, "Labels to assign (key=value)")
	return cmd
}
