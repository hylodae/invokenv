args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

all: vet lint build

golangci-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

build:
	go build -o bin/invokenv main.go

run:
	go run main.go $(call args,defaultstring)

vet:
	go vet ./...

lint: golangci-lint
	golangci-lint run ./...

test:
	go test -v ./...
