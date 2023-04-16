package single

import (
	"net/url"
	"time"

	"gn/issues/user"
)

type IssueDetails struct {
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Author      user.Details
	BaseURL     url.URL
	Assignees   []user.Details
	Labels      []Label
	Discussion  []Comment
}

type Label struct {
	Title string
	Color string
}

type Comment struct {
	Author       user.Details
	Body         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastEditedBy user.Details
	Comments     []Comment
	Resolved     bool
}

func (id *IssueDetails) UpdateUsername(ownUsername string) {
	if len(ownUsername) == 0 {
		return
	}

	you := user.Details{
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
