package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Execute() {
	var issues = &cobra.Command{
		Use:   "issues [command]",
		Short: "View issues",
		Long:  "Either query all issues or a single issue with the given iid.\nIf no argument is passed, query all",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Issues")
		},
	}

	var allIssues = &cobra.Command{
		Use:   "all",
		Short: "Show all issues",
		Long:  "Show all issues of the current directory",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("All issues")
		},
	}

	var singleIssue = &cobra.Command{
		Use:   "single [iid]",
		Short: "Show all issues",
		Long:  "Show all issues of the current directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Single issue: %s\n", args[0])
		},
	}

	var list = &cobra.Command{
		Use:   "list",
		Short: "List all repos that you have access to",
		Long:  "List all repos that you have access to (long description)",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Show all repos")
		},
	}

	var rootCmd = &cobra.Command{
		Use:     "gn",
		Version: "1",
	}
	rootCmd.AddCommand(issues, list)
	issues.AddCommand(allIssues, singleIssue)
	rootCmd.Execute()
}
