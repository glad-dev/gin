package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/repo"
	"github.com/glad-dev/gin/style"
	allIssues "github.com/glad-dev/gin/tui/issues/all"
	singleIssue "github.com/glad-dev/gin/tui/issues/single"

	"github.com/spf13/cobra"
)

func newCmdAllIssues() *cobra.Command {
	cmdAllIssues := &cobra.Command{
		Use:               "issues",
		Short:             "View all issues of a repository",
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

	// Check if the flags were defined
	if urlFlag == nil || pathFlag == nil {
		logger.Log.Error("URL or path flag was not defined.", "urlFlag", urlFlag, "pathFlag", pathFlag)

		return errors.New("flag --path or --url was not defined")
	}

	// Check if flags exists and if they were set
	if urlFlag.Changed && pathFlag.Changed {
		logger.Log.Error("User set both --path and --url.")

		return errors.New("flags --path and --url are mutually exclusive")
	}

	return nil
}

func getDetailsOrURL(cmd *cobra.Command) ([]repo.Details, *url.URL, error) {
	// Check if the flags are ok
	dir, err := cmd.Flags().GetString("path")
	if err != nil {
		logger.Log.Error("Failed to get the path flag", "error", err)

		return nil, nil, fmt.Errorf("failed to get the 'path' flag: %w", err)
	}

	urlStr, err := cmd.Flags().GetString("url")
	if err != nil {
		logger.Log.Error("Failed to get the url flag", "error", err)

		return nil, nil, fmt.Errorf("failed to get the 'url' flag: %w", err)
	}

	// Due to our pre-run hook, we know that only one of the flags is set
	if len(urlStr) > 0 {
		// We were passed a URL flag
		var u *url.URL
		u, err = url.ParseRequestURI(urlStr)
		if err != nil {
			logger.Log.Error("Invalid url passed.", "error", err, "url", urlStr)

			return nil, nil, fmt.Errorf("failed to parse given url: %w", err)
		}

		return nil, u, nil
	}

	// We were either passed a path flag or no flag
	if len(dir) == 0 {
		// Path flag was not set => Use current directory
		dir, err = os.Getwd()
		if err != nil {
			logger.Log.Error("Failed to get current directory", "error", err)

			return nil, nil, err
		}
	}

	details, err := repo.Get(dir)
	if err != nil {
		logger.Log.Error("Failed to get repository details.", "error", err, "directory", dir)

		return nil, nil, fmt.Errorf("failed to get repo details: %w", err)
	}

	// Get the git repository at the current directory
	return details, nil, err
}
