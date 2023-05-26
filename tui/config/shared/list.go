package shared

import (
	"github.com/glad-dev/gin/style"

	"github.com/charmbracelet/bubbles/list"
)

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

func RenderList(list list.Model) string {
	cp := style.List.Copy()
	cp.Width(list.Width())

	return cp.Render(list.View())
}
