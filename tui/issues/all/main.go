package all

import (
	"net/url"

	"gn/issues/single"
	"gn/repo"
	"gn/style"
	"gn/tui/issues/shared"
	"gn/tui/widgets"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details, u *url.URL) {
	conf, err := shared.SelectConfig(details, u)
	if err != nil {
		style.PrintErrAndExit("Failed to select config: " + err.Error())
	}
	initStyles()

	lst := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	lst.SetShowFilter(false) // TODO: Fix filtering bug

	lsts := [3]list.Model{}
	for i := range lsts {
		lsts[i] = lst
	}

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
				Spinner: *widgets.GetSpinner(),
			},
			conf:                conf,
			viewport:            viewport.New(0, 0),
			viewedIssues:        make(map[string]single.IssueDetails),
			currentlyDisplaying: displayingInitialLoading,
			state:               stateRunning,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}

	if m, ok := m.(model); ok && m.state == exitFailure {
		style.PrintErrAndExit(m.error)
	}
}
