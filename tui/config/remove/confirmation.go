package remove

import (
	"fmt"

	"gn/style"

	"github.com/charmbracelet/lipgloss"
)

var (
	buttons     = []string{"[ no ]", "[ yes ]"}
	buttonStyle = lipgloss.NewStyle().Width(10)
)

func (m *model) viewConfirmation() string {
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
					m.oldConfig.Remotes[m.remotes.Index()].Details[m.details.Index()].GetTokenName(),
					m.oldConfig.Remotes[m.remotes.Index()].URL.String(),
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

	btns := make([]string, len(buttons))
	for i := range buttons {
		btns[i] = buttonStyle.Render(buttons[i])

		if i == confirmPosition {
			btns[i] = style.Focused.Render(btns[i])
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		btns...,
	)
}
