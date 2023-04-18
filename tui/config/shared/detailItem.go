package shared

import (
	"fmt"
	"io"
	"strings"

	"gn/logger"
	"gn/style"

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
		logger.Log.Error("Got item that is not a DetailItem.", "item", item)

		return
	}

	str := fmt.Sprintf("%d. Username: '%s', Token name: '%s'", index+1, i.Username, i.TokenName)

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + strings.Join(s, ""))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}
