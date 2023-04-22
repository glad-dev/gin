package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCmdList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all repos that you have access to",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Show all repos - ToDo")
		},
	}
}
