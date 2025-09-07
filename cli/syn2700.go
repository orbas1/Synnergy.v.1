package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var schedule *core.VestingSchedule

func parseEntries(s string) []core.VestingEntry {
	if s == "" {
		return nil
	}
	var entries []core.VestingEntry
	parts := strings.Split(s, ",")
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		t, err := time.Parse(time.RFC3339, kv[0])
		if err != nil {
			continue
		}
		var amt uint64
		fmt.Sscanf(kv[1], "%d", &amt)
		entries = append(entries, core.VestingEntry{ReleaseTime: t, Amount: amt})
	}
	return entries
}

func init() {
	cmd := &cobra.Command{
		Use:   "syn2700",
		Short: "Vesting schedule management",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a vesting schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, _ := cmd.Flags().GetString("entries")
			entries := parseEntries(data)
			if len(entries) == 0 {
				return errors.New("entries required")
			}
			schedule = core.NewVestingSchedule(entries)
			fmt.Println("schedule created")
			return nil
		},
	}
	createCmd.Flags().String("entries", "", "entry as time=amount,comma-separated")
	cmd.AddCommand(createCmd)

	claimCmd := &cobra.Command{
		Use:   "claim",
		Short: "Claim vested amounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if schedule == nil {
				return errors.New("schedule not created")
			}
			amt := schedule.Claim(time.Now())
			fmt.Printf("claimed %d\n", amt)
			return nil
		},
	}
	cmd.AddCommand(claimCmd)

	pendingCmd := &cobra.Command{
		Use:   "pending",
		Short: "Show pending amount",
		RunE: func(cmd *cobra.Command, args []string) error {
			if schedule == nil {
				return errors.New("schedule not created")
			}
			fmt.Printf("pending %d\n", schedule.Pending(time.Now()))
			return nil
		},
	}
	cmd.AddCommand(pendingCmd)

	rootCmd.AddCommand(cmd)
}
