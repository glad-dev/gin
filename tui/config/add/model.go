package add

import (
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type displaying int

const (
	displayingAdd displaying = iota
	displayingError
)

type model struct {
	exitText            string
	error               string
	inputs              []textinput.Model
	currentlyDisplaying displaying
	focusIndex          int
	width               int
	height              int
	submit              bool
	quit                bool
	failure             bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quit = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.currentlyDisplaying {
	case displayingAdd:
		cmd = m.updateAdd(msg)

		return m, cmd

	case displayingError:
		m.updateError(msg)

		return m, nil

	default:
		m.exitText = style.FormatQuitText("Invalid update state")
		m.failure = true

		return m, tea.Quit
	}
}

func (m model) View() string {
	if m.quit {
		return style.FormatQuitText("No changes were made")
	}

	if m.submit {
		return m.exitText
	}

	switch m.currentlyDisplaying {
	case displayingAdd:
		return m.viewAdd()

	case displayingError:
		return m.viewError()

	default:
		return "Invalid view state"
	}
}
