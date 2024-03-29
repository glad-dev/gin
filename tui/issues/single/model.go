package single

import (
	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/issues"
	"github.com/glad-dev/gin/issues/discussion"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/issues/shared"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	shared   *shared.Shared
	conf     *configuration.Config
	content  string
	viewport viewport.Model
	ready    bool
	failure  bool
}

type updateMsg struct {
	discussion *discussion.Details
	errText    string
}

// Init is required for model to be a tea.Model.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		getIssue(&m),
		m.shared.Spinner.Tick,
	)
}

// Update is required for model to be a tea.Model.
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

		// We got the discussion loaded
		m.content = shared.PrettyPrintIssue(msg.discussion, m.viewport.Width, m.viewport.Height)
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

// View required for model to be a tea.Model.
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
				discussion: nil,
				errText:    style.FormatQuitText("Failed to query discussion: " + err.Error()),
			}
		}

		return updateMsg{
			discussion: issue,
			errText:    "",
		}
	}
}
