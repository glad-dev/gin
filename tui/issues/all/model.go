package all

import (
	"fmt"

	"gn/config"
	"gn/issues"
	"gn/repo"
	shared "gn/tui/issues"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	shared       *shared.Shared
	viewedIssues map[string]issues.IssueDetails
	tabs         tabs
	view         view
	isLoading    bool
	viewingList  bool
}

type tabs struct {
	lists     [3]list.Model
	activeTab int
}

type view struct {
	content string
	view    viewport.Model
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
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()

		for i := range m.tabs.lists {
			m.tabs.lists[i].SetSize(msg.Width-h, msg.Height-v-8)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.isLoading {
				return m, nil
			}

			selected, ok := m.tabs.lists[m.tabs.activeTab].Items()[m.tabs.lists[m.tabs.activeTab].Index()].(itemWrapper)
			if !ok {
				return m, tea.Quit
			}

			m.shared.IssueID = selected.issue.Iid

			_, ok = m.viewedIssues[m.shared.IssueID]
			if !ok {
				m.viewedIssues[m.shared.IssueID] = issues.IssueDetails{
					Title: "Title: " + m.shared.IssueID,
				}
			}

		case "esc", "backspace":
			// If an issue is selected, deselect it.
			if len(m.shared.IssueID) != 0 {
				m.shared.IssueID = ""

				return m, nil
			}

			// Otherwise exit program.
			return m, tea.Quit

		case "right", "tab":
			m.tabs.activeTab = min(m.tabs.activeTab+1, len(m.tabs.lists)-1)

			return m, nil
		case "left", "shift+tab":
			m.tabs.activeTab = max(m.tabs.activeTab-1, 0)

			return m, nil
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
		m.tabs.lists[0].SetItems(open)
		m.tabs.lists[0].Title = "Open issues"

		// 1 => Closed issues
		m.tabs.lists[1].SetItems(closed)
		m.tabs.lists[1].Title = "Closed issues"

		// 2 => All issues
		m.tabs.lists[2].SetItems(all)
		m.tabs.lists[2].Title = "All issues"

		m.isLoading = false

	default:
		m.shared.Spinner, cmd = m.shared.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.tabs.lists[m.tabs.activeTab], cmd = m.tabs.lists[m.tabs.activeTab].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.isLoading {
		return lipgloss.Place(
			m.tabs.lists[m.tabs.activeTab].Width(),
			m.tabs.lists[m.tabs.activeTab].Height(),
			lipgloss.Center,
			lipgloss.Center,

			fmt.Sprintf("Loading %s", m.shared.Spinner.View()),
		)
	}

	if len(m.shared.IssueID) > 0 {
		return docStyle.Render("I want to request: ", m.viewedIssues[m.shared.IssueID].Title)
	}

	if m.viewingList {
		return renderTab(&m.tabs)
	}

	return renderViewport(&m)
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
