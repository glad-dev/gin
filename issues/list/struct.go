package list

import (
	"time"

	"gn/remote"
)

type Issue struct {
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Iid       string
	State     string
	Author    remote.User
	Assignees []remote.User
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
			i.Assignees[k] = remote.User{
				Name:     "you",
				Username: "",
			}
		}
	}
}
