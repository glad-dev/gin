package remove

import (
	"fmt"

	"gn/config"
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	exitText  string
	list      list.Model
	oldConfig config.Wrapper
	quitting  bool
	finished  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		if msg.Height > 2*style.InputField.GetVerticalPadding() {
			// We have to subtract our padding
			m.list.SetHeight(msg.Height - 2*style.InputField.GetVerticalPadding())
		}

		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			m.exitText = onSubmit(&m)
			m.finished = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return style.FormatQuitText("No changes were made.")
	}

	if m.finished {
		return m.exitText
	}

	return shared.RenderList(m.list)
}

func onSubmit(m *model) string {
	index := m.list.Index()

	selected, ok := m.list.Items()[index].(shared.ListItem)
	if !ok {
		return style.FormatQuitText("Failed to convert list.Item to item")
	}

	err := config.Remove(&m.oldConfig, index)
	if err != nil {
		return style.FormatQuitText(fmt.Sprintf("Failed to remove remote: %s", err))
	}

	return style.FormatQuitText(fmt.Sprintf("Sucessfully deleted the remote %s\nRemember to delete the API key on Gitlab", selected.Lab.URL.String()))
}
