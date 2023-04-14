package edit

import (
	"errors"
	"fmt"

	"gn/config"
	"gn/style"
	"gn/tui/config/shared"

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

func (m *editModel) set(match *config.Match, listIndex int, detailsIndex int) {
	// Set the new values
	m.inputs[0].SetValue(match.URL.String())
	m.inputs[1].SetValue(match.Token)
	m.listIndex = listIndex
	m.detailsIndex = detailsIndex

	// Set the focus to the first element
	m.focusIndex = 0
	m.updateFocus()
}

func (m *editModel) view() string {
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
			return config.ErrConfigDoesNotExistMsg, true
		} else if errors.Is(err, config.ErrUpdateSameValues) {
			return "No need to update the config: No changes were made.", false
		}

		return fmt.Sprintf("Failed to update remote: %s", err), true
	}

	return fmt.Sprintf("Sucessfully updated the remote %s", oldURL), false
}
