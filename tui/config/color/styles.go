package color

import (
	"fmt"

	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/style/color"

	"github.com/charmbracelet/lipgloss"
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
