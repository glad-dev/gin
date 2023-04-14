package shared

import (
	"gn/style"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
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
	return lipgloss.Place(
		list.Width(),
		list.Height(),
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			style.InputField.Render(list.View()),
		),
	)
}
