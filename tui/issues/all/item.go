package all

import (
	"fmt"

	"github.com/glad-dev/gin/issues/list"
)

// itemWrapper is a list.Item and contains a list.Issue.
type itemWrapper struct {
	issue list.Issue
}

// Title is required for itemWrapper to be a list.Item.
func (i itemWrapper) Title() string {
	return fmt.Sprintf(
		"#%s %s",
		i.issue.Iid,
		i.issue.Title,
	)
}

// Description is required for itemWrapper to be a list.Item.
func (i itemWrapper) Description() string {
	// Use author and creation date as description
	return fmt.Sprintf(
		"Created by %s on %s",
		i.issue.Author.String(),
		i.issue.CreatedAt.Format("2006-01-02 15:04"),
	)
}

// FilterValue is required for itemWrapper to be a list.Item.
func (i itemWrapper) FilterValue() string {
	return i.Title() + "   " + i.Description()
}
