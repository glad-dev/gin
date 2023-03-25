package style

import (
	"fmt"
	"os"
	"strings"
)

func FormatQuitText(str string) string {
	// Max length of 80
	const maxLen = 80

	out := ""
	for _, line := range strings.Split(str, "\n") {
		if len(line) < maxLen {
			out += line + "\n"

			continue
		}

		for len(line) >= maxLen {
			// Get the right most space
			index := strings.LastIndex(line[:maxLen], " ")
			if index == -1 || index > maxLen {
				out += line[:maxLen-2] + "-\n"
				line = line[maxLen-2:]

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
