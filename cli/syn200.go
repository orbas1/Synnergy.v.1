package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var carbonRegistry = tokens.NewCarbonRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn200",
		Short: "Carbon credit registry",
	}

	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a carbon project",
		Run: func(cmd *cobra.Command, args []string) {
			owner, _ := cmd.Flags().GetString("owner")
			name, _ := cmd.Flags().GetString("name")
			total, _ := cmd.Flags().GetUint64("total")
			p := carbonRegistry.Register(owner, name, total)
			fmt.Println(p.ID)
		},
	}
	registerCmd.Flags().String("owner", "", "project owner")
	registerCmd.Flags().String("name", "", "project name")
	registerCmd.Flags().Uint64("total", 0, "total credits")
	cmd.AddCommand(registerCmd)

	issueCmd := &cobra.Command{
		Use:   "issue <project> <to> <amount>",
		Short: "Issue carbon credits",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := carbonRegistry.Issue(args[0], args[1], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(issueCmd)

	retireCmd := &cobra.Command{
		Use:   "retire <project> <holder> <amount>",
		Short: "Retire credits from circulation",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := carbonRegistry.Retire(args[0], args[1], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(retireCmd)

	verifyCmd := &cobra.Command{
		Use:   "verify <project> <verifier> <recordID> [status]",
		Short: "Add a verification record",
		Args:  cobra.RangeArgs(3, 4),
		Run: func(cmd *cobra.Command, args []string) {
			status := ""
			if len(args) == 4 {
				status = args[3]
			}
			if err := carbonRegistry.AddVerification(args[0], args[1], args[2], status); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(verifyCmd)

	verifsCmd := &cobra.Command{
		Use:   "verifications <project>",
		Short: "List project verifications",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			verifs, ok := carbonRegistry.Verifications(args[0])
			if !ok {
				fmt.Println("project not found")
				return
			}
			for _, v := range verifs {
				fmt.Printf("%s %s %s %s\n", v.Verifier, v.RecordID, v.Status, v.Time.Format(time.RFC3339))
			}
		},
	}
	cmd.AddCommand(verifsCmd)

	infoCmd := &cobra.Command{
		Use:   "info <project>",
		Short: "Show project info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := carbonRegistry.ProjectInfo(args[0])
			if !ok {
				fmt.Println("project not found")
				return
			}
			fmt.Printf("ID:%s Owner:%s Name:%s Total:%d Issued:%d Retired:%d\n", p.ID, p.Owner, p.Name, p.TotalCredits, p.IssuedCredits, p.RetiredCredits)
		},
	}
	cmd.AddCommand(infoCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all carbon projects",
		Run: func(cmd *cobra.Command, args []string) {
			projects := carbonRegistry.ListProjects()
			for _, p := range projects {
				fmt.Printf("%s %s %s %d %d %d\n", p.ID, p.Owner, p.Name, p.TotalCredits, p.IssuedCredits, p.RetiredCredits)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
