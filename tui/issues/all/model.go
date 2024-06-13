package all

import (
	"fmt"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/issues"
	"github.com/glad-dev/gin/issues/discussion"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/issues/shared"

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
	conf                *configuration.Config
	viewedIssues        map[string]discussion.Details
	channel             chan int
	tabs                tabs
	error               string
	viewport            viewport.Model
	alreadyLoaded       int
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
	details  *discussion.Details
	errorMsg string
	issueID  string
}

type alreadyLoadedItemsMsg struct {
	newValue int
}

// Init is required for model to be a tea.Model.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		getIssues(&m),
		updateLoadedCount(&m),
		m.shared.Spinner.Tick,
	)
}

// Update is required for model to be a tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()

		for i := range m.tabs.lists {
			m.tabs.lists[i].SetSize(msg.Width-h, msg.Height-v-8)
		}

		shared.ViewportSetSize(&m.viewport, &msg, m.shared.IssueID)

		return m, nil

	case alreadyLoadedItemsMsg:
		if msg.newValue > m.alreadyLoaded {
			m.alreadyLoaded = msg.newValue
		}

		return m, nil

	case allIssuesUpdateMsg:
		if len(msg.errorMsg) > 0 {
			m.error = msg.errorMsg
			m.state = exitFailure

			return m, tea.Quit
		}

		return initList(&m, &msg)

	case singleIssueUpdateMsg:
		if len(msg.errorMsg) > 0 {
			m.error = msg.errorMsg
			m.state = exitFailure

			return m, tea.Quit
		}

		// Store discussion in map
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
		cmds[1] = updateLoadedCount(&m)

		return m, tea.Batch(cmds...)

	case displayingList:
		cmds[1] = updateList(&m, msg)

		return m, tea.Batch(cmds...)

	case displaySecondLoading:
		return m, tea.Batch(
			cmds[0],
			loadDetails(&m),
		)

	case displayingDetails:
		cmds[1] = updateViewport(&m, msg)

		return m, tea.Batch(cmds...)

	default:
		m.state = exitFailure
		m.error = "Invalid update state"

		return m, tea.Quit
	}
}

// View required for model to be a tea.Model.
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

			fmt.Sprintf("%s Loading issues (%d so far)", m.shared.Spinner.View(), m.alreadyLoaded),
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
		allIssues, err := issues.QueryList(m.conf, m.shared.Details, m.shared.URL, m.channel)
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

func updateLoadedCount(m *model) func() tea.Msg {
	return func() tea.Msg {
		select {
		case loaded := <-m.channel:
			return alreadyLoadedItemsMsg{newValue: loaded}
		default:
			// Channel is empty
		}

		return nil
	}
}
