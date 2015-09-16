
godep:
	go get github.com/tools/godep
	godep restore ./...

test: godep
	godep go test ./...

