package cli

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
	"synnergy/internal/nodes"
)

var baseNode = core.NewBaseNode(nodes.Address("base1"))

func init() {
	bnCmd := &cobra.Command{
		Use:   "basenode",
		Short: "Manage base node lifecycle and peers",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the node",
		RunE:  func(cmd *cobra.Command, args []string) error { return baseNode.Start() },
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the node",
		RunE:  func(cmd *cobra.Command, args []string) error { return baseNode.Stop() },
	}

	runningCmd := &cobra.Command{
		Use:   "running",
		Short: "Check if the node is running",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println(baseNode.IsRunning()) },
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List known peers",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range baseNode.Peers() {
				fmt.Println(p)
			}
		},
	}

	var dialPub, dialSig string
	dialCmd := &cobra.Command{
		Use:   "dial [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Dial a seed peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			pub, err := hex.DecodeString(dialPub)
			if err != nil {
				return err
			}
			sig, err := hex.DecodeString(dialSig)
			if err != nil {
				return err
			}
			return baseNode.DialSeedSigned(nodes.Address(args[0]), sig, ed25519.PublicKey(pub))
		},
	}
	dialCmd.Flags().StringVar(&dialPub, "pub", "", "hex-encoded public key")
	dialCmd.Flags().StringVar(&dialSig, "sig", "", "hex-encoded signature")
	dialCmd.MarkFlagRequired("pub")
	dialCmd.MarkFlagRequired("sig")

	bnCmd.AddCommand(startCmd, stopCmd, runningCmd, peersCmd, dialCmd)
	rootCmd.AddCommand(bnCmd)
}
