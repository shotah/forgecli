.PHONY : upgrade
upgrade :
	go get -u .
	go mod download
	go mod tidy

.PHONY : build
build :
	go build

.PHONY : lint
lint : 
	golint ./...
	gofmt -w -s .
	revive -config revive.toml ./...
	staticcheck ./...

.PHONY: test
test:
	go test ./...
