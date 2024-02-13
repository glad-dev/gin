package shared

import (
	"net/url"

	"github.com/glad-dev/gin/repository"

	"github.com/charmbracelet/bubbles/spinner"
)

// Shared contains data that is used by the all issue and single issue TUI.
type Shared struct {
	IssueID string
	URL     *url.URL
	Details []repository.Details
	Spinner spinner.Model
}
