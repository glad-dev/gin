package shared

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

func ViewportView(vp *viewport.Model, issueID string) string {
	return fmt.Sprintf(
		"%s\n%s\n%s",
		viewportHeader(issueID, vp.Width),
		vp.View(),
		viewportFooter(vp),
	)
}

func ViewportSetSize(vp *viewport.Model, msg *tea.WindowSizeMsg, issueID string) {
	headerHeight := lipgloss.Height(viewportHeader(issueID, vp.Width))
	footerHeight := lipgloss.Height(viewportFooter(vp))
	verticalMarginHeight := headerHeight + footerHeight

	vp.Width = msg.Width
	vp.Height = msg.Height - verticalMarginHeight
}

func ViewportInitSize(msg *tea.WindowSizeMsg, issueID string) viewport.Model {
	vp := viewport.New(0, 0)
	headerHeight := lipgloss.Height(viewportHeader(issueID, vp.Width))

	ViewportSetSize(&vp, msg, issueID)
	vp.YPosition = headerHeight

	return vp
}

func viewportHeader(issueID string, width int) string {
	title := titleStyle.Render("Details of issue #" + issueID)
	line := strings.Repeat("─", max(0, width-lipgloss.Width(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func viewportFooter(vp *viewport.Model) string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", vp.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, vp.Width-lipgloss.Width(info)))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
