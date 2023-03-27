package all

import (
	"fmt"
	"os"

	"gn/config"
	"gn/issues"
	"gn/repo"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	selectedIid  string
	list         list.Model
	details      []repo.Details
	viewedIssues map[string]issues.IssueDetails
	spinner      spinner.Model // ToDo: Spinner is not spinning
	isLoading    bool
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
	return tea.Batch(getIssues(m.details), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() { //nolint:gocritic
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.isLoading {
				return m, nil
			}

			selected, ok := m.list.Items()[m.list.Index()].(itemWrapper)
			if !ok {
				return m, tea.Quit
			}

			m.selectedIid = selected.issue.Iid

		case "esc", "backspace":
			// If an issue is selected, deselect it.
			if len(m.selectedIid) != 0 {
				m.selectedIid = ""

				return m, nil
			}

			// Otherwise exit program.
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case updateMsg:
		m.list.SetItems(msg.items)
		m.list.Title = "Issues of " + msg.projectPath
		m.isLoading = false

	default:
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.isLoading {
		return lipgloss.Place(
			m.list.Width(),
			m.list.Height(),
			lipgloss.Center,
			lipgloss.Center,

			fmt.Sprintf("Loading %s", m.spinner.View()),
		)
	}

	if len(m.selectedIid) > 0 {
		// Pull logic should be in update, not view but leaving it here for now until everything is connected.

		// Check if issue was requested in this session
		_, ok := m.viewedIssues[m.selectedIid]
		if ok {
			return docStyle.Render("I have already requested issue " + m.selectedIid)
		}

		return docStyle.Render("I want to request: ", m.selectedIid)
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
