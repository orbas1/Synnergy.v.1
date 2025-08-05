package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var investorReg = tokens.NewInvestorRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn2600",
		Short: "Investor token registry",
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue investor token",
		Run: func(cmd *cobra.Command, args []string) {
			asset, _ := cmd.Flags().GetString("asset")
			owner, _ := cmd.Flags().GetString("owner")
			shares, _ := cmd.Flags().GetUint64("shares")
			expiryStr, _ := cmd.Flags().GetString("expiry")
			expiry, _ := time.Parse(time.RFC3339, expiryStr)
			tok := investorReg.Issue(asset, owner, shares, expiry)
			fmt.Println(tok.ID)
		},
	}
	issueCmd.Flags().String("asset", "", "underlying asset")
	issueCmd.Flags().String("owner", "", "token owner")
	issueCmd.Flags().Uint64("shares", 0, "share amount")
	issueCmd.Flags().String("expiry", time.Now().Add(24*time.Hour).Format(time.RFC3339), "expiry time")
	cmd.AddCommand(issueCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <tokenID> <newOwner>",
		Short: "Transfer token ownership",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := investorReg.Transfer(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("transferred")
			}
		},
	}
	cmd.AddCommand(transferCmd)

	returnCmd := &cobra.Command{
		Use:   "record-return <tokenID> <amount>",
		Short: "Record investor return",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := investorReg.RecordReturn(args[0], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("recorded")
			}
		},
	}
	cmd.AddCommand(returnCmd)

	getCmd := &cobra.Command{
		Use:   "get <tokenID>",
		Short: "Get token info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tok, ok := investorReg.Get(args[0])
			if !ok {
				fmt.Println("token not found")
				return
			}
			fmt.Printf("%+v\n", *tok)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List tokens",
		Run: func(cmd *cobra.Command, args []string) {
			for _, tok := range investorReg.List() {
				fmt.Printf("%+v\n", *tok)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
