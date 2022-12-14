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
	staticcheck ./...

.PHONY: test
test:
	go test -cover -vet=all ./...

.PHONY: test-failed
test-failed:
	go test -failfast -v -cover ./...

.PHONY: coverage
coverage:
# go gocov convert coverage.out | go gocov-xml > coverage.xml
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: release
release:
	goreleaser build --single-target --skip-before --skip-validate --rm-dist

.PHONY: release-deploy
release-deploy:
# goreleaser release --rm-dist --config .goreleaser-actions.yml --parallelism 6 --skip-before --skip-validate
	goreleaser release --parallelism 6 --skip-before --skip-validate --rm-dist --debug
