package all

import (
	"gn/repo"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details) {
	lst := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	lst.Title = ""

	s := spinner.New()
	s.Spinner = spinner.Points

	m := model{
		list:      lst,
		details:   details,
		isLoading: true,
		spinner:   s,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}
}
