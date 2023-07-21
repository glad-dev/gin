package add

import (
	"fmt"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateLoading(m *model) tea.Cmd {
	return func() tea.Msg {
		err := config.Append(m.inputs[0].Value(), m.inputs[1].Value(), m.remoteType)
		if err != nil {
			return updateMsg{
				str:     fmt.Sprintf("Could not add config: %s", err),
				success: false,
			}
		}

		return updateMsg{
			str:     style.FormatQuitText(fmt.Sprintf("Successfully added config for %s", m.inputs[0].Value())),
			success: true,
		}
	}
}

func viewLoading(m *model) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,

		"Verifying token "+m.spinner.View(),
	)
}
