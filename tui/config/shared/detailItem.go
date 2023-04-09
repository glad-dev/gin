package shared

import (
	"fmt"
	"io"
	"strings"

	"gn/tui/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type DetailItem struct {
	Username  string
	TokenName string
}

func (d DetailItem) FilterValue() string {
	return ""
}

type DetailsItemDelegate struct{}

func (d DetailsItemDelegate) Height() int                             { return 1 }
func (d DetailsItemDelegate) Spacing() int                            { return 0 }
func (d DetailsItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d DetailsItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(DetailItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. Username: '%s', Token name: '%s'", index+1, i.Username, i.TokenName)

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + strings.Join(s, ""))
		}
	}

	fmt.Fprint(w, fn(str))
}
