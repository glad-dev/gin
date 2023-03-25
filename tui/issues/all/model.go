package all

import (
	"fmt"

	"gn/issues"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list list.Model
}

// item is a wrapper for issues.Issue that implements all functions required by the list.Item interface.
type item struct {
	issue issues.Issue
}

func (i item) Title() string {
	status := ""
	if i.issue.State == "closed" {
		status = "[closed] "
	}

	return fmt.Sprintf(
		"#%s %s%s by %s on %s",
		i.issue.Iid, status,
		i.issue.Title,
		i.issue.Author.String(),
		i.issue.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func (i item) Description() string {
	return i.issue.Description
}

func (i item) FilterValue() string {
	return i.issue.Title + i.issue.Description
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() { //nolint:gocritic
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}
