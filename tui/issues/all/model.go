package all

import (
	"fmt"
	"os"

	"gn/config"
	"gn/issues"
	"gn/repo"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	selectedIid string
	list        list.Model
	details     []repo.Details
}

type updateMsg struct {
	projectPath string
	items       []list.Item
}

// itemWrapper is a wrapper for issues.Issue that implements all functions required by the list.Item interface.
type itemWrapper struct {
	issue issues.Issue
}

func (i itemWrapper) Title() string {
	status := ""
	if i.issue.State == "closed" {
		status = "[closed] "
	}

	return fmt.Sprintf(
		"#%s %s%s",
		i.issue.Iid,
		status,
		i.issue.Title,
	)
}

func (i itemWrapper) Description() string {
	// Use author and creation date as description
	return fmt.Sprintf(
		"Created by %s on %s",
		i.issue.Author.String(),
		i.issue.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func (i itemWrapper) FilterValue() string {
	return i.issue.Title
}

func (m model) Init() tea.Cmd {
	return getIssues(m.details)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() { //nolint:gocritic
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected, ok := m.list.Items()[m.list.Index()].(itemWrapper)
			if !ok {
				return m, tea.Quit
			}

			m.selectedIid = selected.issue.Iid
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.Msg:
		val, ok := msg.(updateMsg)
		if !ok {
			return m, nil
		}

		m.list.SetItems(val.items)
		m.list.Title = "Issues of " + val.projectPath
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if len(m.list.Title) == 0 {
		return "Loading..."
	}

	if m.selectedIid != "" {
		return docStyle.Render("I want: ", m.selectedIid)
	}

	return docStyle.Render(m.list.View())
}

func getIssues(details []repo.Details) func() tea.Msg {
	return func() tea.Msg {
		conf, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
			os.Exit(1)
		}

		allIssues, projectPath, err := issues.QueryAll(conf, details)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failure: %s\n", err)
			os.Exit(1)
		}

		issueList := make([]list.Item, len(allIssues))
		for i, issue := range allIssues {
			issueList[i] = itemWrapper{
				issue: issue,
			}
		}

		return updateMsg{
			projectPath: projectPath,
			items:       issueList,
		}
	}
}
