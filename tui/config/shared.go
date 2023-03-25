package config

import (
	"fmt"
	"strings"

	"gn/constants"
	"gn/tui/style"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = style.Focused.Copy()

	focusedButton = style.Focused.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func GetTextInputs() []textinput.Model {
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
			lipgloss.JoinVertical(
				lipgloss.Left,
				"Gitlab URL",
				style.InputField.Render(inputs[0].View()),
			),
			"\n",
			lipgloss.JoinVertical(
				lipgloss.Left,
				"API Key",
				style.InputField.Render(inputs[1].View()),
			),
			"\n",
			lipgloss.JoinVertical(
				lipgloss.Left,
				*button,
			),
		),
	)
}
