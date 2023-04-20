package shared

import (
	"fmt"
	"io"
	"strings"

	"gn/config"
	"gn/logger"
	"gn/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem struct {
	Remote *config.Remote
}

func (i ListItem) FilterValue() string { return "" }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(ListItem)
	if !ok {
		logger.Log.Error("Got item that is not a ListItem.", "item", item)

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
