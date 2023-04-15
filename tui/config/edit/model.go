package edit

import (
	"gn/style"
	"gn/tui/config/shared"

	"github.com/charmbracelet/bubbles/spinner"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	displaying uint8
	state      uint8
)

const (
	displayingList displaying = iota
	displayingDetails
	displayingEdit
	displayingLoading
	displayingError
)

const (
	stateRunning state = iota
	exitFailure
	exitSuccess
	exitNoChange
)

type model struct {
	remotes             list.Model
	details             list.Model
	text                string
	spinner             spinner.Model
	edit                editModel
	currentlyDisplaying displaying
	state               state
}

type updateMsg struct {
	str     string
	success bool
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
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
	case displayingList:
		cmds[1] = m.updateList(msg)

		return m, tea.Batch(cmds...)

	case displayingDetails:
		cmds[1] = m.updateDetails(msg)

		return m, tea.Batch(cmds...)

	case displayingEdit:
		cmds[1] = m.updateEdit(msg)

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
		m.state = exitFailure
		m.text = style.FormatQuitText("Invalid display type")

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
	case displayingList:
		return shared.RenderList(m.remotes)

	case displayingDetails:
		return shared.RenderList(m.details)

	case displayingEdit:
		return m.edit.view()

	case displayingLoading:
		return m.viewLoading()

	case displayingError:
		return m.viewError()

	default:
		return "Unkown view"
	}
}
