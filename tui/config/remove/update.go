package remove

import (
	"gn/style"
	"gn/tui/config/shared"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateRemoteList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quit = true

			return tea.Quit

		case "enter":
			// User selected a config
			selected, ok := m.remotes.Items()[m.remotes.Index()].(shared.ListItem)
			if !ok {
				m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")
				m.failure = true

				return tea.Quit
			}

			if len(selected.Remote.Details) > 1 {
				m.currentlyDisplaying = displayingDetails

				items := make([]list.Item, len(selected.Remote.Details))
				for i, details := range selected.Remote.Details {
					items[i] = shared.DetailItem{
						Username:  details.Username,
						TokenName: details.TokenName,
					}
				}

				m.details.SetItems(items)
				m.details.ResetSelected()

				return nil
			}

			m.currentlyDisplaying = displayingConfirmation
			m.details.SetItems([]list.Item{})
			m.details.ResetSelected()

			return nil
		}
	}

	m.remotes, cmd = m.remotes.Update(msg)

	return cmd
}

func (m *model) updateDetailsList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.currentlyDisplaying = displayingList
			m.details.SetItems([]list.Item{})

			return nil

		case "enter":
			m.currentlyDisplaying = displayingConfirmation
			m.confirmPosition = 0

			return nil
		}
	}

	m.details, cmd = m.details.Update(msg)

	return cmd
}

func (m *model) updateConfirmation(msg tea.Msg) tea.Cmd {
	goBack := func() {
		if len(m.details.Items()) == 0 {
			m.currentlyDisplaying = displayingList

			return
		}

		m.currentlyDisplaying = displayingDetails
	}

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			goBack()

			return nil

		case "right", "left", "tab", "shift+tab":
			if keypress == "right" || keypress == "tab" {
				m.confirmPosition++
				if m.confirmPosition > 1 {
					m.confirmPosition = 1
				}
			} else {
				m.confirmPosition--
				if m.confirmPosition < 0 {
					m.confirmPosition = 0
				}
			}

		case "enter":
			if m.confirmPosition == 0 {
				goBack()

				return nil
			}

			var failure bool
			m.exitText, failure = submit(m)
			if failure {
				m.failure = true
			} else {
				m.finished = true
			}

			return tea.Quit
		}
	}

	return nil
}
