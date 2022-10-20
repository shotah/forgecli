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
	golangci-lint --version
	golangci-lint

.PHONY: test
test:
	go test ./...
