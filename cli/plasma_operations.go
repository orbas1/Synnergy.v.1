package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	opsCmd := &cobra.Command{
		Use:   "plasma-ops",
		Short: "Plasma bridge deposits and exits",
	}

	depositCmd := &cobra.Command{
		Use:   "deposit [owner] [token] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit tokens into the Plasma bridge",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			if err := plasmaBridge.Deposit(args[0], args[1], amt); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	exitCmd := &cobra.Command{
		Use:   "exit [owner] [token] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Start an exit from the bridge",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			nonce, err := plasmaBridge.StartExit(args[0], args[1], amt)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(nonce)
		},
	}

	finalizeCmd := &cobra.Command{
		Use:   "finalize [nonce]",
		Args:  cobra.ExactArgs(1),
		Short: "Finalize a pending exit",
		Run: func(cmd *cobra.Command, args []string) {
			n, _ := strconv.ParseUint(args[0], 10, 64)
			if err := plasmaBridge.FinalizeExit(n); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [nonce]",
		Args:  cobra.ExactArgs(1),
		Short: "Get details of an exit",
		Run: func(cmd *cobra.Command, args []string) {
			n, _ := strconv.ParseUint(args[0], 10, 64)
			ex, err := plasmaBridge.GetExit(n)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Printf("%+v\n", ex)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [owner]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List exits, optionally by owner",
		Run: func(cmd *cobra.Command, args []string) {
			owner := ""
			if len(args) == 1 {
				owner = args[0]
			}
			for _, ex := range plasmaBridge.ListExits(owner) {
				fmt.Printf("%+v\n", ex)
			}
		},
	}

	opsCmd.AddCommand(depositCmd, exitCmd, finalizeCmd, getCmd, listCmd)
	plasmaCmd.AddCommand(opsCmd)
}
