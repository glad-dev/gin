package remove

import (
	"fmt"

	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/lipgloss"
)

var (
	buttons     = []string{"[ no ]", "[ yes ]"}
	buttonStyle = lipgloss.NewStyle().Width(10)
)

func viewConfirmation(m *model) string {
	return lipgloss.Place(
		m.remotes.Width(),
		m.remotes.Height(),
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Center,

			lipgloss.NewStyle().Width(m.remotes.Width()).Align(lipgloss.Center, lipgloss.Center).Render(
				fmt.Sprintf(
					"Are you sure that you want to delete the token '%s' for %s?",
					m.originalConfig.Remotes[m.remotes.Index()].Details[m.details.Index()].TokenName,
					m.originalConfig.Remotes[m.remotes.Index()].URL.String(),
				),
			),
			"\n",
			renderButtons(m.confirmPosition),
		),
	)
}

func renderButtons(confirmPosition int) string {
	if confirmPosition > len(buttons) {
		confirmPosition = 0
	}

	renderedButtons := make([]string, len(buttons))
	for i := range buttons {
		renderedButtons[i] = buttonStyle.Render(buttons[i])

		if i == confirmPosition {
			renderedButtons[i] = style.Focused.Render(renderedButtons[i])
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		renderedButtons...,
	)
}
