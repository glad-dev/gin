package cmd

import (
	"fmt"
	"gn/constants"

	"github.com/spf13/cobra"
)

func Execute() error {
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

	var config = &cobra.Command{
		Use:   "config [command]",
		Short: "Interact with config",
		Long:  "Long - edit config",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("config")
		},
	}

	var editConfig = &cobra.Command{
		Use:   "add",
		Short: "Add remote",
		Long:  "Long - add remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("config - edit")
		},
	}

	var removeConfig = &cobra.Command{
		Use:   "remove",
		Short: "Remove a remote",
		Long:  "Long - Remove a remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("config - remove remote")
		},
	}

	var rootCmd = &cobra.Command{
		Use:     "gn",
		Version: constants.Version,
	}

	rootCmd.AddCommand(issues, list, config)
	issues.AddCommand(allIssues, singleIssue)
	config.AddCommand(editConfig, removeConfig)

	return rootCmd.Execute()
}
