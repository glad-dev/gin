package shared

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/glad-dev/gin/issues/discussion"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var issueTitleStyle = lipgloss.NewStyle().Bold(true).Underline(true)

// PrettyPrintIssue renders the passed discussion as markdown.
func PrettyPrintIssue(details *discussion.Details, width int, height int) string {
	_, w := style.Comment.GetFrameSize()
	availableWidth := width - w

	style.Comment.Width(width - style.Comment.GetHorizontalFrameSize())
	style.Discussion.Width(availableWidth - style.Comment.GetHorizontalFrameSize() - style.Discussion.GetHorizontalFrameSize())

	markdownOptions := []glamour.TermRendererOption{
		glamour.WithWordWrap(style.Comment.GetWidth() - 2),
		glamour.WithAutoStyle(),
		glamour.WithEmoji(),
		glamour.WithBaseURL(details.BaseURL.String()),
	}

	markdownOuter, err := glamour.NewTermRenderer(
		markdownOptions...,
	)
	if err != nil {
		log.Error("Failed to create markdown renderer", "error", err.Error())

		return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
	}
	defer markdownOuter.Close()

	// Update the word wrap length
	markdownOptions[0] = glamour.WithWordWrap(style.Discussion.GetWidth() - 2)
	markdownInner, err := glamour.NewTermRenderer(
		markdownOptions...,
	)
	if err != nil {
		log.Error("Failed to create markdown renderer", "error", err.Error())

		return style.FormatQuitText("Failed to create markdown renderer: " + err.Error())
	}
	defer markdownInner.Close()

	desc, err := markdownOuter.Render(details.Description)
	if err != nil {
		log.Error("Failed to render description.", "error", err, "input", details.Description)

		return style.FormatQuitText("Failed to render markdown: " + err.Error())
	}

	updatedAt := ""
	if details.CreatedAt != details.UpdatedAt {
		updatedAt = fmt.Sprintf("Updated on %s\n", details.UpdatedAt.Format("2006-01-02 15:04"))
	}

	// Issue details
	out := style.Comment.Render(fmt.Sprintf(
		"%s\nCreated by %s on %s\n%s%s%s\n%s",
		getTitle(details, width),
		details.Author.String(),
		details.CreatedAt.Format("2006-01-02 15:04"),
		updatedAt,
		getAssignees(details),
		getLabels(details),
		desc,
	)) + "\n"

	// Comments
	var commentBody string
	for _, comment := range details.Discussion {
		commentBody, err = markdownOuter.Render(comment.Body)
		if err != nil {
			log.Error("Failed to render comment.", "error", err, "input", details.Description)

			return style.FormatQuitText("Failed to render markdown: " + err.Error())
		}

		currentDiscussion := fmt.Sprintf(
			"Created by %s on %s\n\n%s\n",
			comment.Author.String(),
			comment.CreatedAt.Format("2006-01-02 15:04"),
			strings.TrimSpace(commentBody),
		)

		if len(comment.Comments) > 0 {
			currentDiscussion += "\n"
		}

		// comments on the comments
		for i, innerComment := range comment.Comments {
			commentBody, err = markdownInner.Render(innerComment.Body)
			if err != nil {
				log.Error("Failed to render inner comment.", "error", err, "input", details.Description)

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

			currentDiscussion += style.Discussion.Render(fmt.Sprintf(
				"Created by %s on %s\n%s\n%s",
				innerComment.Author.String(),
				innerComment.CreatedAt.Format("2006-01-02 15:04"),
				editedBy,
				strings.TrimSpace(commentBody),
			))

			if i < len(comment.Comments)-1 {
				currentDiscussion += "\n"
			}
		}

		out += style.Comment.Render(currentDiscussion) + "\n"
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,

		out,
	)
}

func getTitle(details *discussion.Details, width int) string {
	return lipgloss.PlaceHorizontal(
		width-style.Comment.GetHorizontalFrameSize(),
		lipgloss.Center,
		issueTitleStyle.Render(details.Title),
	)
}

func getAssignees(details *discussion.Details) string {
	if len(details.Assignees) == 0 {
		return ""
	}

	assignees := make([]string, len(details.Assignees))
	for i, assignee := range details.Assignees {
		assignees[i] = assignee.String()
	}

	return fmt.Sprintf("Assigned to: %s\n", strings.Join(assignees, ", "))
}

func getLabels(details *discussion.Details) string {
	if len(details.Labels) == 0 {
		return ""
	}

	labels := make([]string, len(details.Labels))
	for i, label := range details.Labels {
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
