package color

import (
	"gn/config"
	"gn/style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	displaying uint8
	state      uint8
)

const (
	displayingColors displaying = iota
	displayingError
)

const (
	stateRunning state = iota
	exitNoChange
	exitFailure
	exitSuccess
)

type model struct {
	wrapper             *config.Wrapper
	text                string
	inputs              []textinput.Model
	width               int
	height              int
	focusIndex          int
	currentlyDisplaying displaying
	state               state
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
			m.state = exitNoChange

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.currentlyDisplaying {
	case displayingColors:
		cmd = m.updateColor(msg)

		return m, cmd

	case displayingError:
		m.updateError(msg)

		return m, nil

	default:
		m.state = exitFailure
		m.text = "Invalid update state"

		return m, nil
	}
}

func (m model) View() string {
	switch m.state {
	case stateRunning:
		break

	case exitNoChange:
		return style.FormatQuitText("No changes were made")

	case exitFailure, exitSuccess:
		return m.text
	}

	switch m.currentlyDisplaying {
	case displayingColors:
		return m.viewAdd()

	case displayingError:
		return m.viewError()

	default:
		return "Invalid view state"
	}
}
