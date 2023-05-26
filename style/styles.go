package style

import (
	"github.com/glad-dev/gin/style/color"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	Title         lipgloss.Style
	Item          lipgloss.Style
	SelectedItem  lipgloss.Style
	Pagination    lipgloss.Style
	Help          lipgloss.Style
	PrintOnlyList lipgloss.Style
	List          lipgloss.Style
	ListDetails   lipgloss.Style
	InputField    lipgloss.Style
	Focused       lipgloss.Style
	None          lipgloss.Style
	Error         lipgloss.Style

	Comment    lipgloss.Style
	Discussion lipgloss.Style

	quitText lipgloss.Style
)

// Init initialized the styles with the global color variables.
func Init() {
	Title = lipgloss.NewStyle().MarginLeft(2)
	Item = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(color.Focused)
	Pagination = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	Help = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	PrintOnlyList = lipgloss.NewStyle().PaddingLeft(4)
	List = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.NormalBorder()).Padding(1, 0, 1, 0)
	ListDetails = lipgloss.NewStyle().PaddingLeft(6)
	InputField = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	Focused = lipgloss.NewStyle().Foreground(color.Focused)
	None = lipgloss.NewStyle()
	Error = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

	Comment = lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.RoundedBorder()).Padding(1)
	Discussion = Comment.Copy().PaddingRight(0)

	quitText = lipgloss.NewStyle().Padding(1, 2, 1, 2)
}
