package single

import (
	"gn/config"
	"gn/issues"
	"gn/issues/single"
	"gn/style"
	"gn/tui/issues/shared"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	shared   *shared.Shared
	conf     *config.Wrapper
	content  string
	viewport viewport.Model
	ready    bool
	failure  bool
}

type updateMsg struct {
	issue   *single.IssueDetails
	errText string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		getIssue(&m),
		m.shared.Spinner.Tick,
	)
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
		} else {
			shared.ViewportSetSize(&m.viewport, &msg, m.shared.IssueID)
		}

	case updateMsg:
		// Check if there was an error
		if len(msg.errText) != 0 {
			m.content = msg.errText
			m.failure = true

			return m, tea.Quit
		}

		// We got the issue loaded
		m.content = shared.PrettyPrintIssue(msg.issue, m.viewport.Width, m.viewport.Height)
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
	if !m.ready || len(m.content) == 0 {
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
		issue, err := issues.QuerySingle(m.conf, m.shared.Details, m.shared.URL, m.shared.IssueID)
		if err != nil {
			return updateMsg{
				issue:   nil,
				errText: style.FormatQuitText("Failed to query issue: " + err.Error()),
			}
		}

		return updateMsg{
			issue:   issue,
			errText: "",
		}
	}
}
