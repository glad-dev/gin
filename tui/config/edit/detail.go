package edit

import (
	"fmt"
	"io"
	"strings"

	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type detail struct {
	username  string
	tokenName string
}

func (d detail) FilterValue() string {
	return ""
}

type detailsItemDelegate struct{}

func (d detailsItemDelegate) Height() int                             { return 1 }
func (d detailsItemDelegate) Spacing() int                            { return 0 }
func (d detailsItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d detailsItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(detail)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. Username: '%s', Token name: '%s'", index+1, i.username, i.tokenName)

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + strings.Join(s, ""))
		}
	}

	fmt.Fprint(w, fn(str))
}
