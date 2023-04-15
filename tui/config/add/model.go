package add

import (
	"gn/style"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	displaying uint8
	state      uint8
)

const (
	displayingAdd displaying = iota
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
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

type updateMsg struct {
	str     string
	success bool
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
	case displayingAdd:
		cmds[1] = m.updateAdd(msg)

		return m, tea.Batch(cmds...)

	case displayingLoading:
		return m, tea.Batch(
			cmds[0],
			m.updateLoading(),
		)

	case displayingError:
		m.updateError(msg)

		return m, cmds[0]

	default:
		m.text = style.FormatQuitText("Invalid update state")
		m.state = exitFailure

		return m, tea.Quit
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
	case displayingAdd:
		return m.viewAdd()

	case displayingLoading:
		return m.viewLoading()

	case displayingError:
		return m.viewError()

	default:
		return "Invalid view state"
	}
}
