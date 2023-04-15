package add

import (
	"os"

	"gn/style"
	"gn/tui/config/shared"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	s := spinner.New()
	s.Spinner = spinner.Points

	p := tea.NewProgram(model{
		inputs:              shared.GetTextInputs(),
		spinner:             s,
		currentlyDisplaying: displayingAdd,
		failure:             false,
	})

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Failed to start program: " + err.Error())
	}

	if m, ok := m.(model); ok && m.failure {
		os.Exit(1)
	}
}
