# Git Issue Navigator (gin)

View the issues for a given repository on your command line.

![Gif showing gin's functionality](examples/issues.gif)

## Creating tokens

### GitLab

Preferences > Access Tokens > New Token with scopes "read_api", "read_user", "read_repository"

### GitHub

Settings > Developer settings > Personal access token > Generate new token (classic) with scopes public_repo and read_user.

# Acknowledgements

This project would not have been possible without [charmbracelet](https://github.com/charmbracelet) and their
[bubbletea](https://github.com/charmbracelet/bubbletea) and [glow](https://github.com/charmbracelet/glow) libraries.
