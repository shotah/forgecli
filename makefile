include .env
export

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
	go test -v -vet=all ./...

.PHONY: release
release:
	goreleaser build --single-target --skip-before --skip-validate --rm-dist

.PHONY: release-deploy
release-deploy:
	echo $(GITHUB_TOKEN)
	goreleaser release --skip-before --skip-validate --rm-dist
