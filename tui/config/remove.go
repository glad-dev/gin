package config

import (
	"fmt"
	"gn/config"
	"gn/tui/style"
	"io"
	"os"

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

	fn := style.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list      list.Model
	oldConfig config.General
	quitting  bool
	delete    bool
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
			m.delete = true

			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return style.QuitTextStyle.Render("No changes were made.")
	}

	if m.delete {
		index := m.list.Index()

		selected, ok := m.list.Items()[index].(item)
		if !ok {
			return style.QuitTextStyle.Render("Failed to convert list.Item to item")
		}

		err := config.RemoveOnce(&m.oldConfig, index) // This is called multiple times?
		if err != nil {
			return style.QuitTextStyle.Render(fmt.Sprintf("Failed to remove remote: %s", err))
		}

		return style.QuitTextStyle.Render(fmt.Sprintf("Sucessfully deleted the remote %s\nRemember to delete the API key on Gitlab", selected.lab.URL.String()))
	}

	return "\n" + m.list.View()
}

func Remove() {
	// Load current config
	generalConfig, err := config.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, style.QuitTextStyle.Render(fmt.Sprintf("Failed to load config: %s", err)))

		return
	}

	if len(generalConfig.Configs) == 0 {
		fmt.Fprint(os.Stderr, style.QuitTextStyle.Render("The config file contains no remotes."))

		return
	}

	items := make([]list.Item, len(generalConfig.Configs))
	for i, conf := range generalConfig.Configs {
		items[i] = item{
			lab: conf,
		}
	}

	// Not sure what these numbers mean, but the TUI looks better with them
	const defaultWidth = 20
	const listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Which remote do you want to delete?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = style.TitleStyle
	l.Styles.PaginationStyle = style.PaginationStyle
	l.Styles.HelpStyle = style.HelpStyle

	m := model{
		list:      l,
		oldConfig: *generalConfig,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
