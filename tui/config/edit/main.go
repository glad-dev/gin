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
	for i := range wrapper.Configs {
		items[i] = shared.ListItem{
			Remote: &wrapper.Configs[i],
		}
	}

	lst := list.New(items, shared.ItemDelegate{}, 0, 0)
	lst.Title = "Which remote do you want to edit?"
	lst.SetFilteringEnabled(false)
	lst.SetShowStatusBar(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	detailsLst := list.New([]list.Item{}, detailsItemDelegate{}, 0, 0)
	detailsLst.Title = "Which token do you want to edit?"
	detailsLst.SetFilteringEnabled(false)
	detailsLst.SetShowStatusBar(false)
	detailsLst.Styles.Title = style.Title
	detailsLst.Styles.PaginationStyle = style.Pagination
	detailsLst.Styles.HelpStyle = style.Help

	p := tea.NewProgram(model{
		remotes: lst,
		details: detailsLst,
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
