package single

import (
	"fmt"
	"os"

	"gn/repo"

	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details, issueID string) {
	p := tea.NewProgram(
		model{
			content: "",
			shared: Shared{
				details: details,
				issueID: issueID,
			},
		},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
