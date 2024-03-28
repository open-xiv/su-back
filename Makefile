.PHONY: env-show build clean

# environment variables
ROOT = $(PWD)
SRC = $(ROOT)/cmd/su
BIN = $(ROOT)/build

# print env
env-show:
	@echo "ROOT = $(ROOT), SRC = $(SRC), BIN = $(BIN)"

# build check
build-check:
	@ls $(BIN)

# default command: only build
all: build

# build for all platforms
build: env-show build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-check

# build for linux amd64
build-linux-amd64:
	@echo "Building for linux/amd64"
	@mkdir -p build
	@cd $(SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN)/su_amd64

# build on linux arm64
build-linux-arm64:
	@echo "Building for linux/arm64"
	@mkdir -p build
	@cd $(SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(BIN)/su_arm64

# build for macOS with Intel
build-darwin-amd64:
	@echo "Building for darwin/amd64"
	@mkdir -p build
	@cd $(SRC) && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BIN)/su_darwin_amd64

build-darwin-arm64:
	@echo "Building for darwin/arm64"
	@mkdir -p build
	@cd $(SRC) && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(BIN)/su_darwin_arm64


clean:
	@echo "Cleaning up"
	@go clean
