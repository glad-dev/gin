package issues

import (
	"fmt"
	"time"
)

type Issue struct {
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Iid       string
	State     string
	Author    User
	Assignees []User
}

type IssueDetails struct {
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Author      User
	Assignees   []User
	Labels      []Label
	Discussion  []Comment
}

type Comment struct {
	Author       User
	Body         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastEditedBy interface{} // ToDo: Remove interface
	Comments     []Comment
	Resolved     bool
}

type User struct {
	Name     string
	Username string
}

type Label struct {
	Title string
	Color string
}

func (u *User) String() string {
	if u.Username == "" {
		return u.Name
	}

	return fmt.Sprintf("%s (%s)", u.Name, u.Username)
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
			i.Assignees[k] = User{
				Name:     "you",
				Username: "",
			}
		}
	}
}

func (id *IssueDetails) UpdateUsername(ownUsername string) {
	if len(ownUsername) == 0 {
		return
	}

	// Update author
	if id.Author.Username == ownUsername {
		id.Author = User{
			Name:     "you",
			Username: "",
		}
	}

	// Update assignees
	for k, assignee := range id.Assignees {
		if assignee.Username == ownUsername {
			id.Assignees[k] = User{
				Name:     "you",
				Username: "",
			}
		}
	}

	// Update comments
	for k, comment := range id.Discussion {
		if comment.Author.Username == ownUsername {
			id.Discussion[k].Author = User{
				Name:     "you",
				Username: "",
			}
		}
	}
}
