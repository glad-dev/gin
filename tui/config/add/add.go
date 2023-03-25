package add

import (
	"fmt"
	"os"

	"gn/tui/config/shared"

	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	m := model{
		inputs: shared.GetTextInputs(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
