package add

import (
	"gn/style"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type displaying int8

const (
	displayingAdd displaying = iota
	displayingLoading
	displayingError
)

type model struct {
	text                string
	inputs              []textinput.Model
	spinner             spinner.Model
	focusIndex          int
	width               int
	height              int
	currentlyDisplaying displaying
	done                bool
	noChanges           bool
	failure             bool
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

type updateMsg struct {
	str     string
	failure bool
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.noChanges = true

			return m, tea.Quit
		}

	case updateMsg:
		m.text = msg.str
		if msg.failure {
			m.currentlyDisplaying = displayingError

			return m, nil
		}

		m.done = true
		m.text = "Test: " + m.text

		return m, tea.Quit
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
		m.failure = true

		return m, tea.Quit
	}
}

func (m model) View() string {
	if m.noChanges {
		return style.FormatQuitText("No changes were made")
	}

	if m.done {
		return "What?"
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
