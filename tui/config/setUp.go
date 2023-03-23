package config

import (
	"fmt"
	"os"

	"gn/config"
	style "gn/tui/style/config"

	"github.com/charmbracelet/bubbles/list"
)

func setUp(title string, viewFunction func(*model) string) model {
	// Load current config
	generalConfig, err := config.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, style.QuitText.Render(fmt.Sprintf("Failed to load config: %s", err)))
		os.Exit(1)
	}

	if len(generalConfig.Configs) == 0 {
		fmt.Fprint(os.Stderr, style.QuitText.Render("The config file contains no remotes."))
		os.Exit(1)
	}

	items := make([]list.Item, len(generalConfig.Configs))
	for i, conf := range generalConfig.Configs {
		items[i] = item{
			lab: conf,
		}
	}

	// Not sure what these numbers do, but the TUI looks better with them
	const defaultWidth = 20
	const listHeight = 14

	lst := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	lst.Title = title
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	return model{
		list:      lst,
		oldConfig: *generalConfig,
		view:      viewFunction,
	}
}