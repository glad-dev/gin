package remove

import (
	"fmt"

	"gn/config"
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type displaying int

const (
	displayingList displaying = iota
	displayingDetails
	displayingEdit
	displayingError
)

type model struct {
	remotes             list.Model
	exitText            string
	oldConfig           config.Wrapper
	currentlyDisplaying displaying
	quitting            bool
	finished            bool
	failure             bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.InputField.GetFrameSize()
		m.remotes.SetSize(msg.Width-h, msg.Height-v)

		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			m.exitText, m.failure = onSubmit(&m)
			m.finished = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.remotes, cmd = m.remotes.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return style.FormatQuitText("No changes were made.")
	}

	if m.finished {
		return m.exitText
	}

	return shared.RenderList(m.remotes)
}

func onSubmit(m *model) (string, bool) {
	index := m.remotes.Index()

	selected, ok := m.remotes.Items()[index].(shared.ListItem)
	if !ok {
		return style.FormatQuitText("Failed to convert list.Item to item"), true
	}

	err := config.Remove(&m.oldConfig, index)
	if err != nil {
		return style.FormatQuitText(fmt.Sprintf("Failed to remove remote: %s", err)), true
	}

	return style.FormatQuitText(fmt.Sprintf("Sucessfully deleted the remote %s\nRemember to delete the API key on Gitlab", selected.Remote.URL.String())), false
}
