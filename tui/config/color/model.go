package color

import (
	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/style"

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

// Init is required for model to be a tea.Model.
func (m model) Init() tea.Cmd {
	return nil
}

// Update is required for model to be a tea.Model.
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
		cmd = updateColor(&m, msg)

		return m, cmd

	case displayingError:
		updateError(&m, msg)

		return m, nil

	default:
		m.state = exitFailure
		m.text = "Invalid update state"

		return m, nil
	}
}

// View is required for model to be a tea.Model.
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
		return viewAdd(&m)

	case displayingError:
		return viewError(&m)

	default:
		return "Invalid view state"
	}
}
