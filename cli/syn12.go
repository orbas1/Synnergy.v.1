package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
	"time"
)

var syn12Token *tokens.SYN12Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn12",
		Short: "SYN12 treasury bill token",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise SYN12 token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint32("decimals")
			billID, _ := cmd.Flags().GetString("bill")
			issuer, _ := cmd.Flags().GetString("issuer")
			issueStr, _ := cmd.Flags().GetString("issue")
			maturityStr, _ := cmd.Flags().GetString("maturity")
			discount, _ := cmd.Flags().GetFloat64("discount")
			face, _ := cmd.Flags().GetUint64("face")
			issue, _ := time.Parse(time.RFC3339, issueStr)
			maturity, _ := time.Parse(time.RFC3339, maturityStr)
			meta := tokens.SYN12Metadata{BillID: billID, Issuer: issuer, IssueDate: issue, Maturity: maturity, Discount: discount, FaceValue: face}
			id := tokenRegistry.NextID()
			syn12Token = tokens.NewSYN12Token(id, name, symbol, meta, uint8(dec))
			tokenRegistry.Register(syn12Token)
			fmt.Println("syn12 initialised")
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().Uint32("decimals", 2, "decimal places")
	initCmd.Flags().String("bill", "", "bill identifier")
	initCmd.Flags().String("issuer", "", "issuer")
	initCmd.Flags().String("issue", time.Now().Format(time.RFC3339), "issue date RFC3339")
	initCmd.Flags().String("maturity", time.Now().Add(30*24*time.Hour).Format(time.RFC3339), "maturity date RFC3339")
	initCmd.Flags().Float64("discount", 0, "discount rate")
	initCmd.Flags().Uint64("face", 0, "face value")
	cmd.AddCommand(initCmd)

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show token info",
		Run: func(cmd *cobra.Command, args []string) {
			if syn12Token == nil {
				fmt.Println("token not initialised")
				return
			}
			meta := syn12Token.Metadata
			fmt.Printf("Bill:%s Issuer:%s Issue:%s Maturity:%s Discount:%f Face:%d Supply:%d\n", meta.BillID, meta.Issuer, meta.IssueDate.Format(time.RFC3339), meta.Maturity.Format(time.RFC3339), meta.Discount, meta.FaceValue, syn12Token.TotalSupply())
		},
	}
	cmd.AddCommand(infoCmd)

	mintCmd := &cobra.Command{
		Use:   "mint <to> <amt>",
		Short: "Mint tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn12Token == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := syn12Token.Mint(args[0], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("minted")
		},
	}
	cmd.AddCommand(mintCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amt>",
		Short: "Transfer tokens",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if syn12Token == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := syn12Token.Transfer(args[0], args[1], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("transferred")
		},
	}
	cmd.AddCommand(transferCmd)

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show balance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn12Token == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn12Token.BalanceOf(args[0]))
		},
	}
	cmd.AddCommand(balanceCmd)

	rootCmd.AddCommand(cmd)
}
