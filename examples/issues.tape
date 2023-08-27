Output issues.gif

Require go

Set Shell "fish"
Set FontSize 20
Set Width 1750
Set Height 1250

# Build the program
Hide
Type "cd .."
Enter
Type "go build -o gin"
Enter
Type "clear"
Enter
Show

# Open a repo
Type "./gin issues --url 'https://github.com/charmbracelet/bubbletea'"
Sleep 500ms
Enter
## Wait for issues to load
Sleep 5s

# Select an open issue
Down@1s 2
Enter
## Wait for issue discussion to load
Sleep 2s
## Scroll down
Down@250ms 8
Sleep 1s
## Go back
Type "q"

# Select a closed issue
Right@1s
Sleep 1s
Down@250ms 4
Enter
## Wait for issue discussion to load
Sleep 2s
Down@250ms 10
Sleep 1s
## Go back
Type "q"

## Show all issues
Right@1s
Sleep 1s
Down@250ms 10
Sleep 2s

# Delete the binary
Hide
Type "q"
Type "rm gin"
Enter
