package edit

import (
	"errors"
	"fmt"

	"gn/config"
	"gn/tui/config/shared"
	"gn/tui/style"

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

func (m *editModel) Update(msg tea.Msg) (string, bool, tea.Cmd) {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		s := msg.String()
		switch s {
		case "tab", "shift+tab", "enter", "up", "down": //nolint: goconst
			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				str, failure := m.submit()

				return str, failure, nil
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

			return "", false, m.updateFocus()
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return "", false, tea.Batch(cmds...)
}

func (m *editModel) Set(match *config.Match, listIndex int, detailsIndex int) {
	// Set the new values
	m.inputs[0].SetValue(match.URL.String())
	m.inputs[1].SetValue(match.Token)
	m.listIndex = listIndex
	m.detailsIndex = detailsIndex

	// Set the focus to the first element
	m.focusIndex = 0
	m.updateFocus()
}

func (m *editModel) View() string {
	return shared.RenderInputFields(
		m.inputs,
		m.focusIndex,
		m.width,
		m.height-2*style.InputField.GetVerticalFrameSize(),
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

// submit returns exit text and if the operation failed.
func (m *editModel) submit() (string, bool) {
	oldURL := m.oldConfig.Configs[m.listIndex].URL.String()

	err := config.Update(m.oldConfig, m.listIndex, m.detailsIndex, m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		if errors.Is(err, config.ErrConfigDoesNotExist) {
			return style.FormatQuitText(config.ErrConfigDoesNotExistMsg), true
		} else if errors.Is(err, config.ErrUpdateSameValues) {
			return style.FormatQuitText("No need to update the config: No changes were made."), false
		}

		return style.FormatQuitText(fmt.Sprintf("Failed to update remote: %s", err)), true
	}

	return style.FormatQuitText(fmt.Sprintf("Sucessfully updated the remote %s", oldURL)), false
}
