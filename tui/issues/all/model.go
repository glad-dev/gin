package all

import (
	"fmt"

	"gn/config"
	"gn/issues"
	"gn/repo"
	shared "gn/tui/issues"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	lists        [3]list.Model
	shared       *shared.Shared
	activeTab    int
	viewedIssues map[string]issues.IssueDetails
	isLoading    bool
}

type updateMsg struct {
	projectPath string
	items       []itemWrapper
}

// itemWrapper is a wrapper for issues.Issue that implements all functions required by the list.Item interface.
type itemWrapper struct {
	issue issues.Issue
}

func (i itemWrapper) Title() string {
	return fmt.Sprintf(
		"#%s %s",
		i.issue.Iid,
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
	return tea.Batch(getIssues(m.shared.Details), m.shared.Spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() { //nolint:gocritic
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.isLoading {
				return m, nil
			}

			selected, ok := m.lists[m.activeTab].Items()[m.lists[m.activeTab].Index()].(itemWrapper)
			if !ok {
				return m, tea.Quit
			}

			m.shared.IssueID = selected.issue.Iid

		case "esc", "backspace":
			// If an issue is selected, deselect it.
			if len(m.shared.IssueID) != 0 {
				m.shared.IssueID = ""

				return m, nil
			}

			// Otherwise exit program.
			return m, tea.Quit

		case "right", "tab":
			m.activeTab = min(m.activeTab+1, len(m.lists)-1)
			return m, nil
		case "left", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil

		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()

		for i, _ := range m.lists {
			m.lists[i].SetSize(msg.Width-h, msg.Height-v)
		}

	case updateMsg:
		open := make([]list.Item, 0)
		closed := make([]list.Item, 0)
		all := make([]list.Item, 0)
		for _, item := range msg.items {
			switch item.issue.State {
			case "open", "opened":
				open = append(open, item)
				all = append(all, item)
			case "closed":
				closed = append(closed, item)

				item.issue.Title = "[closed] " + item.issue.Title
				all = append(all, item)
			}
		}

		// 0 => Open
		m.lists[0].SetItems(open)
		m.lists[0].Title = "Open issues"

		// 1 => Closed issues
		m.lists[1].SetItems(closed)
		m.lists[1].Title = "Closed issues"

		// 2 => All issues
		m.lists[2].SetItems(all)
		m.lists[2].Title = "All issues"

		m.isLoading = false

	default:
		m.shared.Spinner, cmd = m.shared.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.lists[m.activeTab], cmd = m.lists[m.activeTab].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.isLoading {
		return lipgloss.Place(
			m.lists[m.activeTab].Width(),
			m.lists[m.activeTab].Height(),
			lipgloss.Center,
			lipgloss.Center,

			fmt.Sprintf("Loading %s", m.shared.Spinner.View()),
		)
	}

	if len(m.shared.IssueID) > 0 {
		// Pull logic should be in update, not view but leaving it here for now until everything is connected.

		// Check if issue was requested in this session
		_, ok := m.viewedIssues[m.shared.IssueID]
		if ok {
			return docStyle.Render("I have already requested issue " + m.viewedIssues[m.shared.IssueID].Title)
		}

		m.viewedIssues[m.shared.IssueID] = issues.IssueDetails{
			Title: "Title: " + m.shared.IssueID,
		}
		return docStyle.Render("I want to request: ", m.viewedIssues[m.shared.IssueID].Title)
	}

	return docStyle.Render(m.lists[m.activeTab].View())
}

func getIssues(details []repo.Details) func() tea.Msg {
	return func() tea.Msg {
		conf, err := config.Load()
		if err != nil {
			style.PrintErrAndExit("Failed to load config: " + err.Error())
		}

		allIssues, projectPath, err := issues.QueryAll(conf, details)
		if err != nil {
			style.PrintErrAndExit("Failed to query issues: " + err.Error())
		}

		issueList := make([]itemWrapper, len(allIssues))
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
