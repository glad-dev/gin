package single

import (
	"net/url"

	"github.com/glad-dev/gin/repo"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/issues/shared"
	"github.com/glad-dev/gin/tui/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details, u *url.URL, issueID string) {
	conf, err := shared.SelectConfig(details, u)
	if err != nil {
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
