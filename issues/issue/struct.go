package issue

import (
	"net/url"
	"time"

	"gn/remote"
)

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

type Label struct {
	Title string
	Color string
}

type Comment struct {
	Author       remote.User
	Body         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastEditedBy remote.User
	Comments     []Comment
	Resolved     bool
}

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
