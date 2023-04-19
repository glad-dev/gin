package all

import (
	"fmt"

	"gn/config"
	"gn/issues"
	"gn/issues/issue"
	"gn/style"
	"gn/tui/issues/shared"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	displaying uint8
	state      uint8
)

const (
	displayingInitialLoading displaying = iota
	displayingList
	displaySecondLoading
	displayingDetails
)

const (
	stateRunning state = iota
	exitFailure
)

type model struct {
	shared              *shared.Shared
	conf                *config.Wrapper
	viewedIssues        map[string]issue.Details
	tabs                tabs
	error               string
	viewport            viewport.Model
	state               state
	currentlyDisplaying displaying
}

type tabs struct {
	lists     [3]list.Model
	activeTab int
}

type allIssuesUpdateMsg struct {
	errorMsg string
	items    []itemWrapper
}

type singleIssueUpdateMsg struct {
	details  *issue.Details
	errorMsg string
	issueID  string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		getIssues(&m),
		m.shared.Spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()

		for i := range m.tabs.lists {
			m.tabs.lists[i].SetSize(msg.Width-h, msg.Height-v-8)
		}

		shared.ViewportSetSize(&m.viewport, &msg, m.shared.IssueID)

		return m, nil

	case allIssuesUpdateMsg:
		if len(msg.errorMsg) > 0 {
			m.error = msg.errorMsg
			m.state = exitFailure

			return m, tea.Quit
		}

		return m.initList(&msg)

	case singleIssueUpdateMsg:
		if len(msg.errorMsg) > 0 {
			m.error = msg.errorMsg
			m.state = exitFailure

			return m, tea.Quit
		}

		// Store issue in map
		m.viewedIssues[m.shared.IssueID] = *msg.details

		m.shared.IssueID = msg.issueID
		m.viewport.SetContent(shared.PrettyPrintIssue(msg.details, m.viewport.Width, m.viewport.Height))
		m.currentlyDisplaying = displayingDetails
		m.viewport.GotoTop()

		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	cmds := make([]tea.Cmd, 2)
	m.shared.Spinner, cmds[0] = m.shared.Spinner.Update(msg)

	switch m.currentlyDisplaying {
	case displayingInitialLoading:
		return m, cmds[0]

	case displayingList:
		cmds[1] = m.updateList(msg)

		return m, tea.Batch(cmds...)

	case displaySecondLoading:
		return m, tea.Batch(
			cmds[0],
			m.loadDetails(),
		)

	case displayingDetails:
		cmds[1] = m.updateViewport(msg)

		return m, tea.Batch(cmds...)

	default:
		m.state = exitFailure
		m.error = "Invalid update state"

		return m, tea.Quit
	}
}

func (m model) View() string {
	if m.state == exitFailure {
		// This isn't shown?
		return style.FormatQuitText("An error occurred: " + m.error)
	}

	switch m.currentlyDisplaying {
	case displayingInitialLoading, displaySecondLoading:
		return lipgloss.Place(
			m.tabs.lists[m.tabs.activeTab].Width(),
			m.tabs.lists[m.tabs.activeTab].Height(),
			lipgloss.Center,
			lipgloss.Center,

			fmt.Sprintf("Loading %s", m.shared.Spinner.View()),
		)

	case displayingList:
		return renderTab(&m.tabs)

	case displayingDetails:
		return shared.ViewportView(&m.viewport, m.shared.IssueID)

	default:
		return "Unknown view"
	}
}

func getIssues(m *model) func() tea.Msg {
	return func() tea.Msg {
		allIssues, err := issues.QueryList(m.conf, m.shared.Details, m.shared.URL)
		if err != nil {
			return allIssuesUpdateMsg{
				items:    nil,
				errorMsg: "Failed to query issues: " + err.Error(),
			}
		}

		issueList := make([]itemWrapper, len(allIssues))
		for i, issue := range allIssues {
			issueList[i] = itemWrapper{
				issue: issue,
			}
		}

		return allIssuesUpdateMsg{
			items: issueList,
		}
	}
}
