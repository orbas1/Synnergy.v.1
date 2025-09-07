package cli

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			asset, _ := cmd.Flags().GetString("asset")
			owner, _ := cmd.Flags().GetString("owner")
			shares, _ := cmd.Flags().GetUint64("shares")
			expStr, _ := cmd.Flags().GetString("expiry")
			if asset == "" || owner == "" || shares == 0 {
				return errors.New("asset, owner and shares are required")
			}
			expiry, err := time.Parse(time.RFC3339, expStr)
			if err != nil {
				return err
			}
			tok, err := investorRegistry.Issue(asset, owner, shares, expiry)
			if err != nil {
				return err
			}
			fmt.Println(tok.ID)
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return investorRegistry.Transfer(args[0], args[1])
		},
	}
	cmd.AddCommand(transferCmd)

	returnCmd := &cobra.Command{
		Use:   "return <id> <amount>",
		Short: "Record a return payment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var amt uint64
			if _, err := fmt.Sscanf(args[1], "%d", &amt); err != nil {
				return err
			}
			return investorRegistry.RecordReturn(args[0], amt)
		},
	}
	cmd.AddCommand(returnCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get token info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tok, ok := investorRegistry.Get(args[0])
			if !ok {
				return errors.New("token not found")
			}
			fmt.Printf("ID:%s Asset:%s Owner:%s Shares:%d Active:%t\n", tok.ID, tok.Asset, tok.Owner, tok.Shares, tok.Active)
			for _, r := range tok.Returns {
				fmt.Printf("return %d %s\n", r.Amount, r.Time.Format(time.RFC3339))
			}
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List investor tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			tokens := investorRegistry.List()
			for _, tok := range tokens {
				fmt.Printf("%s %s %s %d\n", tok.ID, tok.Asset, tok.Owner, tok.Shares)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	deactivateCmd := &cobra.Command{
		Use:   "deactivate <id>",
		Short: "Deactivate an investor token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return investorRegistry.Deactivate(args[0])
		},
	}
	cmd.AddCommand(deactivateCmd)

	rootCmd.AddCommand(cmd)
}
