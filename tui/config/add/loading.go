package add

import (
	"fmt"

	"gn/config"
	"gn/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateLoading() tea.Cmd {
	return func() tea.Msg {
		err := config.Append(m.inputs[0].Value(), m.inputs[1].Value())
		if err != nil {
			return updateMsg{
				str:     fmt.Sprintf("Could not add config: %s", err),
				failure: true,
			}
		}

		return updateMsg{
			str:     style.FormatQuitText(fmt.Sprintf("Successfully added config for %s", m.inputs[0].Value())),
			failure: false,
		}
	}
}

func (m *model) viewLoading() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,

		"Verifying token "+m.spinner.View(),
	)
}
