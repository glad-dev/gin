package remove

import (
	"errors"
	"os"

	"gn/config"
	"gn/style"
	"gn/tui/config/shared"

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

	items := make([]list.Item, len(wrapper.Remotes))
	for i := range wrapper.Remotes {
		items[i] = shared.ListItem{
			Remote: &wrapper.Remotes[i],
		}
	}

	p := tea.NewProgram(model{
		remotes:   shared.NewList(items, shared.ItemDelegate{}, "Which remote do you want to delete?"),
		details:   shared.NewList([]list.Item{}, shared.DetailsItemDelegate{}, "Which token do you want to delete?"),
		oldConfig: *wrapper,
	})

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}

	if m, ok := m.(model); ok && m.failure {
		os.Exit(1)
	}
}
