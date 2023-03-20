package cmd

import (
	"gn/constants"

	"github.com/spf13/cobra"
)

//nolint:goconst
func Execute() error {
	rootCmd := &cobra.Command{
		Use:     "gn",
		Version: constants.Version,
	}

	cmdAllIssues := newCmdAllIssues()
	cmdSingleIssue := newCmdSingleIssue()
	cmdConfig := newCmdConfig()
	cmdList := newCmdList()

	rootCmd.AddCommand(cmdAllIssues, cmdSingleIssue, cmdConfig, cmdList)

	return rootCmd.Execute()
}
