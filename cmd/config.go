package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/constants"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/add"
	"github.com/glad-dev/gin/tui/config/color"
	"github.com/glad-dev/gin/tui/config/edit"
	"github.com/glad-dev/gin/tui/config/remove"

	"github.com/spf13/cobra"
)

func newCmdConfig() *cobra.Command {
	root := &cobra.Command{
		Use:   "config [command]",
		Short: "Interact with config",
		Long:  "Long - edit config",
		Args:  cobra.ExactArgs(0),
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

	addDesc := fmt.Sprintf(
		"Add a new token.\nA GitLab token needs the following scopes: %s\nA GitHub token needs the following scopes: %s\n",
		strings.Join(constants.RequiredGitLabScopes, ", "),
		strings.Join(constants.RequiredGitHubScopes, ", "),
	)

	cmdAdd := &cobra.Command{
		Use:   "add",
		Short: "Add a new remote",
		Long:  addDesc,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			add.Config()
		},
	}

	cmdRemove := &cobra.Command{
		Use:   "remove",
		Short: "Remove a remote",
		Long:  "Remove a remote and its API token from the local storage",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			remove.Config()
		},
	}

	cmdEdit := &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration of an existing remote",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			edit.Config()
		},
	}

	cmdVerify := &cobra.Command{
		Use:   "verify",
		Short: "Check the validity of all stored tokens",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.VerifyTokens()
			if err != nil {
				style.PrintErrAndExit(err.Error())
			}

			fmt.Print(style.FormatQuitText("All tokens are valid."))
		},
	}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: "Update the username and token names",
		Long:  "Update the username and token names if they changed",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.UpdateRemote()
			if err != nil {
				style.PrintErrAndExit(err.Error())
			}

			fmt.Print(style.FormatQuitText("All tokens were successfully updated."))
		},
	}

	cmdColors := &cobra.Command{
		Use:   "colors",
		Short: "Update the colors",
		Long:  "Update the colors. Delete the input field to revert back to the default color.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			color.Config()
		},
	}

	root.AddCommand(cmdList, cmdAdd, cmdRemove, cmdEdit, cmdVerify, cmdUpdate, cmdColors)

	return root
}
