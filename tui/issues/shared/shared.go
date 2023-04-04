package shared

import (
	"gn/repo"
	"net/url"

	"github.com/charmbracelet/bubbles/spinner"
)

type Shared struct {
	IssueID string
	URL     *url.URL
	Details []repo.Details
	Spinner spinner.Model
}
