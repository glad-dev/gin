package add

import (
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"

	tea "github.com/charmbracelet/bubbletea"
)

func updateAdd(m *model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.currentlyDisplaying = displayingType
			m.inputs[0].SetValue("")
			m.inputs[1].SetValue("")
			m.focusIndex = 0

			return nil

			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.currentlyDisplaying = displayingLoading

				return nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			return updateInputs(m)
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func viewAdd(m *model) string {
	return shared.RenderInputFields(
		m.inputs,
		m.remoteType.TokenRequirements(),
		m.focusIndex,
		m.width,
		m.height,
	)
}

func updateInputs(m *model) tea.Cmd {
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
