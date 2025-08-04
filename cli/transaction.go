package cli

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction utilities",
	}

	createCmd := &cobra.Command{
		Use:   "create [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Create a transaction",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			fmt.Printf("tx: %+v\n", tx)
		},
	}

	signCmd := &cobra.Command{
		Use:   "sign [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Create and sign a transaction with a new wallet",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			w, err := core.NewWallet()
			if err != nil {
				fmt.Println("wallet err:", err)
				return
			}
			sig, err := w.Sign(tx)
			if err != nil {
				fmt.Println("sign err:", err)
				return
			}
			pubBytes := elliptic.Marshal(elliptic.P256(), w.PrivateKey.PublicKey.X, w.PrivateKey.PublicKey.Y)
			fmt.Printf("txID: %s\npublicKey: %x\nsignature: %x\n", tx.ID, pubBytes, sig)
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [from] [to] [amount] [fee] [nonce] [pubhex] [sighex]",
		Args:  cobra.ExactArgs(7),
		Short: "Verify a transaction signature",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			pubBytes, _ := hex.DecodeString(args[5])
			sig, _ := hex.DecodeString(args[6])
			x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
			pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
			fmt.Println("valid:", core.VerifySignature(tx, sig, pub))
		},
	}

	feeCmd := &cobra.Command{
		Use:   "fee [type] [value] [base] [varRate] [tip]",
		Args:  cobra.ExactArgs(5),
		Short: "Estimate transaction fees and distribution",
		Run: func(cmd *cobra.Command, args []string) {
			val, _ := strconv.ParseUint(args[1], 10, 64)
			base, _ := strconv.ParseUint(args[2], 10, 64)
			rate, _ := strconv.ParseUint(args[3], 10, 64)
			tip, _ := strconv.ParseUint(args[4], 10, 64)
			var fb core.FeeBreakdown
			switch args[0] {
			case "transfer":
				fb = core.FeeForTransfer(val, base, rate, tip)
			case "purchase":
				fb = core.FeeForPurchase(val, base, rate, tip)
			case "token":
				fb = core.FeeForTokenUsage(val, base, rate, tip)
			case "contract":
				fb = core.FeeForContract(val, base, rate, tip)
			case "verify":
				fb = core.FeeForWalletVerification(val, base, rate, tip)
			default:
				fmt.Println("unknown type")
				return
			}
			dist := core.DistributeFees(fb.Total)
			fmt.Printf("fee breakdown: %+v\nfee distribution: %+v\n", fb, dist)
		},
	}

	txCmd.AddCommand(createCmd, signCmd, verifyCmd, feeCmd)
	rootCmd.AddCommand(txCmd)
}
