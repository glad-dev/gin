package remove

import (
	"fmt"

	"gn/config"
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type displaying int

const (
	displayingList displaying = iota
	displayingDetails
	displayingConfirmation
)

type model struct {
	remotes             list.Model
	details             list.Model
	exitText            string
	oldConfig           config.Wrapper
	currentlyDisplaying displaying
	confirmPosition     int
	quit                bool
	finished            bool
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
		cmd = m.updateRemoteList(msg)

		return m, cmd

	case displayingDetails:
		cmd = m.updateDetailsList(msg)

		return m, cmd

	case displayingConfirmation:
		cmd = m.updateConfirmation(msg)

		return m, cmd

	default:
		m.failure = true
		m.exitText = "Invalid display state"

		return m, tea.Quit
	}
}

func (m model) View() string {
	if m.quit {
		return style.FormatQuitText("No changes were made.")
	}

	if m.finished || m.failure {
		return m.exitText
	}

	switch m.currentlyDisplaying {
	case displayingList:
		return shared.RenderList(m.remotes)

	case displayingDetails:
		return shared.RenderList(m.details)

	case displayingConfirmation:
		return m.viewConfirmation()

	default:
		return "Invalid state!"
	}
}

func submit(m *model) (string, bool) {
	selected, ok := m.remotes.Items()[m.remotes.Index()].(shared.ListItem)
	if !ok {
		return style.FormatQuitText("Failed to convert list.Item to item"), true
	}

	tokenName := selected.Remote.Details[m.details.Index()].TokenName

	err := config.Remove(&m.oldConfig, m.remotes.Index(), m.details.Index())
	if err != nil {
		return style.FormatQuitText(fmt.Sprintf("Failed to remove remote: %s", err)), true
	}

	return style.FormatQuitText(fmt.Sprintf("Sucessfully deleted the token '%s'\nRemember to delete the API key on Gitlab", tokenName)), false
}
