package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
	"synnergy/core"
)

var connectionManager = core.NewChainConnectionManager()

func init() {
	cmd := &cobra.Command{
		Use:   "cross_chain_connection",
		Short: "Create and monitor chain connections",
	}

	var listJSON bool
	var getJSON bool

	openCmd := &cobra.Command{
		Use:   "open <local_chain> <remote_chain>",
		Args:  cobra.ExactArgs(2),
		Short: "Establish a new connection",
		Run: func(cmd *cobra.Command, args []string) {
			id := connectionManager.Open(args[0], args[1])
			fmt.Printf("%d gas:%d\n", id, synnergy.GasCost("OpenConnection"))
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
			if err := connectionManager.Close(id); err != nil {
				return err
			}
			fmt.Printf("gas:%d\n", synnergy.GasCost("CloseConnection"))
			return nil
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
			if getJSON {
				enc, _ := json.Marshal(c)
				fmt.Println(string(enc))
				return nil
			}
			fmt.Printf("%d: %s <-> %s open=%v\n", c.ID, c.LocalChain, c.RemoteChain, c.Open)
			return nil
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List active and historic connections",
		Run: func(cmd *cobra.Command, args []string) {
			cs := connectionManager.List()
			if listJSON {
				enc, _ := json.Marshal(cs)
				fmt.Println(string(enc))
				return
			}
			for _, c := range cs {
				fmt.Printf("%d: %s <-> %s open=%v\n", c.ID, c.LocalChain, c.RemoteChain, c.Open)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	cmd.AddCommand(openCmd, closeCmd, getCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
