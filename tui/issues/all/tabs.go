package all

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var tabTitles = []string{
	"Open",
	"Closed",
	"All",
}

func renderTab(m *model) string {
	doc := strings.Builder{}

	renderedTabs := make([]string, len(tabTitles))

	activeList := m.tabs.lists[m.tabs.activeTab]

	// width per item
	itemWidth := (activeList.Width() - len(tabTitles[0]) - len(tabTitles[1]) - len(tabTitles[2])) / 3

	contentWidth := 0
	for i, t := range tabTitles {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabTitles)-1, i == m.tabs.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}

		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive { //nolint:gocritic
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)

		spaces := strings.Repeat(" ", max(0, itemWidth-len(t))/2)
		toPrint := fmt.Sprintf("%s%s%s", spaces, t, spaces)

		contentWidth += len(toPrint)
		if i == len(tabTitles)-1 && contentWidth < activeList.Width() {
			toPrint += strings.Repeat(" ", max(0, activeList.Width()-contentWidth-12))
		}

		renderedTabs[i] = style.Render(toPrint)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")

	doc.WriteString(
		windowStyle.Width(activeList.Width() - windowStyle.GetHorizontalFrameSize()).
			Render(activeList.View()),
	)

	return docStyle.Render(doc.String())
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right

	return border
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
