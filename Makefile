all:
	go build -o bin/gophernaut ...cmd/gophernaut
	go test -v $(go list ./... | grep -v /vendor/) -coverpkg github.com/steder/gophernaut

deps:
	dep ensure

stringer: ${GOPATH}/bin/stringer
	go get golang.org/x/tools/cmd/stringer

generate: stringer
	go install
	go generate
