package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var carbonReg = tokens.NewCarbonRegistry()

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
			p := carbonReg.Register(owner, name, total)
			fmt.Println(p.ID)
		},
	}
	registerCmd.Flags().String("owner", "", "project owner")
	registerCmd.Flags().String("name", "", "project name")
	registerCmd.Flags().Uint64("total", 0, "total credits")
	cmd.AddCommand(registerCmd)

	issueCmd := &cobra.Command{
		Use:   "issue <projectID> <holder> <amount>",
		Short: "Issue credits to holder",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := carbonReg.Issue(args[0], args[1], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("issued")
			}
		},
	}
	cmd.AddCommand(issueCmd)

	retireCmd := &cobra.Command{
		Use:   "retire <projectID> <holder> <amount>",
		Short: "Retire credits from holder",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := carbonReg.Retire(args[0], args[1], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("retired")
			}
		},
	}
	cmd.AddCommand(retireCmd)

	verifyCmd := &cobra.Command{
		Use:   "verify <projectID> <verifier> <recordID> <status>",
		Short: "Add verification record",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			if err := carbonReg.AddVerification(args[0], args[1], args[2], args[3]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("verification recorded")
			}
		},
	}
	cmd.AddCommand(verifyCmd)

	verifsCmd := &cobra.Command{
		Use:   "verifications <projectID>",
		Short: "List verifications for project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			v, ok := carbonReg.Verifications(args[0])
			if !ok {
				fmt.Println("project not found")
				return
			}
			for _, rec := range v {
				fmt.Printf("%+v\n", rec)
			}
		},
	}
	cmd.AddCommand(verifsCmd)

	infoCmd := &cobra.Command{
		Use:   "info <projectID>",
		Short: "Show project information",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := carbonReg.ProjectInfo(args[0])
			if !ok {
				fmt.Println("project not found")
				return
			}
			fmt.Printf("%+v\n", *p)
		},
	}
	cmd.AddCommand(infoCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range carbonReg.ListProjects() {
				fmt.Printf("%+v\n", *p)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
