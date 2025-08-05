package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var event *core.EventMetadata

func init() {
	cmd := &cobra.Command{
		Use:   "syn1700",
		Short: "Event ticket token",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise event metadata",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("desc")
			location, _ := cmd.Flags().GetString("location")
			start, _ := cmd.Flags().GetInt64("start")
			end, _ := cmd.Flags().GetInt64("end")
			supply, _ := cmd.Flags().GetUint64("supply")
			event = core.NewEvent(name, desc, location, start, end, supply)
			fmt.Println("event initialised")
		},
	}
	initCmd.Flags().String("name", "", "event name")
	initCmd.Flags().String("desc", "", "description")
	initCmd.Flags().String("location", "", "location")
	initCmd.Flags().Int64("start", 0, "start unix time")
	initCmd.Flags().Int64("end", 0, "end unix time")
	initCmd.Flags().Uint64("supply", 0, "ticket supply")
	cmd.AddCommand(initCmd)

	issueCmd := &cobra.Command{
		Use:   "issue <owner> <class> <type> <price>",
		Args:  cobra.ExactArgs(4),
		Short: "Issue a ticket",
		Run: func(cmd *cobra.Command, args []string) {
			if event == nil {
				fmt.Println("event not initialised")
				return
			}
			price, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				fmt.Println("invalid price")
				return
			}
			id, err := event.IssueTicket(args[0], args[1], args[2], price)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(id)
		},
	}
	cmd.AddCommand(issueCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <from> <to>",
		Args:  cobra.ExactArgs(3),
		Short: "Transfer a ticket",
		Run: func(cmd *cobra.Command, args []string) {
			if event == nil {
				fmt.Println("event not initialised")
				return
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid id")
				return
			}
			if err := event.TransferTicket(id, args[1], args[2]); err != nil {
				fmt.Println(err)
			}
		},
	}
	cmd.AddCommand(transferCmd)

	verifyCmd := &cobra.Command{
		Use:   "verify <id> <holder>",
		Args:  cobra.ExactArgs(2),
		Short: "Verify ticket ownership",
		Run: func(cmd *cobra.Command, args []string) {
			if event == nil {
				fmt.Println("event not initialised")
				return
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid id")
				return
			}
			fmt.Println(event.VerifyTicket(id, args[1]))
		},
	}
	cmd.AddCommand(verifyCmd)

	rootCmd.AddCommand(cmd)
}
