package shared

import (
	"fmt"
	"strconv"
	"strings"

	"gn/issues/single"
	"gn/logger"
	"gn/style"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var issueTitleStyle = lipgloss.NewStyle().Bold(true).Underline(true)

func PrettyPrintIssue(issue *single.IssueDetails, width int, height int) string {
	_, w := style.Comment.GetFrameSize()
	availableWidth := width - w

	style.Comment.Width(width - style.Comment.GetHorizontalFrameSize())
	style.Discussion.Width(availableWidth - style.Comment.GetHorizontalFrameSize() - style.Discussion.GetHorizontalFrameSize())

	markdownOptions := []glamour.TermRendererOption{
		glamour.WithWordWrap(style.Comment.GetWidth() - 2),
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithBaseURL(issue.BaseURL.String()),
	}

	markdownOuter, err := glamour.NewTermRenderer(
		markdownOptions...,
	)
	if err != nil {
		logger.Log.Errorf("Failed to create markdown renderer: %s", err.Error())

		return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
	}
	defer markdownOuter.Close()

	// Update the word wrap length
	markdownOptions[0] = glamour.WithWordWrap(style.Discussion.GetWidth() - 2)
	markdownInner, err := glamour.NewTermRenderer(
		markdownOptions...,
	)
	if err != nil {
		logger.Log.Errorf("Failed to create markdown renderer: %s", err.Error())

		return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
	}
	defer markdownInner.Close()

	desc, err := markdownOuter.Render(issue.Description)
	if err != nil {
		logger.Log.Error("Failed to render description.", "error", err, "input", issue.Description)

		return style.FormatQuitText("Failed to render markdown: " + err.Error())
	}

	updatedAt := ""
	if issue.CreatedAt != issue.UpdatedAt {
		updatedAt = fmt.Sprintf("Updated on %s\n", issue.UpdatedAt.Format("2006-01-02 15:04"))
	}

	// Issue details
	out := style.Comment.Render(fmt.Sprintf(
		"%s\nCreated by %s on %s\n%s%s%s\n%s",
		getTitle(issue, width),
		issue.Author.String(),
		issue.CreatedAt.Format("2006-01-02 15:04"),
		updatedAt,
		getAssignees(issue),
		getLabels(issue),
		desc,
	)) + "\n"

	// Comments
	var commentBody string
	for _, comment := range issue.Discussion {
		commentBody, err = markdownOuter.Render(comment.Body)
		if err != nil {
			logger.Log.Error("Failed to render comment.", "error", err, "input", issue.Description)

			return style.FormatQuitText("Failed to render markdown: " + err.Error())
		}

		discussion := fmt.Sprintf(
			"Created by %s on %s\n\n%s\n",
			comment.Author.String(),
			comment.CreatedAt.Format("2006-01-02 15:04"),
			strings.TrimSpace(commentBody),
		)

		if len(comment.Comments) > 0 {
			discussion += "\n"
		}

		// comments on the comments
		for i, innerComment := range comment.Comments {
			commentBody, err = markdownInner.Render(innerComment.Body)
			if err != nil {
				logger.Log.Error("Failed to render inner comment.", "error", err, "input", issue.Description)

				return style.FormatQuitText("Failed to render markdown: " + err.Error())
			}

			editedBy := ""
			if innerComment.LastEditedBy.Username != "" || innerComment.LastEditedBy.Name != "" {
				editedBy = fmt.Sprintf(
					"Last edit by %s on %s\n",
					innerComment.LastEditedBy.String(),
					innerComment.UpdatedAt.Format("2006-01-02 15:04"),
				)
			}

			discussion += style.Discussion.Render(fmt.Sprintf(
				"Created by %s on %s\n%s\n%s",
				innerComment.Author.String(),
				innerComment.CreatedAt.Format("2006-01-02 15:04"),
				editedBy,
				strings.TrimSpace(commentBody),
			))

			if i < len(comment.Comments)-1 {
				discussion += "\n"
			}
		}

		out += style.Comment.Render(discussion) + "\n"
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,

		out,
	)
}

func getTitle(issue *single.IssueDetails, width int) string {
	return lipgloss.PlaceHorizontal(
		width-style.Comment.GetHorizontalFrameSize(),
		lipgloss.Center,
		issueTitleStyle.Render(issue.Title),
	)
}

func getAssignees(issue *single.IssueDetails) string {
	if len(issue.Assignees) == 0 {
		return ""
	}

	assignees := make([]string, len(issue.Assignees))
	for i, assignee := range issue.Assignees {
		assignees[i] = assignee.String()
	}

	return fmt.Sprintf("Assigned to: %s\n", strings.Join(assignees, ", "))
}

func getLabels(issue *single.IssueDetails) string {
	if len(issue.Labels) == 0 {
		return ""
	}

	labels := make([]string, len(issue.Labels))
	for i, label := range issue.Labels {
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
