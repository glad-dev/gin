package widgets

import "github.com/charmbracelet/bubbles/spinner"

func GetSpinner() *spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Points

	return &s
}
