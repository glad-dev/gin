package config

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
	view      func(m *model) string
	list      list.Model
	oldConfig config.General
	quitting  bool
	action    bool
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
			m.action = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return m.view(&m)
}
