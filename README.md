# Git Issue Navigator (gin)

View the issues foa given repository on your command line.

## Known issues

- All issues:
	- List items that are too wide, break the view
	- Filtering does not show text since the UI would break otherwise
- Single issue:
	- List items that are too wide, look weird => Solution: Change placement, avoid lipgloss.Place

### Won't do

- Repo based query without token (How do I get the correct details?)

## v2

- [ ] Add support for PR
- [ ] Add support for Bitbucket, ...

# Acknowledgements

This project would not have been possible without [charmbracelet](https://github.com/charmbracelet) and their
[bubbletea](https://github.com/charmbracelet/bubbletea) and [glow](https://github.com/charmbracelet/glow) libraries.
