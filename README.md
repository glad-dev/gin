# To do

## v1

- [ ] Issues view all => Fix filtering view

## Known issues

- All issues:
	- List items that are too wide, break the view
	- Filtering does not show text since the UI would break otherwise
- Single issue:
	- List items that are too wide, look weird => Solution: Change placement, avoid lipgloss.Place

### Done

- All issues:
	- [x] Change order, newest first
- [x] GitHub: What if ID passed that is not an issue? => Query fails and returns NOT_FOUND
- [x] Check gitlab error checking code. Example url that should fail: https://gitlab.com/zerkc/whatsdesk/-/issues
- [x] Add GitHub support
- [x] Add selection for when multiple configs match
- [x] Color config
- [x] Config add show error in TUI, not after
- [x] Pagination for all issues
- [x] Config add/edit show error in TUI, not after
- [x] Loading screens during transition
- [x] Add cli command to check token validity
- [x] Log errors

### Not blocking

- Rendering the warning emoji leads to whitespace error

### Won't do

- Repo based query without token (How do I get the correct details?)

## v2

- [ ] Add PR support

# Acknowledgements

This project would not have been possible without [charmbracelet](https://github.com/charmbracelet) and their
[bubbletea](https://github.com/charmbracelet/bubbletea) and [glow](https://github.com/charmbracelet/glow) libraries.
