package all

import (
	"fmt"

	"gn/config"
	"gn/issues"
	"gn/tui/issues/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	shared       *shared.Shared
	conf         *config.Wrapper
	viewedIssues map[string]issues.IssueDetails
	tabs         tabs
	error        string
	viewport     viewport.Model
	isLoading    bool
	viewingList  bool
	failure      bool
}

type tabs struct {
	lists     [3]list.Model
	activeTab int
}

type updateMsg struct {
	conf  *config.Wrapper
	items []itemWrapper
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		getIssues(&m),
		m.shared.Spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()

		for i := range m.tabs.lists {
			m.tabs.lists[i].SetSize(msg.Width-h, msg.Height-v-8)
		}

		shared.ViewportSetSize(&m.viewport, &msg, m.shared.IssueID)

		if m.viewingList {
			m.tabs.lists[m.tabs.activeTab], cmd = m.tabs.lists[m.tabs.activeTab].Update(msg)

			return m, cmd
		}

		m.viewport, cmd = m.viewport.Update(msg)

		return m, cmd

	case updateMsg:
		return updateList(&m, &msg)

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.viewingList {
			return handleListUpdate(&m, msg)
		}

		// Not needed but added for clarity
		break //nolint:gosimple

	case tea.MouseMsg:
		// Mouse msg are only intended for the viewport
		if m.viewingList {
			// We don't handle MouseMsg for the list
			return m, nil
		}

		// Not needed but added for clarity
		break //nolint:gosimple

	default:
		m.shared.Spinner, cmd = m.shared.Spinner.Update(msg)

		return m, cmd
	}

	// Can only be reached if we are displaying the viewport, and we either got a tea.MouseMsg or a tea.KeyMsg
	return handleViewportUpdate(&m, msg)
}

func (m model) View() string {
	if m.failure {
		return style.FormatQuitText("An error occurred: " + m.error)
	}

	if m.isLoading {
		return lipgloss.Place(
			m.tabs.lists[m.tabs.activeTab].Width(),
			m.tabs.lists[m.tabs.activeTab].Height(),
			lipgloss.Center,
			lipgloss.Center,

			fmt.Sprintf("Loading %s", m.shared.Spinner.View()),
		)
	}

	if m.viewingList {
		return renderTab(&m.tabs)
	}

	return shared.ViewportView(&m.viewport, m.shared.IssueID)
}

func getIssues(m *model) func() tea.Msg {
	return func() tea.Msg {
		conf, err := config.Load()
		if err != nil {
			style.PrintErrAndExit("Failed to load config: " + err.Error())
		}

		allIssues, err := issues.QueryAll(conf, m.shared.Details, m.shared.URL)
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
			items: issueList,
			conf:  conf,
		}
	}
}
