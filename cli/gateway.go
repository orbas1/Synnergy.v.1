package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
	nodes "synnergy/internal/nodes"
)

var gateway = core.NewGatewayNode(nodes.Address("gw1"), core.GatewayConfig{})

func init() {
	gwCmd := &cobra.Command{
		Use:   "gateway",
		Short: "Gateway node endpoint management",
	}

	registerCmd := &cobra.Command{
		Use:   "register [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Register an endpoint",
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			gateway.RegisterEndpoint(name, func(b []byte) error {
				fmt.Printf("%s received: %s\n", name, string(b))
				return nil
			})
		},
	}

	callCmd := &cobra.Command{
		Use:   "call [name] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Invoke an endpoint",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gateway.Handle(args[0], []byte(args[1])); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove an endpoint",
		Run: func(cmd *cobra.Command, args []string) {
			gateway.RemoveEndpoint(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered endpoints",
		Run: func(cmd *cobra.Command, args []string) {
			for _, e := range gateway.Endpoints() {
				fmt.Println(e)
			}
		},
	}

	gwCmd.AddCommand(registerCmd, callCmd, removeCmd, listCmd)
	rootCmd.AddCommand(gwCmd)
}
