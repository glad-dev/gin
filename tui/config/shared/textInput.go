package shared

import (
	"fmt"

	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/style/color"

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

// GetTextInputs returns two textinput.Model. The first model is meant for a URL. The second model is meant for a token.
func GetTextInputs() []textinput.Model {
	initStyles()

	inputs := make([]textinput.Model, 2)

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64

		switch i {
		case 0:
			t.Placeholder = "E.g. https://gitlab.com"
			t.Focus()
			t.PromptStyle = style.Focused
			t.TextStyle = style.Focused
		case 1:
			t.Placeholder = "E.g. bfgQeo8JAnMSz4SnDq8l"
		}

		inputs[i] = t
	}

	return inputs
}

// RenderInputFields renders the input fields created by GetTextInputs.
func RenderInputFields(inputs []textinput.Model, apiDetails string, focusIndex int, width int, height int) string {
	button := &blurredButton
	if focusIndex == len(inputs) {
		button = &focusedButton
	}

	if len(apiDetails) > 0 {
		apiDetails = fmt.Sprintf(" (requires following permissions: %s)", apiDetails)
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Left,

			"Base URL",
			style.InputField.Render(inputs[0].View()),
			"\n",
			"API Key"+apiDetails,
			style.InputField.Render(inputs[1].View()),
			"\n",
			*button,
		),
	)
}
