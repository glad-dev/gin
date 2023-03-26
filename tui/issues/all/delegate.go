package all

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var iid string

		if i, ok := m.SelectedItem().(itemWrapper); ok {
			iid = i.issue.Iid
		} else {
			return nil
		}

		selectItem := key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "View more details on selected issue"),
		)

		switch msg := msg.(type) { //nolint:gocritic
		case tea.KeyMsg:
			switch { //nolint:gocritic
			case key.Matches(msg, selectItem):
				return m.NewStatusMessage("You chose " + iid)
			}
		}

		return nil
	}

	return d
}
