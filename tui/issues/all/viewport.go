package all

import tea "github.com/charmbracelet/bubbletea"

func handleViewportUpdate(m *model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc", "backspace", "q":
			m.viewingList = true

			return m, nil
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}
