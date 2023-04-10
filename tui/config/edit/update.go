package edit

import (
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

			match, err := selected.Remote.ToMatch()
			if err != nil {
				m.exitText = style.FormatQuitText("Failed to convert item to match: " + err.Error())
				m.failure = true

				return tea.Quit
			}

			m.currentlyDisplaying = displayingEdit
			m.edit.set(match, m.remotes.Index(), 0)

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
			selected, ok := m.remotes.Items()[m.remotes.Index()].(shared.ListItem)
			if !ok {
				m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")
				m.failure = true

				return tea.Quit
			}

			match, err := selected.Remote.ToMatchAtIndex(m.details.Index())
			if err != nil {
				m.exitText = style.FormatQuitText("Failed to convert item to match: " + err.Error())
				m.failure = true

				return tea.Quit
			}

			m.currentlyDisplaying = displayingEdit
			m.edit.set(match, m.remotes.Index(), m.details.Index())

			return nil
		}
	}

	m.details, cmd = m.details.Update(msg)

	return cmd
}

func (m *model) updateEdit(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch s := msg.String(); s {
		case "esc":
			if len(m.details.Items()) == 0 {
				m.currentlyDisplaying = displayingList

				return nil
			}

			m.currentlyDisplaying = displayingDetails

			return nil

		case "tab", "shift+tab", "enter", "up", "down": //nolint: goconst
			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.edit.focusIndex == len(m.edit.inputs) {
				str, failure := m.edit.submit()
				if !failure {
					m.exitText = style.FormatQuitText(str)

					return tea.Quit
				}

				m.error = str
				m.currentlyDisplaying = displayingError

				return nil
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.edit.focusIndex--
			} else {
				m.edit.focusIndex++
			}

			if m.edit.focusIndex > len(m.edit.inputs) {
				m.edit.focusIndex = 0
			} else if m.edit.focusIndex < 0 {
				m.edit.focusIndex = len(m.edit.inputs)
			}

			return m.edit.updateFocus()
		}
	}

	cmds := make([]tea.Cmd, len(m.edit.inputs))
	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.edit.inputs {
		m.edit.inputs[i], cmds[i] = m.edit.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) updateError(msg tea.Msg) {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "backspace":
			m.currentlyDisplaying = displayingEdit

			return
		}
	}
}
