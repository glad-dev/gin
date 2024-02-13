package edit

import (
	"errors"
	"os"
	"strings"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"
	"github.com/glad-dev/gin/tui/widgets"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Config is the entry point of this TUI, which allows to edit remotes.
func Config() {
	// Load current config
	config, err := configuration.Load()
	if err != nil {
		if errors.Is(err, configuration.ErrConfigDoesNotExist) {
			style.PrintErrAndExit(configuration.ErrConfigDoesNotExistMsg)
		}

		style.PrintErrAndExit("Failed to load the configuration: " + err.Error())
	}

	items := make([]list.Item, len(config.Remotes))
	for i := range config.Remotes {
		items[i] = shared.ListItem{
			Remote: &config.Remotes[i],
		}
	}

	p := tea.NewProgram(model{
		remotes: shared.NewList(items, shared.ItemDelegate{}, "Which remote do you want to edit?"),
		details: shared.NewList([]list.Item{}, shared.DetailsItemDelegate{}, "Which token do you want to edit?"),
		state:   stateRunning,
		spinner: *widgets.GetSpinner(),
		edit: editModel{
			inputs:         shared.GetTextInputs(),
			originalConfig: config,
			width:          0,
			height:         0,
		},
	})

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}

	if r, ok := m.(model); ok && r.state == exitFailure {
		log.Error(strings.TrimSpace(r.text))
		os.Exit(1)
	}
}
