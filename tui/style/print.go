package style

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func FormatQuitText(str string) string {
	// Get the terminal width
	maxWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		maxWidth = 80
	}

	out := ""
	for _, line := range strings.Split(str, "\n") {
		if len(line) < maxWidth {
			out += line + "\n"

			continue
		}

		for len(line) >= maxWidth {
			// Get the right most space
			index := strings.LastIndex(line[:maxWidth], " ")
			if index == -1 || index > maxWidth {
				out += line[:maxWidth-2] + "-\n"
				line = line[maxWidth-2:]

				continue
			}

			out += line[:index] + "\n"
			line = line[index+1:]
		}

		out += line + "\n"
	}

	return "\n" + quitText.Render(out) + "\n"
}

func PrintErrAndExit(str string) {
	fmt.Fprint(os.Stderr, FormatQuitText(str))
	os.Exit(1)
}
