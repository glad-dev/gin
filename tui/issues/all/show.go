package all

import (
	"gn/issues"
	"gn/repo"
	shared "gn/tui/issues"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details) {
	lsts := []list.Model{
		list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}

	s := spinner.New()
	s.Spinner = spinner.Points

	m := model{
		tabs: tabs{
			lists:     [3]list.Model(lsts),
			activeTab: 0,
		},
		shared: &shared.Shared{
			IssueID: "",
			Details: details,
			Spinner: s,
		},
		viewedIssues: make(map[string]issues.IssueDetails),
		isLoading:    true,
		viewingList:  true,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}
}
