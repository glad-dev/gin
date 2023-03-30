package remove

import (
	"errors"

	"gn/config"
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	// Load current config
	wrapper, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrConfigDoesNotExist) {
			style.PrintErrAndExit(config.ErrConfigDoesNotExistMsg)
		}

		style.PrintErrAndExit("Failed to load the configuration: " + err.Error())
	}

	items := make([]list.Item, len(wrapper.Configs))
	for i, conf := range wrapper.Configs {
		items[i] = shared.ListItem{
			Lab: conf,
		}
	}

	lst := list.New(items, shared.ItemDelegate{}, 0, 0)
	lst.Title = "Which remote do you want to delete?"
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	m := model{
		list:      lst,
		oldConfig: *wrapper,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}
}
