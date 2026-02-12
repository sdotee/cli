.PHONY: build clean test darwin_universal release

VERSION=1.0.1
BIN=see
DOCKER_CMD=docker

GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X main.BuildVersion=$(VERSION) -X 'main.BuildTime=`date`' -extldflags -static"
GO=$(GO_ENV) $(shell which go)

build: main.go
	@$(GO_ENV) $(GO) build $(GO_FLAGS) -o $(BIN) main.go

darwin_universal: main.go
	@GOOS=darwin GOARCH=arm64 $(GO_ENV) $(GO) build $(GO_FLAGS) -o $(BIN)_arm64 main.go
	@GOOS=darwin GOARCH=amd64 $(GO_ENV) $(GO) build $(GO_FLAGS) -o $(BIN)_amd64 main.go
	@lipo -create -output $(BIN) $(BIN)_arm64 $(BIN)_amd64
	@rm -f $(BIN)_arm64 $(BIN)_amd64

build_docker_image: clean
	@$(DOCKER_CMD) build -f ./Dockerfile -t see-cli:$(VERSION) .

install: build
	@$(GO) install $(GO_FLAGS) main.go

test:
	@$(GO) test ./...

# clean all build result
clean:
	@$(GO) clean ./...
	@rm -f $(BIN)
	@rm -rf dist/

# run goreleaser to release snapshot
release:
	@goreleaser release --snapshot --clean

all: clean test build
