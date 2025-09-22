package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synn "synnergy"
	"synnergy/core"
)

var assetRegistry = core.NewAssetRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn800_token",
		Short: "Manage SYN800 asset tokens",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <desc> <valuation> <loc> <type> <cert>",
		Args:  cobra.ExactArgs(6),
		Short: "Register an asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid valuation")
			}
			asset, err := assetRegistry.Register(args[0], args[1], val, args[3], args[4], args[5])
			if err != nil {
				return err
			}
			printOutput(map[string]any{
				"status": "asset registered",
				"id":     asset.ID,
				"value":  asset.Valuation,
				"gas":    synn.GasCost("RegisterTangibleAsset"),
			})
			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update <id> <valuation>",
		Args:  cobra.ExactArgs(2),
		Short: "Update asset valuation",
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid valuation")
			}
			if err := assetRegistry.UpdateValuation(args[0], val); err != nil {
				return err
			}
			printOutput(map[string]any{
				"status": "valuation updated",
				"id":     args[0],
				"value":  val,
				"gas":    synn.GasCost("UpdateAssetValuation"),
			})
			return nil
		},
	}

	custodianCmd := &cobra.Command{
		Use:   "custodian <id> <address>",
		Args:  cobra.ExactArgs(2),
		Short: "Assign asset custodian",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := assetRegistry.AssignCustodian(args[0], args[1]); err != nil {
				return err
			}
			printOutput(map[string]any{
				"status":    "custodian assigned",
				"id":        args[0],
				"custodian": args[1],
				"gas":       synn.GasCost("AssignAssetCustodian"),
			})
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Display asset information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if a, ok := assetRegistry.Get(args[0]); ok {
				snap := map[string]any{}
				b, _ := json.Marshal(a)
				_ = json.Unmarshal(b, &snap)
				snap["gas"] = synn.GasCost("DescribeTangibleAsset")
				printOutput(snap)
				return nil
			}
			return fmt.Errorf("asset not found")
		},
	}

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "List assets sorted by valuation",
		RunE: func(cmd *cobra.Command, args []string) error {
			assets := assetRegistry.Snapshot()
			printOutput(map[string]any{
				"assets": assets,
				"gas":    synn.GasCost("AssetSnapshot"),
			})
			return nil
		},
	}

	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Display recent asset updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			limit, _ := cmd.Flags().GetInt("limit")
			entries := assetRegistry.History(limit)
			printOutput(map[string]any{
				"entries": entries,
				"gas":     synn.GasCost("AssetHistory"),
			})
			return nil
		},
	}
	historyCmd.Flags().Int("limit", 5, "number of entries to return")

	cmd.AddCommand(registerCmd, updateCmd, custodianCmd, infoCmd, snapshotCmd, historyCmd)
	rootCmd.AddCommand(cmd)
}
