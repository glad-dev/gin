package all

import (
	"net/url"
	"os"

	"gn/issues"
	"gn/repo"
	"gn/tui/issues/shared"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details, u *url.URL) {
	lst := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	lst.SetShowFilter(false) // TODO: Fix filtering bug

	lsts := [3]list.Model{}
	for i := range lsts {
		lsts[i] = lst
	}

	s := spinner.New()
	s.Spinner = spinner.Points

	p := tea.NewProgram(
		model{
			tabs: tabs{
				lists:     lsts,
				activeTab: 0,
			},
			shared: &shared.Shared{
				IssueID: "",
				URL:     u,
				Details: details,
				Spinner: s,
			},
			viewport:     viewport.New(0, 0),
			viewedIssues: make(map[string]issues.IssueDetails),
			isLoading:    true,
			viewingList:  true,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}

	if m, ok := m.(model); ok && m.failure {
		os.Exit(1)
	}
}
