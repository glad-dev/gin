package selectconfig

import (
	"errors"
	"fmt"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var ErrUserQuit = errors.New("user quit selection program")

// Select is the entry point of this TUI, which allows the user to choose a remote if more than one exists.
func Select(selectedRemote *configuration.Remote, title string) (*remote.Details, error) {
	items := make([]list.Item, len(selectedRemote.Details))
	for i := range selectedRemote.Details {
		items[i] = itemWrapper{
			item: selectedRemote.Details[i],
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

	return &item.item, nil
}
