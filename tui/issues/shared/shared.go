package shared

import (
	"gn/repo"

	"github.com/charmbracelet/bubbles/spinner"
)

type Shared struct {
	IssueID string
	Details []repo.Details
	Spinner spinner.Model
}
