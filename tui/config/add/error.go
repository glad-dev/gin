package add

import (
	"fmt"

	"gn/tui/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateError(msg tea.Msg) {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "backspace":
			m.error = ""
			m.currentlyDisplaying = displayingAdd
		}
	}
}

func (m *model) viewError() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		0.75,

		fmt.Sprintf(
			"%s\n%s\n\nPress the 'q', 'esc' or 'backspace' key to go back.",
			style.Error.Render("An error occurred:"),
			m.error,
		),
	)
}
