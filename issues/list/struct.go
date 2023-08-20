package list

import (
	"time"

	"github.com/glad-dev/gin/remote"
)

const timeLayout = "2006-01-02T15:04:05Z"

// Issue contains information about an issue.
type Issue struct {
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Iid       string
	State     string
	Author    remote.User
	Assignees []remote.User
}

// UpdateUsername replaces the mentions of the user's username with "you".
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
