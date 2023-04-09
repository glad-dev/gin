package edit

import (
	"fmt"
	"io"
	"strings"

	"gn/config"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type editListItem struct {
	remote *config.Remote
}

func (i editListItem) FilterValue() string { return "" }

type editItemDelegate struct{}

func (d editItemDelegate) Height() int                             { return 1 }
func (d editItemDelegate) Spacing() int                            { return 0 }
func (d editItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d editItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(editListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.remote.URL.String())

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + strings.Join(s, ""))
		}
	}

	fmt.Fprint(w, fn(str))
}
