package single

import (
	"fmt"
	"gn/config"
	"gn/issues"
	"gn/repo"
	"gn/tui/style"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"

		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"

		return titleStyle.Copy().BorderStyle(b)
	}()
)

type model struct {
	issue    *issues.IssueDetails
	content  string
	viewport viewport.Model
	shared   Shared
	ready    bool
}

type Shared struct {
	issueID string
	details []repo.Details
}

func (m model) Init() tea.Cmd {
	return getIssue(&m)
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
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}

	case *issues.IssueDetails:
		// We got the issue loaded
		m.issue = msg
		m.content = prettyPrintIssue(&m)

		m.viewport.SetContent(m.content)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.issue == nil {
		return "Issue is loading"
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) headerView() string {
	title := titleStyle.Render("Details of issue #" + m.shared.issueID)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func getIssue(m *model) func() tea.Msg {
	return func() tea.Msg {
		conf, err := config.Load()
		if err != nil {
			style.PrintErrAndExit("Failed to load config: " + err.Error())
		}

		issue, err := issues.QuerySingle(conf, m.shared.details, m.shared.issueID)
		if err != nil {
			style.PrintErrAndExit("Failed to query issue: " + err.Error())
		}

		return issue
	}
}

var issueTitleStyle = lipgloss.NewStyle().Bold(true)

func prettyPrintIssue(m *model) string {
	_, w := style.InputField.GetFrameSize()

	style.InputField.Width(m.viewport.Width - w)

	// Assignees
	assignees := make([]string, len(m.issue.Assignees))
	for i, assignee := range m.issue.Assignees {
		assignees[i] = assignee.String()
	}

	assigneeStr := ""
	if len(assignees) > 0 {
		assigneeStr = fmt.Sprintf("Assigned to: %s\n", strings.Join(assignees, ", "))
	}

	labels := make([]string, len(m.issue.Labels))
	for i, label := range m.issue.Labels {
		labels[i] = lipgloss.NewStyle().
			Background(lipgloss.Color(label.Color)).
			Foreground(getInverseColor(label.Color)).
			Render(label.Title)
	}

	labelStr := ""
	if len(labels) > 0 {
		labelStr = fmt.Sprintf("Labels: %s\n", strings.Join(labels, ", "))
	}

	// Header card
	out := style.InputField.Render(fmt.Sprintf(
		"%s\nCreated by %s\n%s%s\n%s",
		lipgloss.PlaceHorizontal(m.viewport.Width-style.InputField.GetHorizontalFrameSize(), lipgloss.Center, issueTitleStyle.Render(m.issue.Title)),
		m.issue.Author.String(),
		assigneeStr,
		labelStr,
		m.issue.Description,
	)) + "\n"

	// Comments
	for _, comment := range m.issue.Discussion {
		out += style.InputField.Render(comment.Body) + "\n"
	}

	style.InputField.Width(80)

	return lipgloss.Place(
		m.viewport.Width,
		m.viewport.Height,
		lipgloss.Center,
		lipgloss.Center,

		out,
	)
}

func getInverseColor(hexStr string) lipgloss.Color {
	hexInt64, err := strconv.ParseInt(strings.Replace(hexStr, "#", "", 1), 16, 0)
	if err != nil {
		return ""
	}

	return lipgloss.Color(fmt.Sprintf("#%x", int(hexInt64)^0xffffff))
}
