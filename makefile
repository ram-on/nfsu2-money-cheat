build:
	go build nfsu2-money-cheat.go

build-win:
	GOOS=windows GOARCH=386 go build nfsu2-money-cheat.go

clean:
	rm -f nfsu2-money-cheat nfsu2-money-cheat.exe

