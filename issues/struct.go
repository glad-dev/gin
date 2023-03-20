package issues

import (
	"fmt"
	"time"
)

type Issue struct {
	title       string
	description string
	createdAt   time.Time
	updatedAt   time.Time
	iid         string
	state       string
	author      User
	assignees   []User
}

type IssueDetails struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Assignees   []User    `json:"assignees"`
	Labels      []string  `json:"labels"`
	Discussion  []Comment `json:"discussion"`
}

type Comment struct {
	Author       User        `json:"author"`
	Body         string      `json:"body"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	LastEditedBy interface{} `json:"lastEditedBy"` // ToDo: Remove interface
	Comments     []Comment   `json:"comments"`
	Resolved     bool        `json:"resolved"`
}

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (user *User) String() string {
	return fmt.Sprintf("%s (%s)", user.Name, user.Username)
}

// Title should only be used by bubbletea. To ge the title, use TitleOnly.
func (i Issue) Title() string {
	status := ""
	if i.state == "closed" {
		status = "[closed] "
	}

	return fmt.Sprintf("#%s %s%s by %s on %s", i.iid, status, i.title, i.author.String(), i.createdAt.Format("2006-01-02 15:04"))
}

func (i Issue) Description() string {
	return i.description
}

// FilterValue is needed to allow Issue to be a bubbletea list.Item.
func (i Issue) FilterValue() string {
	return i.Title() + i.description
}

func (i Issue) TitleOnly() string {
	return i.title
}

func (i Issue) CreatedAt() time.Time {
	return i.createdAt
}

func (i Issue) UpdatedAt() time.Time {
	return i.updatedAt
}

func (i Issue) Iid() string {
	return i.iid
}

func (i Issue) State() string {
	return i.state
}

func (i Issue) Author() User {
	return i.author
}

func (i Issue) Assignees() []User {
	return i.assignees
}

func (i Issue) HasBeenUpdated() bool {
	return i.createdAt != i.UpdatedAt()
}
