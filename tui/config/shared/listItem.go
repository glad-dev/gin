package shared

import (
	"fmt"
	"io"
	"strings"

	"github.com/glad-dev/gin/configuration"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ListItem is a list.Item and contains a *configuration.Remote.
type ListItem struct {
	Remote *configuration.Remote
}

// FilterValue is required for ListItem to be a list.Item.
func (i ListItem) FilterValue() string { return "" }

// ItemDelegate is a list.ItemDelegate for ListItem.
type ItemDelegate struct{}

// Height is required for selectItemDelegate to be a list.ItemDelegate.
func (d ItemDelegate) Height() int { return 1 }

// Spacing is required for selectItemDelegate to be a list.ItemDelegate.
func (d ItemDelegate) Spacing() int { return 0 }

// Update is required for selectItemDelegate to be a list.ItemDelegate.
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render is required for selectItemDelegate to be a list.ItemDelegate.
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(ListItem)
	if !ok {
		log.Error("Got item that is not a ListItem.", "item", item)

		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Remote.URL.String())
	if len(str) > m.Width() && m.Width() > 3 {
		width := m.Width() - style.List.GetHorizontalFrameSize() - style.ListDetails.GetHorizontalFrameSize() - 8 // TODO: Get correct value

		str = str[:width] + "..."
	}

	fn := style.Item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return style.SelectedItem.Render("> " + strings.Join(s, ""))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}
