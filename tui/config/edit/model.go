package edit

import (
	"fmt"

	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type displaying int

const (
	displayingList displaying = iota
	displayingDetails
	displayingEdit
	displayingError
)

var errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

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

		m.edit.width = msg.Width
		m.edit.height = msg.Height

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
		return m.edit.view()

	case displayingError:
		return lipgloss.Place(
			m.edit.width,
			m.edit.height,
			lipgloss.Center,
			0.75,

			fmt.Sprintf(
				"%s\n%s\n\nPress the 'q', 'esc' or 'enter' key to go back.",
				errorStyle.Render("An error occurred:"),
				m.error,
			),
		)

	default:
		return "Unkown view"
	}
}
