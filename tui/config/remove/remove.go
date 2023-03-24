package remove

import (
	"fmt"
	"os"

	"gn/config"
	style "gn/tui/style/config"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func Config() {
	// Load current config
	wrapper, err := config.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, style.QuitText.Render(fmt.Sprintf("Failed to load config: %s", err)))
		os.Exit(1)
	}

	if len(wrapper.Configs) == 0 {
		fmt.Fprint(os.Stderr, style.QuitText.Render("The config file contains no remotes."))
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
	lst.Title = "Which remote do you want to delete?"
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	m := model{
		list:      lst,
		oldConfig: *wrapper,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
