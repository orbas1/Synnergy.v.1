package cli

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
)

var privTxMgr = synnergy.NewPrivateTxManager()

func init() {
	privCmd := &cobra.Command{Use: "privtx", Short: "Private transaction utilities"}

	sendCmd := &cobra.Command{
		Use:   "send [hexkey] [plaintext]",
		Args:  cobra.ExactArgs(2),
		Short: "Encrypt and store a private transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			encrypted, err := synnergy.Encrypt(key, []byte(args[1]))
			if err != nil {
				return err
			}
			block, err := aes.NewCipher(key)
			if err != nil {
				return err
			}
			gcm, err := cipher.NewGCM(block)
			if err != nil {
				return err
			}
			nonceSize := gcm.NonceSize()
			tx := synnergy.PrivateTransaction{Nonce: encrypted[:nonceSize], Payload: encrypted[nonceSize:]}
			privTxMgr.Send(tx)
			fmt.Println("transaction stored")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List stored private transactions",
		Run: func(cmd *cobra.Command, args []string) {
			txs := privTxMgr.List()
			for i, tx := range txs {
				fmt.Printf("%d: nonce=%x payload=%x\n", i, tx.Nonce, tx.Payload)
			}
		},
	}

	decryptCmd := &cobra.Command{
		Use:   "decrypt [hexkey] [index]",
		Args:  cobra.ExactArgs(2),
		Short: "Decrypt a stored transaction by index",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			idx, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			txs := privTxMgr.List()
			if idx < 0 || idx >= len(txs) {
				return fmt.Errorf("index out of range")
			}
			data := append(txs[idx].Nonce, txs[idx].Payload...)
			plain, err := synnergy.Decrypt(key, data)
			if err != nil {
				return err
			}
			fmt.Println(string(plain))
			return nil
		},
	}

	privCmd.AddCommand(sendCmd, listCmd, decryptCmd)
	rootCmd.AddCommand(privCmd)
}
