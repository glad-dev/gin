package cmd

import (
	"errors"

	"gn/config"
	"gn/tui/config/add"
	"gn/tui/config/edit"
	"gn/tui/config/remove"
	"gn/tui/style"

	"github.com/spf13/cobra"
)

func newCmdConfig() *cobra.Command {
	root := &cobra.Command{
		Use:   "config [command]",
		Short: "Interact with config",
		Long:  "Long - edit config",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			style.PrintErrAndExit("Use commands like add|list|remove|update")
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
				if errors.Is(err, config.ErrConfigDoesNotExist) {
					style.PrintErrAndExit(config.ErrConfigDoesNotExistMsg)
				}

				style.PrintErrAndExit("An error occurred while attempting to list the configuration: " + err.Error())
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
			edit.Config()
		},
	}

	root.AddCommand(cmdList, cmdEdit, cmdRemove, cmdUpdate)

	return root
}
