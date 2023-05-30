package shared

import (
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/bubbles/list"
)

// NewList returns a list given items, delegate and title. Furthermore, the status bar and filtering are disabled.
// It also sets styles for the title, pagination and help.
func NewList(items []list.Item, delegate list.ItemDelegate, title string) list.Model {
	lst := list.New(items, delegate, 0, 0)
	lst.Title = title
	lst.SetShowStatusBar(false)
	lst.SetFilteringEnabled(false)
	lst.Styles.Title = style.Title
	lst.Styles.PaginationStyle = style.Pagination
	lst.Styles.HelpStyle = style.Help

	return lst
}

// RenderList renders a list in the list style with a max width set by list.Width().
func RenderList(list list.Model) string {
	cp := style.List.Copy()
	cp.Width(list.Width())

	return cp.Render(list.View())
}
