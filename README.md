# Git Issue Navigator (gin)

View the issues foa given repository on your command line.

## Known issues

- All issues:
	- List items that are too wide, break the view
	- Filtering does not show text since the UI would break otherwise
	- Querying can take a lot of time if many issues exist
- Single issue:
	- List items that are too wide, look weird => Solution: Change placement, avoid lipgloss.Place

## v2

- [ ] Add support for PR

# Acknowledgements

This project would not have been possible without [charmbracelet](https://github.com/charmbracelet) and their
[bubbletea](https://github.com/charmbracelet/bubbletea) and [glow](https://github.com/charmbracelet/glow) libraries.
