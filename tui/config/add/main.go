package add

import (
	"fmt"
	"strings"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"
	"github.com/glad-dev/gin/tui/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

// Config is the entry point of this TUI, which allows to add new remotes.
func Config() {
	_, _ = config.Load() // To load the colors

	p := tea.NewProgram(model{
		inputs:              shared.GetTextInputs(),
		spinner:             *widgets.GetSpinner(),
		currentlyDisplaying: displayingType,
		state:               stateRunning,
	}, tea.WithAltScreen())

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Failed to start program: " + err.Error())
	}

	r, ok := m.(model)
	if !ok {
		return
	}

	switch r.state { //nolint:exhaustive
	case exitFailure:
		log.Error(strings.TrimSpace(r.text))
		style.PrintErrAndExit(strings.TrimSpace(r.text))

	case exitSuccess:
		fmt.Print(style.FormatQuitText("Successfully added the remote"))
	}
}
