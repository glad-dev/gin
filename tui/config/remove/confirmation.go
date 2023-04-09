package remove

import (
	"fmt"

	"gn/tui/style"

	"github.com/charmbracelet/lipgloss"
)

var (
	buttons    = []string{"[ no ]", "[ yes ]"}
	widthStyle = lipgloss.NewStyle().Width(10)
)

func (m *model) viewConfirmation() string {
	return lipgloss.Place(
		m.remotes.Width(),
		m.remotes.Height(),
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Center,

			fmt.Sprintf(
				"Are you sure that you want to delete the token '%s'?",
				m.oldConfig.Configs[m.remotes.Index()].Details[m.details.Index()].TokenName,
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

	tmp := make([]string, len(buttons))
	for i := range buttons {
		tmp[i] = widthStyle.Render(buttons[i])

		if i == confirmPosition {
			tmp[i] = style.Focused.Render(tmp[i])
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		tmp...,
	)
}
