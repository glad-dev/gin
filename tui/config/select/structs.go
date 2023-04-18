package selectconfig

import (
	"fmt"
	"io"
	"strings"

	"gn/config/remote"
	"gn/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type itemWrapper struct {
	item remote.Details
}

func (i itemWrapper) FilterValue() string { return "" }

type selectItemDelegate struct{}

func (d selectItemDelegate) Height() int                             { return 1 }
func (d selectItemDelegate) Spacing() int                            { return 0 }
func (d selectItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d selectItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(itemWrapper)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s - %s", index+1, i.item.GetUsername(), i.item.GetTokenName())

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + strings.Join(s, ""))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}
