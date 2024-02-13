package add

import (
	rt "github.com/glad-dev/gin/remote/type"
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	displaying uint8
	state      uint8
)

const (
	displayingType displaying = iota
	displayingAdd
	displayingLoading
	displayingError
)

const (
	stateRunning state = iota
	exitNoChange
	exitFailure
	exitSuccess
)

type model struct {
	text                string
	inputs              []textinput.Model
	spinner             spinner.Model
	focusIndex          int
	width               int
	height              int
	currentlyDisplaying displaying
	state               state
	remoteType          rt.Type
}

type updateMsg struct {
	str     string
	success bool
}

// Init is required for model to be a tea.Model.
func (m model) Init() tea.Cmd {
	return m.spinner.Tick
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

	case updateMsg:
		m.text = msg.str
		if msg.success {
			m.state = exitSuccess

			return m, tea.Quit
		}

		m.currentlyDisplaying = displayingError

		return m, nil
	}

	cmds := make([]tea.Cmd, 2)
	m.spinner, cmds[0] = m.spinner.Update(msg)
	switch m.currentlyDisplaying {
	case displayingType:
		cmds[1] = updateType(&m, msg)

		return m, tea.Batch(cmds...)

	case displayingAdd:
		cmds[1] = updateAdd(&m, msg)

		return m, tea.Batch(cmds...)

	case displayingLoading:
		return m, tea.Batch(
			cmds[0],
			updateLoading(&m),
		)

	case displayingError:
		updateError(&m, msg)

		return m, cmds[0]

	default:
		m.text = style.FormatQuitText("Invalid update state")
		m.state = exitFailure

		return m, tea.Quit
	}
}

// View required for model to be a tea.Model.
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
	case displayingType:
		return viewType(&m)

	case displayingAdd:
		return viewAdd(&m)

	case displayingLoading:
		return viewLoading(&m)

	case displayingError:
		return viewError(&m)

	default:
		return "Invalid view state"
	}
}
