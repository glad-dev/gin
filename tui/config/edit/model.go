package edit

import (
	"gn/tui/config/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type displaying int

const (
	displayingList displaying = iota
	displayingDetails
	displayingEdit
)

type model struct {
	remotes             list.Model
	details             list.Model
	exitText            string
	edit                editModel
	currentlyDisplaying displaying
	quit                bool
	failure             bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tmp *ret

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.InputField.GetFrameSize()
		m.remotes.SetSize(msg.Width-h, msg.Height-v)
		m.details.SetSize(msg.Width-h, msg.Height-v)

		m.edit.width = msg.Width - h
		m.edit.height = msg.Height - v

		return m, nil
	case tea.KeyMsg:
		// q should only quit the program if we're in list view. Otherwise, the user would be unable to enter a URL
		// or token that contains the letter q
		if m.currentlyDisplaying == displayingList && msg.String() == "q" {
			m.quit = true

			return m, tea.Quit
		}

		switch msg.String() {
		case "ctrl+c":
			m.quit = true

			return m, tea.Quit

		case "q":
			if m.currentlyDisplaying == displayingDetails {
				m.currentlyDisplaying = displayingList

				return m, nil
			}

		case "esc":
			if m.currentlyDisplaying == displayingList {
				m.quit = true

				return m, tea.Quit
			}

			m.currentlyDisplaying--

			return m, nil
		case "enter":
			switch m.currentlyDisplaying {
			case displayingList:
				// User selected a config
				selected, ok := m.remotes.Items()[m.remotes.Index()].(editListItem)
				if !ok {
					m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")
					m.failure = true

					return m, tea.Quit
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

					return m, nil
				}

				match, err := selected.remote.ToMatch()
				if err != nil {
					m.exitText = style.FormatQuitText("Failed to convert item to match: " + err.Error())
					m.failure = true

					return m, tea.Quit
				}

				m.currentlyDisplaying = displayingEdit
				m.edit.Set(match, m.remotes.Index(), 0)

				return m, nil

			case displayingDetails:
				selected, ok := m.remotes.Items()[m.remotes.Index()].(editListItem)
				if !ok {
					m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")
					m.failure = true

					return m, tea.Quit
				}

				match, err := selected.remote.ToMatchAtIndex(m.details.Index())
				if err != nil {
					m.exitText = style.FormatQuitText("Failed to convert item to match: " + err.Error())
					m.failure = true

					return m, tea.Quit
				}

				m.currentlyDisplaying = displayingEdit
				m.edit.Set(match, m.remotes.Index(), m.details.Index())

				return m, nil

			case displayingEdit:
				tmp = m.edit.Update(msg)
				m.exitText = tmp.str
				m.failure = tmp.failure

				return m, tmp.cmd

			default:
				m.failure = true
				m.exitText = style.FormatQuitText("Invalid displaying item.")

				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	switch m.currentlyDisplaying {
	case displayingList:
		m.remotes, cmd = m.remotes.Update(msg)

		return m, cmd

	case displayingDetails:
		m.details, cmd = m.details.Update(msg)

		return m, cmd

	case displayingEdit:
		tmp = m.edit.Update(msg)
		m.exitText = tmp.str

		return m, tmp.cmd

	default:
		m.failure = true
		m.exitText = style.FormatQuitText("Invalid dispaly type")

		return m, tea.Quit
	}
}

func (m model) View() string {
	if m.quit {
		return style.FormatQuitText("No changes were made.")
	}

	if len(m.exitText) > 0 {
		return m.exitText
	}

	switch m.currentlyDisplaying {
	case displayingList:
		return shared.RenderList(m.remotes)

	case displayingDetails:
		return shared.RenderList(m.details)

	case displayingEdit:
		return m.edit.View()

	default:
		return "Unkown view"
	}
}
