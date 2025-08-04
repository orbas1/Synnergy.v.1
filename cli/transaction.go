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

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Transaction utilities",
}
var (
	feeCap   uint64
	feeFloor uint64
)

func init() {

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
			case "feeless":
				fb = core.FeeForValidatedTransfer(val, base, rate, tip, true)
			default:
				fmt.Println("unknown type")
				return
			}
			policy := core.FeePolicy{Cap: feeCap, Floor: feeFloor}
			adj, note := policy.Enforce(fb.Total)
			fb.Total = adj
			if note != "" {
				fmt.Println("note:", note)
			}
			dist := core.DistributeFees(fb.Total)
			fmt.Printf("fee breakdown: %+v\nfee distribution: %+v\n", fb, dist)
		},
	}
	feeCmd.Flags().Uint64Var(&feeCap, "cap", 0, "fee cap")
	feeCmd.Flags().Uint64Var(&feeFloor, "floor", 0, "fee floor")

	baseFeeCmd := &cobra.Command{
		Use:   "basefee [adjustment] [fees...]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Calculate base fee from recent block fees",
		Run: func(cmd *cobra.Command, args []string) {
			adj, _ := strconv.ParseFloat(args[0], 64)
			var fees []uint64
			for _, a := range args[1:] {
				v, _ := strconv.ParseUint(a, 10, 64)
				fees = append(fees, v)
			}
			fmt.Println(core.CalculateBaseFee(fees, adj))
		},
	}

	variableFeeCmd := &cobra.Command{
		Use:   "variablefee [gasUnits] [gasPrice]",
		Args:  cobra.ExactArgs(2),
		Short: "Calculate variable fee component",
		Run: func(cmd *cobra.Command, args []string) {
			units, _ := strconv.ParseUint(args[0], 10, 64)
			price, _ := strconv.ParseUint(args[1], 10, 64)
			fmt.Println(core.CalculateVariableFee(units, price))
		},
	}

	txCmd.AddCommand(createCmd, signCmd, verifyCmd, feeCmd, baseFeeCmd, variableFeeCmd)
	rootCmd.AddCommand(txCmd)
}
