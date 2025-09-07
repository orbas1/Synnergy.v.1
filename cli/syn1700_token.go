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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("desc")
			location, _ := cmd.Flags().GetString("location")
			start, _ := cmd.Flags().GetInt64("start")
			end, _ := cmd.Flags().GetInt64("end")
			supply, _ := cmd.Flags().GetUint64("supply")
			if name == "" || desc == "" || location == "" || start == 0 || end == 0 || supply == 0 {
				return fmt.Errorf("name, desc, location, start, end and supply must be provided")
			}
			event = core.NewEvent(name, desc, location, start, end, supply)
			cmd.Println("event initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "event name")
	initCmd.Flags().String("desc", "", "description")
	initCmd.Flags().String("location", "", "location")
	initCmd.Flags().Int64("start", 0, "start unix time")
	initCmd.Flags().Int64("end", 0, "end unix time")
	initCmd.Flags().Uint64("supply", 0, "ticket supply")
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("desc")
	initCmd.MarkFlagRequired("location")
	initCmd.MarkFlagRequired("start")
	initCmd.MarkFlagRequired("end")
	initCmd.MarkFlagRequired("supply")
	cmd.AddCommand(initCmd)

	issueCmd := &cobra.Command{
		Use:   "issue <owner> <class> <type> <price>",
		Args:  cobra.ExactArgs(4),
		Short: "Issue a ticket",
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == nil {
				return fmt.Errorf("event not initialised")
			}
			price, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid price")
			}
			id, err := event.IssueTicket(args[0], args[1], args[2], price)
			if err != nil {
				return err
			}
			cmd.Println(id)
			return nil
		},
	}
	cmd.AddCommand(issueCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <from> <to>",
		Args:  cobra.ExactArgs(3),
		Short: "Transfer a ticket",
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == nil {
				return fmt.Errorf("event not initialised")
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			return event.TransferTicket(id, args[1], args[2])
		},
	}
	cmd.AddCommand(transferCmd)

	verifyCmd := &cobra.Command{
		Use:   "verify <id> <holder>",
		Args:  cobra.ExactArgs(2),
		Short: "Verify ticket ownership",
		RunE: func(cmd *cobra.Command, args []string) error {
			if event == nil {
				return fmt.Errorf("event not initialised")
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			cmd.Println(event.VerifyTicket(id, args[1]))
			return nil
		},
	}
	cmd.AddCommand(verifyCmd)

	rootCmd.AddCommand(cmd)
}
