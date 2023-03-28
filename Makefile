default: build

build:
	go build -v ./...

# See https://golangci-lint.run/
lint:
	golangci-lint run -c .golangci.toml ./...

generate:
	go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

.PHONY: build lint generate fmt test
