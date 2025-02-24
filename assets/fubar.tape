Output assets/fubar.gif
Require go
Set Shell "bash"
Set FontSize 24
Set Width 1150
Set Height 1300
Hide
Type@5ms "PS1=' > '"
Enter@10ms
Type@5ms "clear"
Enter@10ms
Show
Type "go test -run TestPanicModule"
Sleep 1s
Enter
Sleep 10s
