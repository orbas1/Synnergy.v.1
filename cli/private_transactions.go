package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var privTxMgr = core.NewPrivateTxManager()

func init() {
	cmd := &cobra.Command{
		Use:   "private-tx",
		Short: "Manage private transactions",
	}

	encryptCmd := &cobra.Command{
		Use:   "encrypt [key] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Encrypt transaction payload bytes",
		Run: func(cmd *cobra.Command, args []string) {
			ct, err := core.Encrypt([]byte(args[0]), []byte(args[1]))
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(hex.EncodeToString(ct))
		},
	}

	decryptCmd := &cobra.Command{
		Use:   "decrypt [key] [hexdata]",
		Args:  cobra.ExactArgs(2),
		Short: "Decrypt previously encrypted payload",
		Run: func(cmd *cobra.Command, args []string) {
			data, err := hex.DecodeString(args[1])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			pt, err := core.Decrypt([]byte(args[0]), data)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(string(pt))
		},
	}

	sendCmd := &cobra.Command{
		Use:   "send [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit an encrypted transaction JSON file",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			var j struct {
				Payload string `json:"payload"`
				Nonce   string `json:"nonce"`
			}
			if err := json.Unmarshal(b, &j); err != nil {
				fmt.Println("error:", err)
				return
			}
			payload, err := hex.DecodeString(j.Payload)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			nonce, err := hex.DecodeString(j.Nonce)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			privTxMgr.Send(core.PrivateTransaction{Payload: payload, Nonce: nonce})
		},
	}

	cmd.AddCommand(encryptCmd, decryptCmd, sendCmd)
	rootCmd.AddCommand(cmd)
}
