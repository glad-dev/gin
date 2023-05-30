package edit

import (
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func updateDetails(m *model, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.currentlyDisplaying = displayingList
			m.details.SetItems([]list.Item{})

			return nil

		case "enter":
			selected, ok := m.remotes.SelectedItem().(shared.ListItem)
			if !ok {
				m.state = exitFailure
				m.text = style.FormatQuitText("Failed to cast selected item to list.Item")

				return tea.Quit
			}

			match, err := selected.Remote.ToMatchAtIndex(m.details.Index())
			if err != nil {
				m.state = exitFailure
				m.text = style.FormatQuitText("Failed to convert item to match: " + err.Error())

				return tea.Quit
			}

			m.currentlyDisplaying = displayingEdit
			m.edit.init(match, m.remotes.Index(), m.details.Index())

			return nil
		}
	}

	m.details, cmd = m.details.Update(msg)

	return cmd
}
