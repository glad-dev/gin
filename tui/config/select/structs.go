package selectconfig

import (
	"fmt"
	"io"
	"strings"

	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type itemWrapper struct {
	item remote.Details
}

// FilterValue is required for itemWrapper to be a list.Item.
func (i itemWrapper) FilterValue() string { return "" }

type selectItemDelegate struct{}

// Height is required for selectItemDelegate to be a list.ItemDelegate.
func (d selectItemDelegate) Height() int { return 1 }

// Spacing is required for selectItemDelegate to be a list.ItemDelegate.
func (d selectItemDelegate) Spacing() int { return 0 }

// Update is required for selectItemDelegate to be a list.ItemDelegate.
func (d selectItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render is required for selectItemDelegate to be a list.ItemDelegate.
func (d selectItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(itemWrapper)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s - %s", index+1, i.item.Username, i.item.TokenName)

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + strings.Join(s, ""))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}
