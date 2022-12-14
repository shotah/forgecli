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
# goreleaser release --rm-dist --config .goreleaser-actions.yml --parallelism 6 --skip-before --skip-validate
	goreleaser release --parallelism 6 --skip-before --skip-validate --rm-dist --debug
