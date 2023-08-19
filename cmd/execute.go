package cmd

import (
	"github.com/glad-dev/gin/constants"

	"github.com/spf13/cobra"
)

// Execute parses the command line parameters and starts the program.
func Execute() error {
	rootCmd := &cobra.Command{
		Use:     "gin",
		Version: constants.Version,
	}

	cmdAllIssues := newCmdAllIssues()
	cmdSingleIssue := newCmdSingleIssue()
	cmdConfig := newCmdConfig()
	// cmdList := newCmdList()

	rootCmd.AddCommand(cmdAllIssues, cmdSingleIssue, cmdConfig)

	return rootCmd.Execute()
}
