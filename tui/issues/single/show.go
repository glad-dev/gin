package single

import (
	"fmt"
	"os"

	"gn/repo"
	"gn/tui/issues"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details, issueID string) {
	s := spinner.New()
	s.Spinner = spinner.Points

	p := tea.NewProgram(
		model{
			content: "",
			shared: &issues.Shared{
				Details: details,
				IssueID: issueID,
				Spinner: s,
			},
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // turn on mouse support, so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
