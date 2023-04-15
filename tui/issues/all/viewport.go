package all

import tea "github.com/charmbracelet/bubbletea"

func (m *model) updateViewport(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "backspace", "q":
			m.currentlyDisplaying = displayingList

			return nil
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return cmd
}
