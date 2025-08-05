package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn500Token *core.SYN500Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn500",
		Short: "SYN500 utility token",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a SYN500 token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			dec, _ := cmd.Flags().GetUint("dec")
			supply, _ := cmd.Flags().GetUint64("supply")
			syn500Token = core.NewSYN500Token(name, symbol, owner, uint8(dec), supply)
			fmt.Println("token created")
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().Uint("dec", 0, "decimals")
	createCmd.Flags().Uint64("supply", 0, "initial supply")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("owner")

	grantCmd := &cobra.Command{
		Use:   "grant <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Grant a usage tier",
		Run: func(cmd *cobra.Command, args []string) {
			if syn500Token == nil {
				fmt.Println("token not created")
				return
			}
			tier, _ := cmd.Flags().GetInt("tier")
			max, _ := cmd.Flags().GetUint64("max")
			syn500Token.Grant(args[0], tier, max)
			fmt.Println("granted")
		},
	}
	grantCmd.Flags().Int("tier", 0, "service tier")
	grantCmd.Flags().Uint64("max", 0, "max usage")
	grantCmd.MarkFlagRequired("tier")
	grantCmd.MarkFlagRequired("max")

	useCmd := &cobra.Command{
		Use:   "use <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Record usage",
		Run: func(cmd *cobra.Command, args []string) {
			if syn500Token == nil {
				fmt.Println("token not created")
				return
			}
			if err := syn500Token.Use(args[0]); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("usage recorded")
			}
		},
	}

	cmd.AddCommand(createCmd, grantCmd, useCmd)
	rootCmd.AddCommand(cmd)
}
