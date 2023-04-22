package discussion

import (
	"errors"
	"net/url"
	"time"

	"gn/remote"
)

// ErrIssueDoesNotExist is returned if the requested discussion does not exist.
var ErrIssueDoesNotExist = errors.New("discussion with the given iid does not exist")

// Details contains the discussion of an issue.
type Details struct {
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Author      remote.User
	BaseURL     url.URL
	Assignees   []remote.User
	Labels      []Label
	Discussion  []Comment
}

// Label contains the title and color of a GitLab/GitHub label.
type Label struct {
	Title string
	Color string
}

// Comment contains the comments on an issue. Comments has a max depth of one.
type Comment struct {
	Author       remote.User
	Body         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastEditedBy remote.User
	Comments     []Comment
	Resolved     bool
}

// UpdateUsername replaces the username of the user with "you".
func (id *Details) UpdateUsername(ownUsername string) {
	if len(ownUsername) == 0 {
		return
	}

	you := remote.User{
		Name:     "you",
		Username: "",
	}

	// Update author
	if id.Author.Username == ownUsername {
		id.Author = you
	}

	// Update assignees
	for k, assignee := range id.Assignees {
		if assignee.Username == ownUsername {
			id.Assignees[k] = you
		}
	}

	// Update comments
	for k, comment := range id.Discussion {
		if comment.Author.Username == ownUsername {
			id.Discussion[k].Author = you
		}

		if comment.LastEditedBy.Username == ownUsername {
			id.Discussion[k].LastEditedBy = you
		}

		for m, inner := range comment.Comments {
			if inner.Author.Username == ownUsername {
				id.Discussion[k].Comments[m].Author = you
			}

			if inner.LastEditedBy.Username == ownUsername {
				id.Discussion[k].Comments[m].LastEditedBy = you
			}
		}
	}
}
