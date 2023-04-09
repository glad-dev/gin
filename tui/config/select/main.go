package selectconfig

import (
	"errors"
	"fmt"

	"gn/config"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var ErrUserQuit = errors.New("user quit selection program")

func Select(selectedRemote config.Remote, title string) (int, error) {
	items := make([]list.Item, len(selectedRemote.Details))
	for i, conf := range selectedRemote.Details {
		items[i] = selectListItem{
			username:  conf.Username,
			tokenName: conf.TokenName,
		}
	}

	lst := list.New(items, selectItemDelegate{}, 0, 0)

	if len(title) == 0 {
		title = "Select the corresponding token"
	}

	lst.Title = title
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	p := tea.NewProgram(model{
		list: lst,
	})

	m, err := p.Run()
	if err != nil {
		return -1, fmt.Errorf("failed to run select program: %w", err)
	}

	mod, ok := m.(model)
	if !ok {
		return -1, fmt.Errorf("failed to convert m to model")
	}

	if mod.quitting {
		return -1, ErrUserQuit
	}

	return mod.list.Index(), nil
}
