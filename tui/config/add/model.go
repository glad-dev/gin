package add

import (
	"fmt"
	"gn/config"
	"gn/tui/style/color"
	style "gn/tui/style/config"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(color.Focused)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	exitText        string
	inputFieldStyle lipgloss.Style
	inputs          []textinput.Model
	focusIndex      int
	width           int
	height          int
	submit          bool
	quit            bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quit = true

			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.submit = true
				m.exitText = onSubmit(&m)

				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle

					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quit {
		return style.QuitText.Render("No changes were made")
	}

	if m.submit {
		return m.exitText
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Gitlab URL",
				m.inputFieldStyle.Render(m.inputs[0].View()),
			),
			"\n",
			lipgloss.JoinVertical(
				lipgloss.Left,
				"API Key",
				m.inputFieldStyle.Render(m.inputs[1].View()),
			),
			"\n",
			lipgloss.JoinVertical(
				lipgloss.Left,
				*button,
			),
		),
	)
}

func onSubmit(m *model) string {
	err := config.Append(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		return style.QuitText.Render(fmt.Sprintf("Could not add config: %s", err))
	}

	return style.QuitText.Render(fmt.Sprintf("Successfully added config for %s", m.inputs[0].Value()))
}
