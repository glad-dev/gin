package issues

import (
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

func (i Issue) Title() string {
	return i.title
}

func (i Issue) Description() string {
	return i.description
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

// FilterValue is needed to allow Issue to be a bubbletea list.Item.
func (i Issue) FilterValue() string {
	return i.title + i.description
}
