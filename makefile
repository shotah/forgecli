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
	go vet ./...
	gofmt -w -s .
	revive -config revive.toml ./...
	staticcheck ./...

.PHONY: test
test:
	go test ./...

.PHONY: choco-build
choco-build:
# Hash url: https://github.com/shotah/forgecli/releases/download/latest/forgecli-latest-windows-amd64.zip.md5
# Zip url: https://github.com/shotah/forgecli/releases/download/latest/forgecli-latest-windows-amd64.zip
	curl.exe -s $(RELEASEHASHURL) > md5.txt
