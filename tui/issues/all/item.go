package all

import (
	"fmt"

	"gn/issues/list"
)

// itemWrapper is a wrapper for list.Issue that implements all functions required by the list.Item interface.
type itemWrapper struct {
	issue list.Issue
}

func (i itemWrapper) Title() string {
	return fmt.Sprintf(
		"#%s %s",
		i.issue.Iid,
		i.issue.Title,
	)
}

func (i itemWrapper) Description() string {
	// Use author and creation date as description
	return fmt.Sprintf(
		"Created by %s on %s",
		i.issue.Author.String(),
		i.issue.CreatedAt.Format("2006-01-02 15:04"),
	)
}

func (i itemWrapper) FilterValue() string {
	return i.Title() + "   " + i.Description()
}
