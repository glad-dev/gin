package color

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/glad-dev/gin/config"
	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/style"
	"github.com/glad-dev/gin/style/color"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Config is the entry point of this TUI, which allows to configure the colors.
func Config() {
	wrapper, err := config.Load()
	if err != nil {
		style.PrintErrAndExit("Failed to load config: " + err.Error())
	}

	initStyles()

	p := tea.NewProgram(model{
		inputs:              getInputs(),
		wrapper:             wrapper,
		currentlyDisplaying: displayingColors,
		state:               stateRunning,
	})

	m, err := p.Run()
	if err != nil {
		style.PrintErrAndExit("Failed to start program: " + err.Error())
	}

	if r, ok := m.(model); ok && r.state == exitFailure {
		logger.Log.Errorf(strings.TrimSpace(r.text))
		os.Exit(1)
	}
}

func getInputs() []textinput.Model {
	inputs := make([]textinput.Model, 3)
	r := regexp.MustCompile("^#[0-9a-fA-F]{0,6}$")

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.CharLimit = 7
		t.Placeholder = "E.g. #DEADBE"
		t.Validate = func(s string) error {
			if len(s) == 0 {
				return nil
			}

			if s == "#" {
				return nil
			}

			if !r.MatchString(s) {
				return errors.New("invalid")
			}

			return nil
		}

		switch i {
		case 0:
			t.Focus()
			t.PromptStyle = style.Focused
			t.TextStyle = style.Focused

			t.SetValue(string(color.Blurred))

		case 1:
			t.SetValue(string(color.Border))

		case 2:
			t.SetValue(string(color.Focused))
		}

		inputs[i] = t
	}

	return inputs
}
