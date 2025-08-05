package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var connectionManager = core.NewChainConnectionManager()

func init() {
	cmd := &cobra.Command{
		Use:   "cross_chain_connection",
		Short: "Create and monitor chain connections",
	}

	openCmd := &cobra.Command{
		Use:   "open <local_chain> <remote_chain>",
		Args:  cobra.ExactArgs(2),
		Short: "Establish a new connection",
		Run: func(cmd *cobra.Command, args []string) {
			id := connectionManager.Open(args[0], args[1])
			fmt.Println(id)
		},
	}

	closeCmd := &cobra.Command{
		Use:   "close <connection_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Terminate a connection",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			return connectionManager.Close(id)
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <connection_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve connection details",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			c, err := connectionManager.Get(id)
			if err != nil {
				return err
			}
			fmt.Printf("%d: %s <-> %s open=%v\n", c.ID, c.LocalChain, c.RemoteChain, c.Open)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List active and historic connections",
		Run: func(cmd *cobra.Command, args []string) {
			for _, c := range connectionManager.List() {
				fmt.Printf("%d: %s <-> %s open=%v\n", c.ID, c.LocalChain, c.RemoteChain, c.Open)
			}
		},
	}

	cmd.AddCommand(openCmd, closeCmd, getCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
