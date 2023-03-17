package issues

import "time"

type Issue struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Iid         string    `json:"iid"`
	State       string    `json:"state"`
	Author      User      `json:"author"`
	Assignees   []User    `json:"assignees"`
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
	Resolved     bool        `json:"resolved"`
	LastEditedBy interface{} `json:"lastEditedBy"` // ToDo: Remove interface
	Comments     []Comment   `json:"comments"`
}

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}
