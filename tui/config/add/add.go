package add

import (
	"fmt"
	"os"
	"strings"

	"gn/constants"
	"gn/tui/style/color"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Add() {
	m := model{
		inputs:          make([]textinput.Model, 2),
		inputFieldStyle: lipgloss.NewStyle().BorderForeground(color.Border).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "E.g. https://gitlab.com"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Requires the following scopes: " + strings.Join(constants.RequiredScopes, ", ")
		}

		m.inputs[i] = t
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
