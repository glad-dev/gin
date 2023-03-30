package single

import (
	"fmt"
	"strconv"
	"strings"

	"gn/tui/style"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var issueTitleStyle = lipgloss.NewStyle().Bold(true).Underline(true)

func prettyPrintIssue(m *model) string {
	_, w := style.Comment.GetFrameSize()
	availableWidth := m.viewport.Width - w

	style.Comment.Width(m.viewport.Width - style.Comment.GetHorizontalFrameSize())
	style.Discussion.Width(availableWidth - style.Comment.GetHorizontalFrameSize() - style.Discussion.GetHorizontalFrameSize())

	outerSpace := style.Comment.GetWidth() - style.Comment.GetHorizontalFrameSize()
	innerSpace := style.Discussion.GetWidth() - style.Discussion.GetHorizontalFrameSize()

	markdownOuter, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(outerSpace),
	)
	if err != nil {
		return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
	}
	defer markdownOuter.Close()

	markdownInner, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(innerSpace),
	)
	if err != nil {
		return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
	}
	defer markdownInner.Close()

	desc, err := markdownOuter.Render(m.issue.Description)
	if err != nil {
		return style.FormatQuitText("Failed to render markdown: " + err.Error())
	}

	// Issue details
	out := style.Comment.Render(fmt.Sprintf(
		"%s\nCreated by %s\nCreated on %s\n%s%s\n%s",
		getTitle(m),
		m.issue.Author.String(),
		m.issue.CreatedAt.Format("2006-01-02 15:04"),
		getAssignees(m),
		getLabels(m),
		desc,
	)) + "\n"

	// Comments
	var commentBody string
	for _, comment := range m.issue.Discussion {
		commentBody, err = markdownOuter.Render(comment.Body)
		if err != nil {
			return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
		}

		discussion := fmt.Sprintf(
			"Created by %s\nCreated on %s\n%s\n",
			comment.Author.String(),
			comment.CreatedAt.Format("2006-01-02 15:04"),
			commentBody,
		)

		// comments on the comments
		for i, innerComment := range comment.Comments {
			commentBody, err = markdownInner.Render(innerComment.Body)
			if err != nil {
				return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
			}

			discussion += style.Discussion.Render(fmt.Sprintf(
				"Created by %s\nCreated on %s\n%s\n",
				innerComment.Author.String(),
				innerComment.CreatedAt.Format("2006-01-02 15:04"),
				commentBody,
			))

			if i < len(comment.Comments)-1 {
				discussion += "\n"
			}
		}

		out += style.Comment.Render(discussion) + "\n"
	}

	return lipgloss.Place(
		m.viewport.Width,
		m.viewport.Height,
		lipgloss.Center,
		lipgloss.Center,

		out,
	)
}

func getTitle(m *model) string {
	return lipgloss.PlaceHorizontal(
		m.viewport.Width-style.Comment.GetHorizontalFrameSize(),
		lipgloss.Center,
		issueTitleStyle.Render(m.issue.Title),
	)
}

func getAssignees(m *model) string {
	if len(m.issue.Assignees) == 0 {
		return ""
	}

	assignees := make([]string, len(m.issue.Assignees))
	for i, assignee := range m.issue.Assignees {
		assignees[i] = assignee.String()
	}

	return fmt.Sprintf("Assigned to: %s\n", strings.Join(assignees, ", "))
}

func getLabels(m *model) string {
	if len(m.issue.Labels) == 0 {
		return ""
	}

	labels := make([]string, len(m.issue.Labels))
	for i, label := range m.issue.Labels {
		labels[i] = lipgloss.NewStyle().
			Background(lipgloss.Color(label.Color)).
			Foreground(getInverseColor(label.Color)).
			Render(label.Title)
	}

	return fmt.Sprintf("Labels: %s\n", strings.Join(labels, ", "))
}

func getInverseColor(hexStr string) lipgloss.Color {
	hexInt64, err := strconv.ParseInt(strings.Replace(hexStr, "#", "", 1), 16, 0)
	if err != nil {
		return ""
	}

	return lipgloss.Color(fmt.Sprintf("#%x", int(hexInt64)^0xffffff))
}
