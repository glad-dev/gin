package add

import (
	"fmt"
	"os"

	tui "gn/tui/config"

	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	m := model{
		inputs: tui.GetTextInputs(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
