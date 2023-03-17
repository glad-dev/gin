package issues

import "time"

type Issue struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Iid         string    `json:"iid"`
	State       string    `json:"state"`
	Assignees   []string  `json:assignees`
}
