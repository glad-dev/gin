package color

import (
	"github.com/glad-dev/gin/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateError(msg tea.Msg) {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "backspace":
			m.text = ""
			m.currentlyDisplaying = displayingColors
		}
	}
}

func (m *model) viewError() string {
	width := m.width - style.InputField.GetHorizontalFrameSize()

	return lipgloss.Place(
		width,
		m.height,
		lipgloss.Center,
		0.75,

		lipgloss.JoinVertical(
			lipgloss.Center,
			style.Error.Render("An error occurred:"),
			lipgloss.NewStyle().Width(width).Align(lipgloss.Center, lipgloss.Center).Render(m.text),
			"\n",
			"Press the 'q', 'esc' or 'backspace' key to go back.",
		),
	)
}
