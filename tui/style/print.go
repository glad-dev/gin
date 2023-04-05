package style

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func FormatQuitText(str string) string {
	// Get the terminal width
	maxWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		maxWidth = 80
	}

	cp := quitText.Copy()
	cp.Width(maxWidth)

	return cp.Render(str)
}

func PrintErrAndExit(str string) {
	fmt.Fprint(os.Stderr, FormatQuitText(str))
	os.Exit(1)
}
