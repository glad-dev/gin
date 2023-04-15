package edit

import (
	"errors"
	"os"

	"gn/tui/widgets"

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

	items := make([]list.Item, len(wrapper.Configs))
	for i := range wrapper.Configs {
		items[i] = shared.ListItem{
			Remote: &wrapper.Configs[i],
		}
	}

	p := tea.NewProgram(model{
		remotes: shared.NewList(items, shared.ItemDelegate{}, "Which remote do you want to edit?"),
		details: shared.NewList([]list.Item{}, shared.DetailsItemDelegate{}, "Which token do you want to edit?"),
		state:   stateRunning,
		spinner: *widgets.GetSpinner(),
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

	if m, ok := m.(model); ok && m.state == exitFailure {
		os.Exit(1)
	}
}
