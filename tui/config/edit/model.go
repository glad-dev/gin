// nolint: goconst
package edit

import (
	"fmt"
	"io"

	"gn/config"
	tui "gn/tui/config"
	"gn/tui/style"
	"gn/tui/style/color"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	lab config.GitLab
}

func (i item) FilterValue() string { return "" }
func (i item) Title() string {
	return i.lab.URL.String()
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.lab.URL.String())

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	exitText       string
	inputs         []textinput.Model
	selectedConfig *config.GitLab
	list           list.Model
	oldConfig      config.Wrapper
	focusIndex     int
	quit           bool
	submit         bool
}

func (m model) displayingList() bool {
	return m.selectedConfig == nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	tmp := msg
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)

		return m, nil
	case tea.KeyMsg:
		if m.displayingList() {
			return updateList(&m, msg)
		}

		return updateSelection(&m, msg, tmp)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func updateList(m *model, msg tea.KeyMsg) (model, tea.Cmd) {
	switch keypress := msg.String(); keypress {
	case "q", "esc", "ctrl+c":
		m.quit = true

		return *m, tea.Quit
	case "enter":
		selected, ok := m.list.Items()[m.list.Index()].(item)
		if !ok {
			m.exitText = style.FormatQuitText("Failed to cast selected item to list.Item")

			return *m, tea.Quit
		}

		m.selectedConfig = &selected.lab
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return *m, cmd
}

var (
	focusedStyle = lipgloss.NewStyle().Foreground(color.Focused)
	noStyle      = lipgloss.NewStyle()
)

func updateSelection(m *model, msg tea.KeyMsg, orgMsg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg.String() {
	case "ctrl+c":
		m.quit = true

		return *m, tea.Quit

	case "esc":
		m.focusIndex = 0
		m.selectedConfig = nil
		// Delete entered values
		m.inputs[0].SetValue("")
		m.inputs[1].SetValue("")

		return updateFocus(m)

	// Set focus to next input
	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		// Did the user press enter while the submit button was focused?
		// If so, exit.
		if s == "enter" && m.focusIndex == len(m.inputs) {
			m.submit = true
			m.exitText = onSubmit(m)

			return *m, tea.Quit
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

		return updateFocus(m)
	}

	// Handle character input and blinking
	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(orgMsg)
	}

	return *m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quit {
		return style.FormatQuitText("No changes were made.")
	}

	if m.submit {
		return m.exitText
	}

	if m.selectedConfig != nil {
		if m.inputs[0].Value() == "" && m.inputs[1].Value() == "" {
			m.inputs[0].SetValue(m.selectedConfig.URL.String())
			m.inputs[1].SetValue(m.selectedConfig.Token)
		}

		return tui.RenderInputFields(m.inputs, m.focusIndex, m.list.Width(), m.list.Height())
	}

	return "\n" + m.list.View()
}

func updateFocus(m *model) (tea.Model, tea.Cmd) {
	var cmds = make([]tea.Cmd, len(m.inputs))
	for i := 0; i < len(m.inputs); i++ {
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

	return *m, tea.Batch(cmds...)
}

func onSubmit(m *model) string {
	err := config.Update(&m.oldConfig, m.list.Index(), m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		return style.FormatQuitText(fmt.Sprintf("Failed to update remote: %s", err))
	}

	return style.FormatQuitText(fmt.Sprintf("Sucessfully updated the remote %s", m.selectedConfig.URL.String()))
}
