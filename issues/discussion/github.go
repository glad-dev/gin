package discussion

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/remote/match"

	"github.com/google/go-github/v59/github"
	"golang.org/x/oauth2"
)

// QueryGitHub returns the discussion associated with the passed issueID. If the requested issue does not exist, an
// ErrIssueDoesNotExist is returned.
func QueryGitHub(match *match.Match, projectPath string, issueID string) (*Details, error) {
	tmp := strings.Split(projectPath, "/")
	if len(tmp) != 2 {
		log.Error("Project path is invalid", "path", projectPath)

		return nil, errors.New("invalid project path")
	}

	issueNumber, err := strconv.Atoi(issueID)
	if err != nil {
		log.Error("Failed to convert issueID to int.", "error", err, "issueID", issueID)

		return nil, fmt.Errorf("failed to convert issueID to int: %w", err)
	}

	var tc *http.Client
	if len(match.Token) > 0 {
		ctx := context.Background()
		tc = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: match.Token},
		))
	}
	client := github.NewClient(tc)

	// Get issue details
	issue, _, err := client.Issues.Get(context.Background(), tmp[0], tmp[1], issueNumber)
	if err != nil {
		return nil, fmt.Errorf("requesting issue: %w", err)
	}

	assignee := issue.GetAssignee()
	author := issue.GetUser()

	details := &Details{
		Title:       issue.GetTitle(),
		Description: issue.GetBody(),
		CreatedAt:   issue.GetCreatedAt().Time,
		UpdatedAt:   issue.GetUpdatedAt().Time,
		Author: remote.User{
			Name:     author.GetName(),
			Username: author.GetLogin(),
		},
		BaseURL: url.URL{
			Scheme: "https",
			Host:   "github.com",
		},
		Assignees: []remote.User{
			{
				Name:     assignee.GetName(),
				Username: assignee.GetLogin(),
			},
		},
		Labels:     make([]Label, 0),
		Discussion: make([]Comment, 0),
	}

	for _, label := range issue.Labels {
		details.Labels = append(details.Labels, Label{
			Title: label.GetName(),
			Color: label.GetColor(),
		})
	}

	// Get comments
	page := 0
	sort := "created"
	direction := "desc"
	for {
		comments, response, err := client.Issues.ListComments(context.Background(), tmp[0], tmp[1], issueNumber, &github.IssueListCommentsOptions{
			Sort:      &sort,
			Direction: &direction,
			Since:     nil,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("requesting comments: %w", err)
		}

		for _, comment := range comments {
			author = comment.GetUser()

			details.Discussion = append(details.Discussion, Comment{
				Author: remote.User{
					Name:     author.GetName(),
					Username: author.GetLogin(),
				},
				Body:      comment.GetBody(),
				CreatedAt: comment.GetCreatedAt().Time,
				UpdatedAt: comment.GetUpdatedAt().Time,
			})
		}

		page = response.NextPage
		if page == 0 {
			break
		}
	}

	return details, nil
}
