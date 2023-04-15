package single

import (
	"fmt"
	"net/url"
	"os"

	"gn/repo"
	"gn/tui/issues/shared"
	"gn/tui/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details, u *url.URL, issueID string) {
	p := tea.NewProgram(
		model{
			content: "",
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
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}

	if m, ok := m.(model); ok && m.failure {
		fmt.Fprintf(os.Stderr, m.content)
		os.Exit(1)
	}
}
