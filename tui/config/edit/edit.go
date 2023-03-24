package edit

import (
	"fmt"
	"os"

	"gn/config"
	tui "gn/tui/config"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	// Load current config
	wrapper, err := config.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, style.FormatQuitText(fmt.Sprintf("Failed to load config: %s", err)))
		os.Exit(1)
	}

	items := make([]list.Item, len(wrapper.Configs))
	for i, conf := range wrapper.Configs {
		items[i] = item{
			lab: conf,
		}
	}

	// Not sure what these numbers do, but the TUI looks better with them
	const defaultWidth = 20
	const listHeight = 14

	lst := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	lst.Title = "Which remote do you want to edit?"
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	m := model{
		list:      lst,
		oldConfig: *wrapper,
		inputs:    tui.GetTextInputs(),
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
