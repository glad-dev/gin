package selectconfig

import (
	"errors"
	"fmt"

	"gn/config"
	"gn/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var ErrUserQuit = errors.New("user quit selection program")

func Select(selectedRemote *config.Remote, title string) (*config.RemoteDetails, error) {
	items := make([]list.Item, len(selectedRemote.Details))
	for i := range selectedRemote.Details {
		items[i] = itemWrapper{
			item: &selectedRemote.Details[i],
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
		return nil, fmt.Errorf("failed to run select program: %w", err)
	}

	mod, ok := m.(model)
	if !ok {
		return nil, fmt.Errorf("failed to convert m to model")
	}

	if mod.state == exitNoSelection {
		return nil, ErrUserQuit
	}

	item, ok := mod.list.SelectedItem().(itemWrapper)
	if !ok {
		return nil, fmt.Errorf("failed to selected item to selectListItem")
	}

	return item.item, nil
}
