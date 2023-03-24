package config

import (
	"fmt"
	"os"

	"gn/config"
	style "gn/tui/style/config"

	tea "github.com/charmbracelet/bubbletea"
)

func Remove() {
	m := setUp("Which remote do you want to delete?", onSubmitRemove)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func onSubmitRemove(m *model) string {
	index := m.list.Index()

	selected, ok := m.list.Items()[index].(item)
	if !ok {
		return style.QuitText.Render("Failed to convert list.Item to item")
	}

	err := config.Remove(&m.oldConfig, index)
	if err != nil {
		return style.QuitText.Render(fmt.Sprintf("Failed to remove remote: %s", err))
	}

	return style.QuitText.Render(fmt.Sprintf("Sucessfully deleted the remote %s\nRemember to delete the API key on Gitlab", selected.lab.URL.String()))
}
