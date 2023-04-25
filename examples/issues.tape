Output issues.gif

Require go

Set Shell "fish"
Set FontSize 20

Hide
Type "cd .."
Enter
Type "go build -o gn"
Enter
Type "clear"
Enter
Show

Type "./gn issues --url 'https://github.com/charmbracelet/bubbletea'"
Sleep 500ms
Enter

Sleep 3s

Hide
Type "rm gn"
