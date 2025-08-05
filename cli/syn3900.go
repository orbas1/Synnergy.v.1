package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var benefitRegistry = core.NewBenefitRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3900",
		Short: "Manage SYN3900 government benefits",
	}

	registerCmd := &cobra.Command{
		Use:   "register <recipient> <program> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a new benefit",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			id := benefitRegistry.RegisterBenefit(args[0], args[1], amt)
			fmt.Println("benefit registered", id)
		},
	}

	claimCmd := &cobra.Command{
		Use:   "claim <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Claim a benefit",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if err := benefitRegistry.Claim(id); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show benefit details",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if b, ok := benefitRegistry.GetBenefit(id); ok {
				data, _ := json.MarshalIndent(b, "", "  ")
				fmt.Println(string(data))
			} else {
				fmt.Println("not found")
			}
		},
	}

	cmd.AddCommand(registerCmd, claimCmd, getCmd)
	rootCmd.AddCommand(cmd)
}
