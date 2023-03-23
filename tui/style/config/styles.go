package config

import (
	"gn/tui/style/color"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	Title        = lipgloss.NewStyle().MarginLeft(2)
	Item         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(color.Focused)
	Pagination   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	Help         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	List         = list.DefaultStyles().HelpStyle.PaddingLeft(4)
	QuitText     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
