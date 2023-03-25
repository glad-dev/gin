package all

import (
	"fmt"

	"gn/issues"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(allIssues []issues.Issue, projectPath string) {
	issueList := make([]list.Item, len(allIssues))
	for i, issue := range allIssues {
		issueList[i] = item{issue: issue}
	}

	lst := list.New(issueList, newItemDelegate(), 0, 0)
	lst.Title = fmt.Sprintf("Issues of %s", projectPath)
	lst.SetSpinner(spinner.Points)

	m := model{
		list: lst,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}
}
