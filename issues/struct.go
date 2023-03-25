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
	Assignees   []User
	Labels      []string
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

func (u *User) String() string {
	if u.Username == "" {
		return u.Name
	}

	return fmt.Sprintf("%s (%s)", u.Name, u.Username)
}

func (i *Issue) HasBeenUpdated() bool {
	return i.CreatedAt != i.UpdatedAt
}
