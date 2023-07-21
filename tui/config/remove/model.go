package remove

import (
	"fmt"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/config/shared"

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
	displayingConfirmation
)

const (
	stateRunning state = iota
	exitSuccess
	exitFailure
	exitNoChange
)

type model struct {
	remotes             list.Model
	details             list.Model
	text                string
	oldConfig           config.Wrapper
	confirmPosition     int
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
		h, v := style.List.GetFrameSize()
		m.remotes.SetSize(msg.Width-h, msg.Height-v)
		m.details.SetSize(msg.Width-h, msg.Height-v)

		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.state = exitNoChange

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.currentlyDisplaying {
	case displayingList:
		cmd = updateRemoteList(&m, msg)

		return m, cmd

	case displayingDetails:
		cmd = updateDetailsList(&m, msg)

		return m, cmd

	case displayingConfirmation:
		cmd = updateConfirmation(&m, msg)

		return m, cmd

	default:
		m.state = exitFailure
		m.text = "Invalid display state"

		return m, tea.Quit
	}
}

// View required for model to be a tea.Model.
func (m model) View() string {
	switch m.state {
	case stateRunning:
		break

	case exitNoChange:
		return style.FormatQuitText("No changes were made.")

	case exitFailure, exitSuccess:
		return m.text
	}

	switch m.currentlyDisplaying {
	case displayingList:
		return shared.RenderList(m.remotes)

	case displayingDetails:
		return shared.RenderList(m.details)

	case displayingConfirmation:
		return viewConfirmation(&m)

	default:
		return "Invalid state!"
	}
}

func submit(m *model) (string, bool) {
	selected, ok := m.remotes.SelectedItem().(shared.ListItem)
	if !ok {
		return style.FormatQuitText("Failed to convert list.Item to item"), true
	}

	tokenName := selected.Remote.Details[m.details.Index()].TokenName

	err := config.Remove(&m.oldConfig, m.remotes.Index(), m.details.Index())
	if err != nil {
		return style.FormatQuitText(fmt.Sprintf("Failed to remove remote: %s", err)), true
	}

	return style.FormatQuitText(fmt.Sprintf("Sucessfully deleted the token '%s'\nRemember to delete the API key", tokenName)), false
}
