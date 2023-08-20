package list

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/remote"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

type issuesStruct struct {
	PageInfo struct {
		EndCursor   graphql.String
		HasNextPage graphql.Boolean
	}

	Nodes []struct {
		Title     graphql.String
		State     graphql.String
		CreatedAt graphql.String
		UpdatedAt graphql.String
		Author    struct {
			Login graphql.String
		}
		Assignees struct {
			Nodes []struct {
				Login graphql.String
			}
		} `graphql:"assignees(first: 100)"`
		Number graphql.Int
	}
}

type firstQuery struct {
	Repository struct {
		Issues issuesStruct `graphql:"issues(first: 100, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type followingQuery struct { // Needed since GitHub considers an empty "after" to be invalid
	Repository struct {
		Issues issuesStruct `graphql:"issues(first: 100, after: $after, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// QueryGitHub returns all issues, open and closed, of a given repository.
func QueryGitHub(match *remote.Match, projectPath string) ([]Issue, error) {
	tmp := strings.Split(projectPath, "/")
	if len(tmp) != 2 {
		logger.Log.Error("Project path is invalid", "path", projectPath)

		return nil, errors.New("invalid project path")
	}

	var tc *http.Client
	if len(match.Token) > 0 {
		ctx := context.Background()
		tc = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: match.Token},
		))
	}

	apiURL, err := match.Type.ApiURL(&match.URL)
	if err != nil {
		logger.Log.Error("Failed to get API URL", "error", err, "match-url", match.URL.String())

		return nil, fmt.Errorf("invalid API url: %w", err)
	}

	client := graphql.NewClient(apiURL, tc)

	fq := &firstQuery{}
	err = client.Query(context.Background(), fq, map[string]any{
		"owner": graphql.String(tmp[0]), // owner is at tmp[0], repo name is at tmp[1]
		"name":  graphql.String(tmp[1]),
	})
	if err != nil {
		logger.Log.Error("First GitHub query failed", "error", err)

		return nil, fmt.Errorf("query failed: %w", err)
	}

	lst := flatten(fq.Repository.Issues)
	if !fq.Repository.Issues.PageInfo.HasNextPage {
		return lst, nil
	}

	issueList := make([]Issue, 0)
	issueList = append(issueList, lst...)

	cursor := fq.Repository.Issues.PageInfo.EndCursor
	q := &followingQuery{}
	for {
		err = client.Query(context.Background(), q, map[string]any{
			"owner": graphql.String(tmp[0]), // owner is at tmp[0], repo name is at tmp[1]
			"name":  graphql.String(tmp[1]),
			"after": cursor,
		})
		if err != nil {
			logger.Log.Error("GitHub query failed", "error", err)

			return nil, fmt.Errorf("query failed: %w", err)
		}

		issueList = append(issueList, flatten(q.Repository.Issues)...)

		cursor = q.Repository.Issues.PageInfo.EndCursor
		if !q.Repository.Issues.PageInfo.HasNextPage {
			break
		}
	}

	return issueList, nil
}

func flatten(issues issuesStruct) []Issue {
	lst := make([]Issue, 0)

	for _, issue := range issues.Nodes {
		assignees := make([]remote.User, len(issue.Assignees.Nodes))
		for i, u := range issue.Assignees.Nodes {
			assignees[i] = remote.User{
				Username: string(u.Login),
			}
		}

		creationTime, err := time.Parse(timeLayout, string(issue.CreatedAt))
		if err != nil {
			logger.Log.Warn("failed to parse creation time", "time", string(issue.CreatedAt), "error", err)
		}

		updateTime, err := time.Parse(timeLayout, string(issue.UpdatedAt))
		if err != nil {
			logger.Log.Warn("failed to parse update time", "time", string(issue.UpdatedAt), "error", err)
		}

		lst = append(lst, Issue{
			Title:     string(issue.Title),
			CreatedAt: creationTime,
			UpdatedAt: updateTime,
			Iid:       strconv.Itoa(int(issue.Number)),
			State:     string(issue.State),
			Author: remote.User{
				Username: string(issue.Author.Login),
			},
			Assignees: assignees,
		})
	}

	return lst
}
