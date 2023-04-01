package add

import (
	"os"

	"gn/tui/config/shared"
	"gn/tui/style"

	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	p := tea.NewProgram(model{
		inputs:  shared.GetTextInputs(),
		failure: false,
	})

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Failed to start program: " + err.Error())
	}

	if m, ok := m.(model); ok && m.failure {
		os.Exit(1)
	}
}
