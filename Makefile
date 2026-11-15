args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`
version = $(shell git describe --always --dirty)
build_date = $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")

all: vet lint build

golangci-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

build:
	go build -ldflags="-X main.Version=$(version) -X main.BuildTime=$(build_date)" -o bin/invokenv main.go

run:
	go run main.go $(call args,defaultstring)

vet:
	go vet ./...

lint: golangci-lint
	golangci-lint run ./...

test:
	go test -v ./...
