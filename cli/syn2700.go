package cli

import (
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
		kv := strings.SplitN(p, ":", 2)
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
		Run: func(cmd *cobra.Command, args []string) {
			data, _ := cmd.Flags().GetString("entries")
			schedule = core.NewVestingSchedule(parseEntries(data))
			fmt.Println("schedule created")
		},
	}
	createCmd.Flags().String("entries", "", "entry as time:amount,comma-separated")
	cmd.AddCommand(createCmd)

	claimCmd := &cobra.Command{
		Use:   "claim",
		Short: "Claim vested amounts",
		Run: func(cmd *cobra.Command, args []string) {
			if schedule == nil {
				fmt.Println("schedule not created")
				return
			}
			amt := schedule.Claim(time.Now())
			fmt.Printf("claimed %d\n", amt)
		},
	}
	cmd.AddCommand(claimCmd)

	pendingCmd := &cobra.Command{
		Use:   "pending",
		Short: "Show pending amount",
		Run: func(cmd *cobra.Command, args []string) {
			if schedule == nil {
				fmt.Println("schedule not created")
				return
			}
			fmt.Printf("pending %d\n", schedule.Pending(time.Now()))
		},
	}
	cmd.AddCommand(pendingCmd)

	rootCmd.AddCommand(cmd)
}
