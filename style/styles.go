package style

import (
	"github.com/glad-dev/gin/style/color"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	Title         = lipgloss.NewStyle().MarginLeft(2)
	Item          = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItem  = lipgloss.NewStyle().PaddingLeft(2).Foreground(color.Focused)
	Pagination    = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	Help          = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	PrintOnlyList = lipgloss.NewStyle().PaddingLeft(4)
	List          = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.NormalBorder()).Padding(1, 0, 1, 0)
	ListDetails   = lipgloss.NewStyle().PaddingLeft(6)
	InputField    = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	Focused       = lipgloss.NewStyle().Foreground(color.Focused)
	None          = lipgloss.NewStyle()
	Error         = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

	Comment    = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.RoundedBorder()).Padding(1)
	Discussion = Comment.Copy().PaddingRight(0)

	quitText = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

func UpdateColors() {
	SelectedItem = SelectedItem.Foreground(color.Focused)
	List = List.BorderForeground(color.Border)
	InputField = InputField.BorderForeground(color.Border)
	Focused = Focused.Foreground(color.Focused)

	Comment = Comment.BorderForeground(color.Border)
	Discussion = Comment.Copy().PaddingRight(0)
}
