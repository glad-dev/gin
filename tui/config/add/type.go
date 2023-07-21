package add

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/glad-dev/gin/remote"
	"github.com/glad-dev/gin/style"
)

const (
	btnGithub    = "[ GitHub ]"
	btnGitLab    = "[ GitLab ]"
	btnBitbucket = "[ Bitbucket ]"
)

func updateType(m *model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = exitNoChange

			return tea.Quit

			// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			if s == "enter" {
				m.currentlyDisplaying = displayingAdd

				switch m.focusIndex {
				case 0:
					m.focusIndex = 1
					m.remoteType = remote.Github
					m.inputs[0].Placeholder = "https://github.com"
					m.inputs[0].SetValue("https://github.com")

				case 1:
					m.focusIndex = 0
					m.remoteType = remote.Gitlab
					m.inputs[0].Placeholder = "https://gitlab.com"

				case 2:
					m.focusIndex = 0
					m.remoteType = remote.Bitbucket
					m.inputs[0].Placeholder = "https://bitbucket.com"
				}

				return updateInputs(m)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			return nil
		}
	}

	return nil
}

func viewType(m *model) string {
	btns := [3]string{btnGithub, btnGitLab, btnBitbucket}
	btns[m.focusIndex] = style.Focused.Render(btns[m.focusIndex])

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,

		lipgloss.JoinVertical(
			lipgloss.Center,
			"Select the remote type:", // ToDo: Add clearer line
			"\n",
			btns[0],
			"\n",
			btns[1],
			"\n",
			btns[2],
		),
	)
}