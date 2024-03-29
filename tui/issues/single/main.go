package single

import (
	"net/url"

	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/repository"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/issues/shared"
	"github.com/glad-dev/gin/tui/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

// Show is the entry point of this TUI, which displays a single issues of a given repository.
func Show(details []repository.Details, u *url.URL, issueID string) {
	conf, err := shared.SelectConfig(details, u)
	if err != nil {
		log.Error("Failed to select config", "error", err)

		style.PrintErrAndExit("Failed to select config: " + err.Error())
	}

	p := tea.NewProgram(
		model{
			content: "",
			conf:    conf,
			shared: &shared.Shared{
				Details: details,
				URL:     u,
				IssueID: issueID,
				Spinner: *widgets.GetSpinner(),
			},
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // turn on mouse support, so we can track the mouse wheel
	)

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("could not run program:" + err.Error())
	}

	if r, ok := m.(model); ok && r.failure {
		style.PrintErrAndExit(r.content)
	}
}
