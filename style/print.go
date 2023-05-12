package style

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// FormatQuitText formats and returns the passed string with the style.quitText style. Word wrap is determined using
// term.GetSize.
func FormatQuitText(str string) string {
	// Get the terminal width
	maxWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		maxWidth = 80
	}

	cp := quitText.Copy()
	cp.Width(maxWidth)

	return cp.Render(str) + "\n"
}

// PrintErrAndExit formats the passed string using FormatQuitText, prints it to Stderr and calls os.Exit(1).
func PrintErrAndExit(str string) {
	_, _ = fmt.Fprint(os.Stderr, FormatQuitText(str))
	os.Exit(1)
}
