package edit

import (
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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

// Init is required for model to be a tea.Model.
func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update is required for model to be a tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.InputField.GetFrameSize()
		m.remotes.SetSize(msg.Width-h-2, msg.Height-v-2)
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
		cmds[1] = updateList(&m, msg)

		return m, tea.Batch(cmds...)

	case displayingDetails:
		cmds[1] = updateDetails(&m, msg)

		return m, tea.Batch(cmds...)

	case displayingEdit:
		cmds[1] = updateEdit(&m, msg)

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
		m.state = exitFailure
		m.text = style.FormatQuitText("Invalid display type")

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
	case displayingList:
		return shared.RenderList(m.remotes)

	case displayingDetails:
		return shared.RenderList(m.details)

	case displayingEdit:
		return m.edit.view()

	case displayingLoading:
		return viewLoading(&m)

	case displayingError:
		return viewError(&m)

	default:
		return "Unknown view"
	}
}
