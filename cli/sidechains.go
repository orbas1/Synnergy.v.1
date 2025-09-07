package cli

import (
        "fmt"
        "strconv"

        "github.com/spf13/cobra"
        "synnergy/core"
)

var (
	sideReg = core.NewSidechainRegistry()
	sideOps = core.NewSidechainOps(sideReg)
)

func init() {
	cmd := &cobra.Command{
		Use:   "sidechain",
		Short: "Manage side-chains",
	}

	registerCmd := &cobra.Command{
		Use:   "register [id] [meta] [validators...]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Register a new side-chain",
		Run: func(cmd *cobra.Command, args []string) {
                        sc, err := sideReg.Register(args[0], args[1], args[2:])
                        if err != nil {
                                fmt.Println("error:", err)
                                return
                        }
                        printOutput(sc)
		},
	}

	headerCmd := &cobra.Command{
		Use:   "header [id] [header]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a side-chain header",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sideReg.SubmitHeader(args[0], args[1]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	getHeaderCmd := &cobra.Command{
		Use:   "get-header [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch a submitted side-chain header",
		Run: func(cmd *cobra.Command, args []string) {
                        if h, ok := sideReg.GetHeader(args[0]); ok {
                                printOutput(h)
                        } else {
                                printOutput("not found")
                        }
		},
	}

	metaCmd := &cobra.Command{
		Use:   "meta [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Display side-chain metadata",
		Run: func(cmd *cobra.Command, args []string) {
                        if sc, ok := sideReg.Meta(args[0]); ok {
                                printOutput(sc)
                        } else {
                                printOutput("not found")
                        }
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered side-chains",
		Run: func(cmd *cobra.Command, args []string) {
                        scs := sideReg.List()
                        printOutput(scs)
                },
        }

	pauseCmd := &cobra.Command{
		Use:   "pause [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Pause a side-chain",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sideReg.Pause(args[0]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Resume a paused side-chain",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sideReg.Resume(args[0]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update-validators [id] [validators...]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Update side-chain validator set",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sideReg.UpdateValidators(args[0], args[1:]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a side-chain and all data",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sideReg.Remove(args[0]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	depositCmd := &cobra.Command{
		Use:   "deposit [chainID] [from] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit tokens to a side-chain escrow",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := sideOps.Deposit(args[0], args[1], amt); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	withdrawCmd := &cobra.Command{
		Use:   "withdraw [chainID] [from] [amount] [proof]",
		Args:  cobra.ExactArgs(4),
		Short: "Verify a withdrawal proof",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := sideOps.Withdraw(args[0], args[1], amt, args[3]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	cmd.AddCommand(registerCmd, headerCmd, getHeaderCmd, metaCmd, listCmd, pauseCmd, resumeCmd, updateCmd, removeCmd, depositCmd, withdrawCmd)
	rootCmd.AddCommand(cmd)
}
