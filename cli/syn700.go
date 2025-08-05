package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var ipRegistry = core.NewIPRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn700",
		Short: "Manage SYN700 IP tokens",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <title> <desc> <creator> <owner>",
		Args:  cobra.ExactArgs(5),
		Short: "Register an IP asset",
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := ipRegistry.Register(args[0], args[1], args[2], args[3], args[4]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	licenseCmd := &cobra.Command{
		Use:   "license <tokenID> <licID> <type> <licensee> <royalty>",
		Args:  cobra.ExactArgs(5),
		Short: "Create a license",
		Run: func(cmd *cobra.Command, args []string) {
			royalty, _ := strconv.ParseUint(args[4], 10, 64)
			if err := ipRegistry.CreateLicense(args[0], args[1], args[2], args[3], royalty); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	royaltyCmd := &cobra.Command{
		Use:   "royalty <tokenID> <licID> <licensee> <amount>",
		Args:  cobra.ExactArgs(4),
		Short: "Record a royalty payment",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[3], 10, 64)
			if err := ipRegistry.RecordRoyalty(args[0], args[1], args[2], amt); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <tokenID>",
		Args:  cobra.ExactArgs(1),
		Short: "Show token info",
		Run: func(cmd *cobra.Command, args []string) {
			if t, ok := ipRegistry.Get(args[0]); ok {
				b, _ := json.MarshalIndent(t, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	}

	cmd.AddCommand(registerCmd, licenseCmd, royaltyCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
