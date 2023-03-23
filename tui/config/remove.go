package config

import (
	"fmt"
	"os"

	"gn/config"
	style "gn/tui/style/config"

	tea "github.com/charmbracelet/bubbletea"
)

func Remove() {
	m := setUp("Which remote do you want to delete?", removeView)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func removeView(m *model) string {
	if m.quitting {
		return style.QuitText.Render("No changes were made.")
	}

	if m.action {
		index := m.list.Index()

		selected, ok := m.list.Items()[index].(item)
		if !ok {
			return style.QuitText.Render("Failed to convert list.Item to item")
		}

		err := config.RemoveOnce(&m.oldConfig, index) // This is called multiple times?
		if err != nil {
			return style.QuitText.Render(fmt.Sprintf("Failed to remove remote: %s", err))
		}

		return style.QuitText.Render(fmt.Sprintf("Sucessfully deleted the remote %s\nRemember to delete the API key on Gitlab", selected.lab.URL.String()))
	}

	return "\n" + m.list.View()
}
