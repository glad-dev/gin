package cmd

import (
	"fmt"
	"log"
	"os"

	"gn/config"
	"gn/constants"
	"gn/issues"
	"gn/repo"

	"github.com/spf13/cobra"
)

//nolint:goconst
func Execute() error {
	// Load config
	conf, err := config.Get()
	if err != nil {
		return err
	}

	var cmdAllIssues = &cobra.Command{
		Use:   "issues",
		Short: "View all issues of a repository",
		Long:  "Long - Query all issues",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			details, err := getRepo(cmd) //nolint:govet
			if err != nil {
				log.Fatalln(err)
			}

			issueList, err := issues.QueryAll(conf, details) //nolint:govet
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}

			for _, issue := range issueList {
				fmt.Printf("%s) %s [%s]\n", issue.Iid, issue.Title, issue.State)
			}
		},
	}

	var cmdSingleIssue = &cobra.Command{
		Use:   "issue [iid]",
		Short: "View the discussion of an issue",
		Long:  "Long - Show single issue",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			details, err := getRepo(cmd) //nolint:govet
			if err != nil {
				log.Fatalln(err)
			}

			issue, err := issues.QuerySingle(conf, details, args[0]) //nolint:govet
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}

			fmt.Println(issue.Title)
			fmt.Println(issue.Description)
			fmt.Println()

			for _, comment := range issue.Discussion {
				fmt.Println(comment.Body)
				fmt.Printf("- %s\n", comment.Author)
				for _, subComments := range comment.Comments {
					fmt.Printf("\t%s\n", subComments.Body)
				}
			}
		},
	}

	var cmdList = &cobra.Command{
		Use:   "list",
		Short: "List all repos that you have access to",
		Long:  "List all repos that you have access to (long description)",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Show all repos - ToDo")
		},
	}

	var cmdConfig = &cobra.Command{
		Use:   "config [command]",
		Short: "Interact with config",
		Long:  "Long - edit config",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stderr, "Use commands like add|list|remove|update")
			fmt.Printf("config - Remove this?")
		},
	}

	var cmdConfigEdit = &cobra.Command{
		Use:   "add",
		Short: "Add remote",
		Long:  "Long - add remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err = config.Append()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	var cmdConfigRemove = &cobra.Command{
		Use:   "remove",
		Short: "Remove a remote",
		Long:  "Long - Remove a remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err = config.Remove()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	var cmdConfigList = &cobra.Command{
		Use:   "list",
		Short: "List all remotes",
		Long:  "View a list of all existing remotes",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err = config.List()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	var cmdConfigUpdate = &cobra.Command{
		Use:   "update",
		Short: "Update the token of an existing remote",
		Long:  "Long - update config",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err = config.UpdateToken()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	var rootCmd = &cobra.Command{
		Use:     "gn",
		Version: constants.Version,
	}

	cmdAllIssues.PersistentFlags().String("path", "", "Path to the repo")
	cmdSingleIssue.PersistentFlags().String("path", "", "Path to the repo")

	rootCmd.AddCommand(cmdAllIssues, cmdSingleIssue, cmdList, cmdConfig)
	cmdConfig.AddCommand(cmdConfigEdit, cmdConfigRemove, cmdConfigList, cmdConfigUpdate)

	return rootCmd.Execute()
}

func getRepo(cmd *cobra.Command) ([]repo.Details, error) {
	// Get path flag
	dir, err := cmd.Flags().GetString("path")
	if err != nil {
		return nil, err
	}

	if len(dir) == 0 {
		// Path flag was not set => Use current directory
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	// Get the git repository at the current directory
	return repo.Get(dir)
}
