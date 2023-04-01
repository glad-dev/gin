package edit

import (
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list                list.Model
	exitText            string
	edit                editModel
	currentlyDisplaying displaying
	quit                bool
	failure             bool
}

type displaying int

const (
	displayingList displaying = iota
	displayingEdit
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tmp *ret

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.InputField.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.edit.width = msg.Width
		m.edit.height = msg.Height

		return m, nil
	case tea.KeyMsg:
		// q should only quit the program if we're in list view. Otherwise, the user would be unable to enter a URL
		// or token that contains the letter q
		if m.currentlyDisplaying == displayingList && msg.String() == "q" {
			m.quit = true

			return m, tea.Quit
		}

		switch msg.String() {
		case "ctrl+c":
			m.quit = true

			return m, tea.Quit
		case "esc":
			if m.currentlyDisplaying == displayingList {
				m.quit = true

				return m, tea.Quit
			}

			// We are currently displaying the edit view => Move back to list view
			m.currentlyDisplaying = displayingList

			return m, nil
		case "enter":
			if m.currentlyDisplaying == displayingList {
				// User selected a config
				selected, ok := m.list.Items()[m.list.Index()].(shared.ListItem)
				if !ok {
					m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")

					return m, tea.Quit
				}

				m.currentlyDisplaying = displayingEdit
				m.edit.Set(&selected.Lab, m.list.Index())

				return m, nil
			}

			tmp = m.edit.Update(msg)
			m.exitText = tmp.str
			m.failure = tmp.failure

			return m, tmp.cmd
		}
	}

	if m.currentlyDisplaying == displayingList {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		return m, cmd
	}

	tmp = m.edit.Update(msg)
	m.exitText = tmp.str

	return m, tmp.cmd
}

func (m model) View() string {
	if m.quit {
		return style.FormatQuitText("No changes were made.")
	}

	if len(m.exitText) > 0 {
		return m.exitText
	}

	if m.currentlyDisplaying == displayingList {
		return shared.RenderList(m.list)
	}

	return m.edit.View()
}
