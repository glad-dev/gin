package issueList

import (
	"time"

	"gn/issues/user"
)

type Issue struct {
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Iid       string
	State     string
	Author    user.Details
	Assignees []user.Details
}

func (i *Issue) HasBeenUpdated() bool {
	return i.CreatedAt != i.UpdatedAt
}

func (i *Issue) UpdateUsername(ownUsername string) {
	if len(ownUsername) == 0 {
		return
	}

	// Update author
	if i.Author.Username == ownUsername {
		i.Author.Name = "you"
		i.Author.Username = ""
	}

	// Update assignees
	for k, assignee := range i.Assignees {
		if assignee.Username == ownUsername {
			i.Assignees[k] = user.Details{
				Name:     "you",
				Username: "",
			}
		}
	}
}
