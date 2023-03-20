package cmd

import (
	"fmt"
	"os"

	"gn/config"

	"github.com/spf13/cobra"
)

func newCmdConfig() *cobra.Command {
	root := &cobra.Command{
		Use:   "config [command]",
		Short: "Interact with config",
		Long:  "Long - edit config",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stderr, "Use commands like add|list|remove|update")
		},
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "List all remotes",
		Long:  "View a list of all existing remotes",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.List()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	edit := &cobra.Command{
		Use:   "add",
		Short: "Add remote",
		Long:  "Long - add remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.Append()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	remove := &cobra.Command{
		Use:   "remove",
		Short: "Remove a remote",
		Long:  "Long - Remove a remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.Remove()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	update := &cobra.Command{
		Use:   "update",
		Short: "Update the token of an existing remote",
		Long:  "Long - update config",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.UpdateToken()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	root.AddCommand(edit, remove, list, update)

	return root
}
