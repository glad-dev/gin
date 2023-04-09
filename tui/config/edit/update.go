package edit

import (
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) updateList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quit = true

			return tea.Quit

		case "enter":
			// User selected a config
			selected, ok := m.remotes.Items()[m.remotes.Index()].(editListItem)
			if !ok {
				m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")
				m.failure = true

				return tea.Quit
			}

			if len(selected.remote.Details) > 1 {
				m.currentlyDisplaying = displayingDetails

				items := make([]list.Item, len(selected.remote.Details))
				for i, details := range selected.remote.Details {
					items[i] = detail{
						username:  details.Username,
						tokenName: details.TokenName,
					}
				}

				m.details.SetItems(items)
				m.details.ResetSelected()

				return nil
			}

			match, err := selected.remote.ToMatch()
			if err != nil {
				m.exitText = style.FormatQuitText("Failed to convert item to match: " + err.Error())
				m.failure = true

				return tea.Quit
			}

			m.currentlyDisplaying = displayingEdit
			m.edit.Set(match, m.remotes.Index(), 0)

			return nil
		}
	}

	m.remotes, cmd = m.remotes.Update(msg)

	return cmd
}

func (m *model) updateDetails(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.currentlyDisplaying = displayingList
			m.details.SetItems([]list.Item{})

			return nil

		case "enter":
			selected, ok := m.remotes.Items()[m.remotes.Index()].(editListItem)
			if !ok {
				m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")
				m.failure = true

				return tea.Quit
			}

			match, err := selected.remote.ToMatchAtIndex(m.details.Index())
			if err != nil {
				m.exitText = style.FormatQuitText("Failed to convert item to match: " + err.Error())
				m.failure = true

				return tea.Quit
			}

			m.currentlyDisplaying = displayingEdit
			m.edit.Set(match, m.remotes.Index(), m.details.Index())

			return nil
		}
	}

	m.details, cmd = m.details.Update(msg)

	return cmd
}

func (m *model) updateEdit(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if len(m.details.Items()) > 0 {
				m.currentlyDisplaying = displayingDetails

				return nil
			}

			m.currentlyDisplaying = displayingList

			return nil

		case "enter":
			var errorStr string
			var failure bool

			errorStr, failure, cmd = m.edit.Update(msg)
			if !failure {
				return tea.Quit
			}

			m.currentlyDisplaying = displayingError
			m.error = lipgloss.Place(
				m.remotes.Width(),
				m.remotes.Height(),
				lipgloss.Center,
				lipgloss.Center,

				errorStr,
			)

			return cmd
		}
	}

	m.exitText, m.failure, cmd = m.edit.Update(msg)

	return cmd
}

func (m *model) updateError(msg tea.Msg) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "enter":
			m.currentlyDisplaying = displayingEdit

			return
		}
	}
}
