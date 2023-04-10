package add

import (
	"errors"
	"fmt"

	"gn/tui/config/shared"

	"gn/config"
	"gn/tui/style"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateAdd(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.quit = true

			return tea.Quit

			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				str, failure := submit(m)
				if failure {
					m.error = str
					m.currentlyDisplaying = displayingError

					return nil
				}

				m.submit = true
				m.exitText = style.FormatQuitText(str)

				return tea.Quit
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
	}

	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func submit(m *model) (string, bool) {
	err := config.Append(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		if errors.Is(err, config.ErrConfigDoesNotExist) {
			return config.ErrConfigDoesNotExistMsg, true
		}

		return fmt.Sprintf("Could not add config: %s", err), true
	}

	return fmt.Sprintf("Successfully added config for %s", m.inputs[0].Value()), false
}

func (m *model) viewAdd() string {
	return shared.RenderInputFields(
		m.inputs,
		m.focusIndex,
		m.width,
		m.height,
	)
}
