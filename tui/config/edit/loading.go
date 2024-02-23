package edit

import (
	"errors"
	"fmt"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateLoading(m *model) tea.Cmd {
	return func() tea.Msg {
		oldURL := m.edit.originalConfig.Remotes[m.edit.listIndex].URL.String()

		err := configuration.Update(
			m.edit.originalConfig,
			m.edit.listIndex,
			m.edit.detailsIndex,
			m.edit.inputs[0].Value(),
			m.edit.inputs[1].Value(),
			true,
		)
		if err != nil {
			if errors.Is(err, configuration.ErrUpdateSameValues) {
				return updateMsg{
					str:     style.FormatQuitText("No need to update the config: No changes were made."),
					success: true,
				}
			}

			return updateMsg{
				str:     style.FormatQuitText(fmt.Sprintf("Failed to update remote: %s", err)),
				success: false,
			}
		}

		return updateMsg{
			str:     style.FormatQuitText(fmt.Sprintf("Sucessfully updated the remote %s", oldURL)),
			success: true,
		}
	}
}

func viewLoading(m *model) string {
	return lipgloss.Place(
		m.edit.width,
		m.edit.height,
		lipgloss.Center,
		lipgloss.Center,

		"Verifying token "+m.spinner.View(),
	)
}
