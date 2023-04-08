package edit

import (
	"errors"
	"os"

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
		match, err := conf.ToMatch()
		if err != nil {
			style.PrintErrAndExit("Failed to parse config: " + err.Error())
		}

		items[i] = shared.ListItem{
			Match: *match,
		}
	}

	lst := list.New(items, shared.ItemDelegate{}, 0, 0)
	lst.Title = "Which remote do you want to edit?"
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	p := tea.NewProgram(model{
		list:    lst,
		failure: false,
		edit: editModel{
			inputs:    shared.GetTextInputs(),
			oldConfig: wrapper,
			width:     0,
			height:    0,
		},
	})

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}

	if m, ok := m.(model); ok && m.failure {
		os.Exit(1)
	}
}
