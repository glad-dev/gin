# Git Issue Navigator (gin)

View the issues of a GitHub or GitLab repository directly from your command line.

![Gif showing gin's functionality](examples/issues.gif)

## Installation
```go
go install github.com/glad-dev/gin@latest
```

## Usage

### Configuring authentication tokens

Authenticate with the GiHub/GitLab APIs by importing tokens with ``gin config add``.

#### Creating a GitLab token

Create a token with scopes ``read_api``, ``read_user`` and ``read_repository`` in ``Preferences > Access Tokens``. 

#### Creating a GitHub token

Create a **classic** token with scopes ``repo`` and ``read_user`` in ``Settings > Developer settings > Personal access token``.

### Viewing issues

- To view the issues of a remote repository, use ``gin issues --url https://github.com/path/to/your/repo``.
- To view the issues of the repository you're currently in, use ``gin issues``

## Acknowledgements

This project would not have been possible without [charmbracelet](https://github.com/charmbracelet) and their
[bubbletea](https://github.com/charmbracelet/bubbletea) and [glow](https://github.com/charmbracelet/glow) libraries.
