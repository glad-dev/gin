package edit

import (
	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type editModel struct {
	oldConfig    *config.Wrapper
	inputs       []textinput.Model
	focusIndex   int
	listIndex    int
	detailsIndex int
	width        int
	height       int
}

func (m *editModel) init(match *remote.Match, listIndex int, detailsIndex int) {
	// Set the new values
	m.inputs[0].SetValue(match.URL.String())
	m.inputs[0].SetCursor(0)
	m.inputs[1].SetValue(match.Token)
	m.listIndex = listIndex
	m.detailsIndex = detailsIndex

	// Set the focus to the first element
	m.focusIndex = 0
	m.updateFocus()
}

func updateEdit(m *model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch s := msg.String(); s {
		case "esc":
			if len(m.details.Items()) == 0 {
				m.currentlyDisplaying = displayingList

				return nil
			}

			m.currentlyDisplaying = displayingDetails

			return nil

		case "tab", "shift+tab", "enter", "up", "down": //nolint: goconst
			// Did the user press enter while the submit button was focused?
			if s == "enter" && m.edit.focusIndex == len(m.edit.inputs) {
				m.currentlyDisplaying = displayingLoading

				return nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.edit.focusIndex--
			} else {
				m.edit.focusIndex++
			}

			if m.edit.focusIndex > len(m.edit.inputs) {
				m.edit.focusIndex = 0
			} else if m.edit.focusIndex < 0 {
				m.edit.focusIndex = len(m.edit.inputs)
			}

			return m.edit.updateFocus()
		}
	}

	cmds := make([]tea.Cmd, len(m.edit.inputs))
	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.edit.inputs {
		m.edit.inputs[i], cmds[i] = m.edit.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *editModel) view() string {
	return shared.RenderInputFields(
		m.inputs,
		m.focusIndex,
		m.width,
		m.height,
	)
}

func (m *editModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = style.Focused
			m.inputs[i].TextStyle = style.Focused

			continue
		}

		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = style.None
		m.inputs[i].TextStyle = style.None
	}

	return tea.Batch(cmds...)
}
