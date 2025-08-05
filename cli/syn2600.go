package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var investorRegistry = tokens.NewInvestorRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn2600",
		Short: "Investor token registry",
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a new investor token",
		Run: func(cmd *cobra.Command, args []string) {
			asset, _ := cmd.Flags().GetString("asset")
			owner, _ := cmd.Flags().GetString("owner")
			shares, _ := cmd.Flags().GetUint64("shares")
			expStr, _ := cmd.Flags().GetString("expiry")
			expiry, _ := time.Parse(time.RFC3339, expStr)
			tok := investorRegistry.Issue(asset, owner, shares, expiry)
			fmt.Println(tok.ID)
		},
	}
	issueCmd.Flags().String("asset", "", "underlying asset")
	issueCmd.Flags().String("owner", "", "owner address")
	issueCmd.Flags().Uint64("shares", 0, "share quantity")
	issueCmd.Flags().String("expiry", time.Now().Add(24*time.Hour).Format(time.RFC3339), "expiry time")
	cmd.AddCommand(issueCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <newOwner>",
		Short: "Transfer token ownership",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := investorRegistry.Transfer(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(transferCmd)

	returnCmd := &cobra.Command{
		Use:   "return <id> <amount>",
		Short: "Record a return payment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := investorRegistry.RecordReturn(args[0], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(returnCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get token info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tok, ok := investorRegistry.Get(args[0])
			if !ok {
				fmt.Println("token not found")
				return
			}
			fmt.Printf("ID:%s Asset:%s Owner:%s Shares:%d Active:%t\n", tok.ID, tok.Asset, tok.Owner, tok.Shares, tok.Active)
			for _, r := range tok.Returns {
				fmt.Printf("return %d %s\n", r.Amount, r.Time.Format(time.RFC3339))
			}
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List investor tokens",
		Run: func(cmd *cobra.Command, args []string) {
			tokens := investorRegistry.List()
			for _, tok := range tokens {
				fmt.Printf("%s %s %s %d\n", tok.ID, tok.Asset, tok.Owner, tok.Shares)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
