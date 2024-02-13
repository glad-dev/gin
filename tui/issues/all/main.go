package all

import (
	"net/url"

	"github.com/glad-dev/gin/issues/discussion"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/repository"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/tui/issues/shared"
	"github.com/glad-dev/gin/tui/widgets"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Show is the entry point of this TUI, which displays all issues of a given repository.
func Show(details []repository.Details, u *url.URL) {
	conf, err := shared.SelectConfig(details, u)
	if err != nil {
		log.Error("Failed to select config", "error", err)

		style.PrintErrAndExit("Failed to select config: " + err.Error())
	}
	initStyles()

	lst := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	lst.SetShowFilter(false) // TODO: Fix filtering bug

	lists := [3]list.Model{}
	for i := range lists {
		lists[i] = lst
	}

	p := tea.NewProgram(
		model{
			tabs: tabs{
				lists:     lists,
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
			viewedIssues:        make(map[string]discussion.Details),
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

	if r, ok := m.(model); ok && r.state == exitFailure {
		style.PrintErrAndExit(r.error)
	}
}
