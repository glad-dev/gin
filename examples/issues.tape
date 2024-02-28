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

# Open a repo
Type "./gin issues --url 'https://github.com/charmbracelet/bubbletea'"
Enter
## Wait for issues to load
Sleep 1s
Show
Sleep 1.5s

# Select an open issue
PageDown@1s 1
Down@250ms 5
Enter
## Wait for issue discussion to load
Sleep 1.5s
## Scroll down
PageDown@1.75s 2
Sleep 1.5s
## Go back
Type "q"
Sleep 1.5s

# Delete the binary
Hide
Type "q"
Type "rm gin"
Enter
