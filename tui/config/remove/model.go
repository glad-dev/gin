package remove

import (
	"fmt"
	"io"

	"gn/config"
	style "gn/tui/style/config"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
	exitText  string
	list      list.Model
	oldConfig config.Wrapper
	quitting  bool
	finished  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)

		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":
			fallthrough
		case "ctrl+c":
			m.quitting = true

			return m, tea.Quit

		case "enter":
			m.exitText = onSubmit(&m)
			m.finished = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return style.QuitText.Render("No changes were made.")
	}

	if m.finished {
		return m.exitText
	}

	return "\n" + m.list.View()
}

func onSubmit(m *model) string {
	index := m.list.Index()

	selected, ok := m.list.Items()[index].(item)
	if !ok {
		return style.QuitText.Render("Failed to convert list.Item to item")
	}

	err := config.Remove(&m.oldConfig, index)
	if err != nil {
		return style.QuitText.Render(fmt.Sprintf("Failed to remove remote: %s", err))
	}

	return style.QuitText.Render(fmt.Sprintf("Sucessfully deleted the remote %s\nRemember to delete the API key on Gitlab", selected.lab.URL.String()))
}
