package edit

import (
	"gn/style"
	"gn/tui/config/shared"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.state = exitNoChange

			return tea.Quit

		case "enter":
			// User selected a config
			selected, ok := m.remotes.Items()[m.remotes.Index()].(shared.ListItem)
			if !ok {
				m.state = exitFailure
				m.text = style.FormatQuitText("Failed to cast selected item to list.Item")

				return tea.Quit
			}

			if len(selected.Remote.Details) > 1 {
				m.currentlyDisplaying = displayingDetails

				items := make([]list.Item, len(selected.Remote.Details))
				for i, details := range selected.Remote.Details {
					items[i] = shared.DetailItem{
						Username:  details.Username,
						TokenName: details.TokenName,
					}
				}

				m.details.SetItems(items)
				m.details.ResetSelected()

				return nil
			}

			match, err := selected.Remote.ToMatch()
			if err != nil {
				m.state = exitFailure
				m.text = style.FormatQuitText("Failed to convert item to match: " + err.Error())

				return tea.Quit
			}

			m.currentlyDisplaying = displayingEdit
			m.edit.init(match, m.remotes.Index(), 0)

			return nil
		}
	}

	m.remotes, cmd = m.remotes.Update(msg)

	return cmd
}
