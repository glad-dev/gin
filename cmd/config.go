package cmd

import (
	"fmt"
	"os"

	"gn/config"
	"gn/tui/config/add"
	"gn/tui/config/remove"

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

	cmdList := &cobra.Command{
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

	cmdEdit := &cobra.Command{
		Use:   "add",
		Short: "Add remote",
		Long:  "Long - add remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			add.Config()
		},
	}

	cmdRemove := &cobra.Command{
		Use:   "remove",
		Short: "Remove a remote",
		Long:  "Long - Remove a remote and its API token",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			remove.Config()
		},
	}

	cmdUpdate := &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration of an existing remote",
		Long:  "Long - edit config",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.UpdateToken()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
				os.Exit(1)
			}
		},
	}

	root.AddCommand(cmdList, cmdEdit, cmdRemove, cmdUpdate)

	return root
}
