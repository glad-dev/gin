package all

import (
	"log"

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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg.String() {
	case "enter":
		if m.isLoading {
			return m, nil
		}

		m.isLoading = true

		selected, ok := m.tabs.lists[m.tabs.activeTab].Items()[m.tabs.lists[m.tabs.activeTab].Index()].(itemWrapper)
		if !ok {
			return m, tea.Quit
		}

		m.shared.IssueID = selected.issue.Iid

		is, ok := m.viewedIssues[m.shared.IssueID]
		if !ok {
			// Request issue
			tmp, err := issues.QuerySingle(m.conf, m.shared.Details, selected.issue.Iid)
			if err != nil {
				// TODO: Remove log.Fatal
				log.Fatalln(err)
			}

			// Store issue
			m.viewedIssues[m.shared.IssueID] = *tmp

			// Copy tmp to is
			is = *tmp
		}

		m.viewport.SetContent(shared.PrettyPrintIssue(&is, m.viewport.Width, m.viewport.Height))
		m.viewingList = false
		m.isLoading = false

	case "esc", "backspace":
		return m, tea.Quit

	case "right", "tab":
		m.tabs.activeTab = min(m.tabs.activeTab+1, len(m.tabs.lists)-1)

		return m, nil

	case "left", "shift+tab":
		m.tabs.activeTab = max(m.tabs.activeTab-1, 0)

		return m, nil
	}

	m.tabs.lists[m.tabs.activeTab], cmd = m.tabs.lists[m.tabs.activeTab].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
