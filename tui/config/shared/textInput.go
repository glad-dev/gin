package shared

import (
	"fmt"
	"strings"

	"gn/constants"
	"gn/style"
	"gn/style/color"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	blurredStyle lipgloss.Style
	cursorStyle  lipgloss.Style

	focusedButton string
	blurredButton string
)

func initStyles() {
	blurredStyle = lipgloss.NewStyle().Foreground(color.Blurred)
	cursorStyle = style.Focused.Copy()

	focusedButton = style.Focused.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
}

func GetTextInputs() []textinput.Model {
	initStyles()

	inputs := make([]textinput.Model, 2)

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "E.g. https://gitlab.com"
			t.Focus()
			t.PromptStyle = style.Focused
			t.TextStyle = style.Focused
		case 1:
			t.Placeholder = "Requires the following scopes: " + strings.Join(constants.RequiredScopes, ", ")
		}

		inputs[i] = t
	}

	return inputs
}

func RenderInputFields(inputs []textinput.Model, focusIndex int, width int, height int) string {
	button := &blurredButton
	if focusIndex == len(inputs) {
		button = &focusedButton
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Left,

			"Gitlab URL",
			style.InputField.Render(inputs[0].View()),
			"\n",
			"API Key",
			style.InputField.Render(inputs[1].View()),
			"\n",
			*button,
		),
	)
}
