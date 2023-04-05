package all

import (
	"gn/issues"
	"gn/tui/issues/shared"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func updateList(m *model, msg *updateMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Set config
	m.conf = msg.conf

	// Set lists
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

	m.tabs.lists[m.tabs.activeTab], cmd = m.tabs.lists[m.tabs.activeTab].Update(msg)

	return m, cmd
}

func handleListUpdate(m *model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "enter":
		if m.isLoading {
			return m, nil
		}

		m.isLoading = true // Does this do anything?
		pullIssue(m)
		m.isLoading = false

		// Re-set viewport view to the top
		m.viewport.GotoTop()

		if m.failure {
			return m, tea.Quit
		}

	case "esc":
		return m, tea.Quit

	case "right", "tab":
		m.tabs.activeTab = min(m.tabs.activeTab+1, len(m.tabs.lists)-1)

		return m, nil

	case "left", "shift+tab":
		m.tabs.activeTab = max(m.tabs.activeTab-1, 0)

		return m, nil
	}

	m.tabs.lists[m.tabs.activeTab], cmd = m.tabs.lists[m.tabs.activeTab].Update(msg)

	return m, cmd
}

func pullIssue(m *model) {
	selected, ok := m.tabs.lists[m.tabs.activeTab].Items()[m.tabs.lists[m.tabs.activeTab].Index()].(itemWrapper)
	if !ok {
		m.error = "Failed to convert selected item to itemWrapper"
		m.failure = true

		return
	}

	m.shared.IssueID = selected.issue.Iid

	issue, ok := m.viewedIssues[m.shared.IssueID]
	if !ok {
		// Request issue
		tmp, err := issues.QuerySingle(m.conf, m.shared.Details, m.shared.URL, selected.issue.Iid)
		if err != nil {
			m.error = "Failed to query issue: " + err.Error()
			m.failure = true

			return
		}

		// Store issue
		m.viewedIssues[m.shared.IssueID] = *tmp

		// Needed since issue.QuerySingle returns a pointer and map access a struct
		issue = *tmp
	}

	m.viewport.SetContent(shared.PrettyPrintIssue(&issue, m.viewport.Width, m.viewport.Height))
	m.viewingList = false
}
