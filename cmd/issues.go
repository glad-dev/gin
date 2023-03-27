package cmd

import (
	"log"
	"os"

	"gn/repo"
	allIssues "gn/tui/issues/all"
	singleIssue "gn/tui/issues/single"

	"github.com/spf13/cobra"
)

func newCmdAllIssues() *cobra.Command {
	cmdAllIssues := &cobra.Command{
		Use:   "issues",
		Short: "View all issues of a repository",
		Long:  "Long - Query all issues",
		Args:  cobra.ExactArgs(0),
		Run:   runAllIssues,
	}

	cmdAllIssues.PersistentFlags().String("path", "", "Path to the repo")

	return cmdAllIssues
}

func newCmdSingleIssue() *cobra.Command {
	cmdSingleIssue := &cobra.Command{
		Use:   "issue [iid]",
		Short: "View the discussion of an issue",
		Long:  "Long - Show single issue",
		Args:  cobra.ExactArgs(1),
		Run:   runSingleIssue,
	}

	cmdSingleIssue.PersistentFlags().String("path", "", "Path to the repo")

	return cmdSingleIssue
}

func runAllIssues(cmd *cobra.Command, _ []string) {
	details, err := getRepo(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	allIssues.Show(details)
}

func runSingleIssue(cmd *cobra.Command, args []string) {
	details, err := getRepo(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	singleIssue.Show(details, args[0])
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
