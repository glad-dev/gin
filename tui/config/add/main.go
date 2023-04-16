package add

import (
	"os"
	"strings"

	"gn/config"
	"gn/logger"
	"gn/style"
	"gn/tui/config/shared"
	"gn/tui/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	_, _ = config.Load() // To load the colors

	p := tea.NewProgram(model{
		inputs:              shared.GetTextInputs(),
		spinner:             *widgets.GetSpinner(),
		currentlyDisplaying: displayingAdd,
		state:               stateRunning,
	})

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Failed to start program: " + err.Error())
	}

	if m, ok := m.(model); ok && m.state == exitFailure {
		logger.Log.Errorf(strings.TrimSpace(m.text))
		os.Exit(1)
	}
}
