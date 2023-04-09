package edit

import (
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
	details             list.Model
	exitText            string
	error               string
	edit                editModel
	currentlyDisplaying displaying
	quit                bool
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
		m.details.SetSize(msg.Width-h, msg.Height-v)

		m.edit.width = msg.Width - h
		m.edit.height = msg.Height - v

		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quit = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.currentlyDisplaying {
	case displayingList:
		cmd = m.updateList(msg)

		return m, cmd

	case displayingDetails:
		cmd = m.updateDetails(msg)

		return m, cmd

	case displayingEdit:
		cmd = m.updateEdit(msg)

		return m, cmd

	case displayingError:
		m.updateError(msg)

		return m, nil

	default:
		m.failure = true
		m.exitText = style.FormatQuitText("Invalid display type")

		return m, tea.Quit
	}
}

func (m model) View() string {
	if m.quit {
		return style.FormatQuitText("No changes were made.")
	}

	if len(m.exitText) > 0 {
		return m.exitText
	}

	switch m.currentlyDisplaying {
	case displayingList:
		return shared.RenderList(m.remotes)

	case displayingDetails:
		return shared.RenderList(m.details)

	case displayingEdit:
		return m.edit.View()

	case displayingError:
		return m.error

	default:
		return "Unkown view"
	}
}
