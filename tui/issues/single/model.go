package single

import (
	"gn/config"
	"gn/issues"
	"gn/tui/issues/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	issue    *issues.IssueDetails
	shared   *shared.Shared
	content  string
	viewport viewport.Model
	ready    bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getIssue(&m), m.shared.Spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = shared.ViewportInitSize(&msg, m.shared.IssueID)
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			// m.viewport.YPosition = headerHeight + 1
		} else {
			shared.ViewportSetSize(&m.viewport, &msg, m.shared.IssueID)
		}

	case *issues.IssueDetails:
		// We got the issue loaded
		m.issue = msg
		m.content = shared.PrettyPrintIssue(m.issue, m.viewport.Width, m.viewport.Height)

		m.viewport.SetContent(m.content)

	default:
		m.shared.Spinner, cmd = m.shared.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready || m.issue == nil {
		return lipgloss.Place(
			m.viewport.Width,
			m.viewport.Height,
			lipgloss.Center,
			lipgloss.Center,

			"Issue is loading "+m.shared.Spinner.View(),
		)
	}

	return shared.ViewportView(&m.viewport, m.shared.IssueID)
}

func getIssue(m *model) func() tea.Msg {
	return func() tea.Msg {
		conf, err := config.Load()
		if err != nil {
			style.PrintErrAndExit("Failed to load config: " + err.Error())
		}

		issue, err := issues.QuerySingle(conf, m.shared.Details, m.shared.IssueID)
		if err != nil {
			style.PrintErrAndExit("Failed to query issue: " + err.Error())
		}

		return issue
	}
}
