package shared

import (
	"fmt"
	"io"
	"strings"

	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// DetailItem is a list.Item and contains the user and token name of a remote.
type DetailItem struct {
	Username  string
	TokenName string
}

// FilterValue is required for DetailItem to be a list.Item.
func (d DetailItem) FilterValue() string {
	return ""
}

// DetailsItemDelegate is a list.ItemDelegate for DetailItem.
type DetailsItemDelegate struct{}

// Height is required for selectItemDelegate to be a list.ItemDelegate.
func (d DetailsItemDelegate) Height() int { return 1 }

// Spacing is required for selectItemDelegate to be a list.ItemDelegate.
func (d DetailsItemDelegate) Spacing() int { return 0 }

// Update is required for selectItemDelegate to be a list.ItemDelegate.
func (d DetailsItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render is required for selectItemDelegate to be a list.ItemDelegate.
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
