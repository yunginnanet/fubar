test:
	go test -cover -v ./...

gif:
	vhs assets/fubar.tape
	git add assets/fubar.gif
