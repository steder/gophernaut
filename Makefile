all:
	go generate
	go build -o bin/gophernaut ...cmd/gophernaut
	go test -v $(go list ./... | grep -v /vendor/) -coverpkg github.com/steder/gophernaut

deps:
	go get -u github.com/Masterminds/glide
	go get -u golang.org/x/tools/cmd/stringer
	glide install
