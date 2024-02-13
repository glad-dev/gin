package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCmdListRepos() *cobra.Command { // nolint:unused
	return &cobra.Command{
		Use:   "list",
		Short: "List all repos that you have access to",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Show all repos - ToDo")
		},
	}
}
