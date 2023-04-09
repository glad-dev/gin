package selectconfig

import (
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list     list.Model
	done     bool
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.InputField.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			m.done = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.quitting || m.done {
		return ""
	}

	return shared.RenderList(m.list)
}
