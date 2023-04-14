package style

import (
	"gn/style/color"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	Title        = lipgloss.NewStyle().MarginLeft(2)
	Item         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(color.Focused)
	Pagination   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	Help         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	List         = lipgloss.NewStyle().PaddingLeft(4)
	ListDetails  = lipgloss.NewStyle().PaddingLeft(6)
	InputField   = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	Focused      = lipgloss.NewStyle().Foreground(color.Focused)
	None         = lipgloss.NewStyle()
	Error        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

	quitText = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

var (
	Comment    = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.RoundedBorder()).Padding(1)
	Discussion = Comment.Copy()
)
