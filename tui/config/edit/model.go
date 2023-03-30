package edit

import (
	"errors"
	"fmt"

	"gn/config"
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	exitText       string
	inputs         []textinput.Model
	selectedConfig *config.GitLab
	list           list.Model
	oldConfig      config.Wrapper
	focusIndex     int
	quit           bool
	submit         bool
}

func (m model) displayingList() bool {
	return m.selectedConfig == nil
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
		if m.displayingList() {
			return updateList(&m, msg)
		}

		return updateSelection(&m, msg)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func updateList(m *model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "q", "esc", "ctrl+c":
		m.quit = true

		return *m, tea.Quit
	case "enter": //nolint:goconst
		selected, ok := m.list.Items()[m.list.Index()].(shared.ListItem)
		if !ok {
			m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")

			return *m, tea.Quit
		}

		m.selectedConfig = &selected.Lab
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return *m, cmd
}

func updateSelection(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quit = true

		return *m, tea.Quit

	case "esc":
		m.focusIndex = 0
		m.selectedConfig = nil
		// Delete entered values
		m.inputs[0].SetValue("")
		m.inputs[1].SetValue("")

		return updateFocus(m)

	// Set focus to next input
	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		// Did the user press enter while the submit button was focused?
		// If so, exit.
		if s == "enter" && m.focusIndex == len(m.inputs) {
			m.submit = true
			m.exitText = onSubmit(m)

			return *m, tea.Quit
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

		return updateFocus(m)
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return *m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quit {
		return style.FormatQuitText("No changes were made.")
	}

	if m.submit {
		return m.exitText
	}

	if m.selectedConfig == nil {
		// No config is selected => We must be in list view
		return shared.RenderList(m.list)
	}

	// There is a selected config => We must be in edit view
	if m.inputs[0].Value() == "" && m.inputs[1].Value() == "" {
		// Write the selected config's values in the input fields
		m.inputs[0].SetValue(m.selectedConfig.URL.String())
		m.inputs[1].SetValue(m.selectedConfig.Token)
	}

	return shared.RenderInputFields(
		m.inputs,
		m.focusIndex,
		m.list.Width(),
		m.list.Height()+2*style.InputField.GetVerticalPadding(),
	)
}

func updateFocus(m *model) (tea.Model, tea.Cmd) {
	var cmds = make([]tea.Cmd, len(m.inputs))
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

	return *m, tea.Batch(cmds...)
}

func onSubmit(m *model) string {
	err := config.Update(&m.oldConfig, m.list.Index(), m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		if errors.Is(err, config.ErrConfigDoesNotExist) {
			return style.FormatQuitText(config.ErrConfigDoesNotExistMsg)
		} else if errors.Is(err, config.ErrUpdateSameValues) {
			return style.FormatQuitText("No need to update the config: No changes were made.")
		}

		return style.FormatQuitText(fmt.Sprintf("Failed to update remote: %s", err))
	}

	return style.FormatQuitText(fmt.Sprintf("Sucessfully updated the remote %s", m.selectedConfig.URL.String()))
}
