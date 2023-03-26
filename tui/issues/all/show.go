package all

import (
	"gn/repo"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func Show(details []repo.Details) {
	lst := list.New([]list.Item{}, newItemDelegate(), 0, 0)
	lst.Title = ""

	m := model{
		list:    lst,
		details: details,
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		style.PrintErrAndExit("Error running program: " + err.Error())
	}
}
