package cmd

import (
	"fmt"
	"log"
	"os"

	"gn/config"
	"gn/issues"
	"gn/repo"

	"github.com/spf13/cobra"
)

func newCmdAllIssues() *cobra.Command {
	allIssues := &cobra.Command{
		Use:   "issues",
		Short: "View all issues of a repository",
		Long:  "Long - Query all issues",
		Args:  cobra.ExactArgs(0),
		Run:   runAllIssues,
	}

	allIssues.PersistentFlags().String("path", "", "Path to the repo")

	return allIssues
}

func newCmdSingleIssue() *cobra.Command {
	singleIssue := &cobra.Command{
		Use:   "issue [iid]",
		Short: "View the discussion of an issue",
		Long:  "Long - Show single issue",
		Args:  cobra.ExactArgs(1),
		Run:   runSingleIssue,
	}

	singleIssue.PersistentFlags().String("path", "", "Path to the repo")

	return singleIssue
}

func runAllIssues(cmd *cobra.Command, _ []string) {
	conf, err := config.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
		os.Exit(1)
	}

	details, err := getRepo(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	issueList, err := issues.QueryAll(conf, details)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
		os.Exit(1)
	}

	for _, issue := range issueList {
		fmt.Printf("%s) %s [%s]\n", issue.Iid, issue.Title, issue.State)
	}
}

func runSingleIssue(cmd *cobra.Command, args []string) {
	conf, err := config.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
		os.Exit(1)
	}

	details, err := getRepo(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	issue, err := issues.QuerySingle(conf, details, args[0])
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
