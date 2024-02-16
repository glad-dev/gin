package color

import (
	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/style"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func updateColor(m *model, msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = exitNoChange

			return tea.Quit

			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return submit(m)
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

func viewAdd(m *model) string {
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Left,

			"Blurred color",
			style.InputField.Render(m.inputs[0].View()),
			"\n",
			"Border color",
			style.InputField.Render(m.inputs[1].View()),
			"\n",
			"Focused color",
			style.InputField.Render(m.inputs[2].View()),
			"\n",
			*button,
		),
	)
}

func submit(m *model) tea.Cmd {
	c := configuration.Colors{
		Blurred: m.inputs[0].Value(),
		Border:  m.inputs[1].Value(),
		Focused: m.inputs[2].Value(),
	}

	err := c.CheckValidity()
	if err != nil {
		m.text = style.FormatQuitText("Invalid color passed: " + err.Error())
		m.currentlyDisplaying = displayingError

		return nil
	}

	err = configuration.UpdateColors(c, m.config)
	if err != nil {
		m.text = style.FormatQuitText("Failed to write config: " + err.Error())
		m.text = err.Error()
		m.currentlyDisplaying = displayingError

		return nil
	}

	m.state = exitSuccess
	m.text = style.FormatQuitText("Successfully updated the colors.")

	return tea.Quit
}
