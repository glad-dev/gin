package single

import (
	"net/url"

	"gn/repo"
	"gn/style"
	"gn/tui/issues/shared"
	"gn/tui/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details, u *url.URL, issueID string) {
	conf, err := shared.SelectConfig(details)
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

	if m, ok := m.(model); ok && m.failure {
		style.PrintErrAndExit(m.content)
	}
}
