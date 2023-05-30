package selectconfig

import (
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type state uint8

const (
	stateRunning state = iota
	exitSuccess
	exitNoSelection
)

type model struct {
	list  list.Model
	state state
}

// Init is required for model to be a tea.Model.
func (m model) Init() tea.Cmd {
	return nil
}

// Update is required for model to be a tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.InputField.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.state = exitNoSelection

			return m, tea.Quit

		case "enter":
			m.state = exitSuccess

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

// View required for model to be a tea.Model.
func (m model) View() string {
	if m.state == stateRunning {
		return shared.RenderList(m.list)
	}

	return ""
}
