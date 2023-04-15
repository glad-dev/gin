package color

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"gn/style"
	"gn/style/color"
)

var (
	blurredStyle  lipgloss.Style
	focusedButton string
	blurredButton string
)

func initStyles() {
	blurredStyle = lipgloss.NewStyle().Foreground(color.Blurred)
	focusedButton = style.Focused.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
}
