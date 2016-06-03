all:
	go build -o bin/gophernaut ...cmd/gophernaut
	go test -v $(go list ./... | grep -v /vendor/) -coverpkg github.com/steder/gophernaut
