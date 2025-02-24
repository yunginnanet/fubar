Output assets/fubar.gif
Require go
Set Shell "bash"
Set FontSize 20
Set Width 1250
Set Height 1300
Hide
Type@5ms "PS1=' > '"
Enter@10ms
Type@5ms "clear"
Enter@10ms
Show
Type "go test -v -run TestExampleCode/exit_1_has_name"
Sleep 1s
Enter
Sleep 10s
