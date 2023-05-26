package add

import (
	"os"
	"strings"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"
	"github.com/glad-dev/gin/tui/widgets"

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

	if r, ok := m.(model); ok && r.state == exitFailure {
		logger.Log.Errorf(strings.TrimSpace(r.text))
		os.Exit(1)
	}
}
