package all

import (
	"fmt"
	"os"

	"gn/issues"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(allIssues []issues.Issue, projectPath string) {
	issueList := make([]list.Item, len(allIssues))
	for i, issue := range allIssues {
		issueList[i] = issue
	}

	delegate := newItemDelegate()
	m := model{list: list.New(issueList, delegate, 0, 0)}
	m.list.Title = fmt.Sprintf("Issues of %s", projectPath)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
