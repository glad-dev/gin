package all

import (
	"strings"

	"gn/issues"
	"gn/logger"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) initList(msg *allIssuesUpdateMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Set lists
	open := make([]list.Item, 0)
	closed := make([]list.Item, 0)
	all := make([]list.Item, 0)
	for _, item := range msg.items {
		switch strings.ToLower(item.issue.State) {
		case "open", "opened":
			open = append(open, item)
			all = append(all, item)
		case "closed":
			closed = append(closed, item)
			item.issue.Title = "[closed] " + item.issue.Title
			all = append(all, item)
		default:
			logger.Log.Warnf("Got item with unkown state: %s", item.issue.State)
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

	m.currentlyDisplaying = displayingList

	m.tabs.lists[m.tabs.activeTab], cmd = m.tabs.lists[m.tabs.activeTab].Update(msg)

	return m, cmd
}

func (m *model) updateList(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) { //nolint:gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.currentlyDisplaying = displaySecondLoading

			return nil

		case "esc":
			return tea.Quit

		case "right", "tab":
			m.tabs.activeTab = min(m.tabs.activeTab+1, len(m.tabs.lists)-1)

			return nil

		case "left", "shift+tab":
			m.tabs.activeTab = max(m.tabs.activeTab-1, 0)

			return nil
		}
	}

	var cmd tea.Cmd
	m.tabs.lists[m.tabs.activeTab], cmd = m.tabs.lists[m.tabs.activeTab].Update(msg)

	return cmd
}

func (m *model) loadDetails() tea.Cmd {
	return func() tea.Msg {
		selected, ok := m.tabs.lists[m.tabs.activeTab].Items()[m.tabs.lists[m.tabs.activeTab].Index()].(itemWrapper)
		if !ok {
			return singleIsseUpdateMsg{
				errorMsg: "Failed to convert selected item to itemWrapper",
				issueID:  "",
				details:  nil,
			}
		}

		issue, ok := m.viewedIssues[selected.issue.Iid]
		if ok {
			return singleIsseUpdateMsg{
				errorMsg: "",
				issueID:  selected.issue.Iid,
				details:  &issue,
			}
		}

		// Request issue
		tmp, err := issues.QuerySingle(m.conf, m.shared.Details, m.shared.URL, selected.issue.Iid)
		if err != nil {
			return singleIsseUpdateMsg{
				errorMsg: "Failed to query issue: " + err.Error(),
				issueID:  "",
				details:  nil,
			}
		}

		return singleIsseUpdateMsg{
			errorMsg: "",
			issueID:  selected.issue.Iid,
			details:  tmp,
		}
	}
}
