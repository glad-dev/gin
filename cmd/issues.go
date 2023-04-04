package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"gn/repo"
	allIssues "gn/tui/issues/all"
	singleIssue "gn/tui/issues/single"
	"gn/tui/style"

	"github.com/spf13/cobra"
)

func newCmdAllIssues() *cobra.Command {
	cmdAllIssues := &cobra.Command{
		Use:               "issues",
		Short:             "View all issues of a repository",
		Long:              "Long - Query all issues",
		Args:              cobra.ExactArgs(0),
		PersistentPreRunE: preRun,
		Run:               runAllIssues,
	}

	cmdAllIssues.PersistentFlags().String("path", "", "Path to the repo")
	cmdAllIssues.PersistentFlags().String("url", "", "URL of the repo")

	return cmdAllIssues
}

func newCmdSingleIssue() *cobra.Command {
	cmdSingleIssue := &cobra.Command{
		Use:               "issue [iid]",
		Short:             "View the discussion of an issue",
		Long:              "Long - Show single issue",
		Args:              cobra.ExactArgs(1),
		PersistentPreRunE: preRun,
		Run:               runSingleIssue,
	}

	cmdSingleIssue.PersistentFlags().String("path", "", "Path to the repo")
	cmdSingleIssue.PersistentFlags().String("url", "", "URL of the repo")

	return cmdSingleIssue
}

func runAllIssues(cmd *cobra.Command, _ []string) {
	details, u, err := getDetailsOrURL(cmd)
	if err != nil {
		style.PrintErrAndExit(err.Error())
	}

	allIssues.Show(details, u)
}

func runSingleIssue(cmd *cobra.Command, args []string) {
	details, u, err := getDetailsOrURL(cmd)
	if err != nil {
		style.PrintErrAndExit(err.Error())
	}

	singleIssue.Show(details, u, args[0])
}

func preRun(cmd *cobra.Command, _ []string) error {
	urlFlag := cmd.Flags().Lookup("url")
	pathFlag := cmd.Flags().Lookup("path")

	// Check if flags exists and if they were set
	if (urlFlag != nil && pathFlag != nil) && (urlFlag.Changed && pathFlag.Changed) {
		return errors.New("flags --path and --url are mutually exclusive")
	}

	return nil
}

func getDetailsOrURL(cmd *cobra.Command) ([]repo.Details, *url.URL, error) {
	// Check if the flags are ok
	dir, err := cmd.Flags().GetString("path")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get the 'path' flag: %w", err)
	}

	urlStr, err := cmd.Flags().GetString("url")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get the 'url' flag: %w", err)
	}

	// Due to our pre-run hook, we know that only one of the flags is set
	if len(urlStr) > 0 {
		// We were passed a URL flag
		var u *url.URL
		u, err = url.ParseRequestURI(urlStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse given url: %w", err)
		}

		return nil, u, nil
	}

	// We were either passed a path flag or no flag
	if len(dir) == 0 {
		// Path flag was not set => Use current directory
		dir, err = os.Getwd()
		if err != nil {
			return nil, nil, err
		}
	}

	details, err := repo.Get(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get repo details: %w", err)
	}

	// Get the git repository at the current directory
	return details, nil, err
}
