package all

import (
	"gn/style/color"

	"github.com/charmbracelet/lipgloss"
)

var (
	inactiveTabBorder lipgloss.Border
	activeTabBorder   lipgloss.Border
	inactiveTabStyle  lipgloss.Style
	activeTabStyle    lipgloss.Style

	windowStyle lipgloss.Style
	docStyle    lipgloss.Style
)

func initStyles() {
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder = tabBorderWithBottom("┘", " ", "└")
	inactiveTabStyle = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(color.Border).Padding(0, 1)
	activeTabStyle = inactiveTabStyle.Copy().Border(activeTabBorder, true)

	windowStyle = lipgloss.NewStyle().BorderForeground(color.Border).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
	docStyle = lipgloss.NewStyle().Margin(1, 2, 1, 2)
}
