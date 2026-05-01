.PHONY: build test clean install fmt vet run help

BINARY     := kwcli
MODULE     := github.com/shawn0915/kwcli
BUILDTIME  := $(shell date '+%Y-%m-%d %H:%M:%S')
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS    := -ldflags "-X '$(MODULE)/cmd.buildTime=$(BUILDTIME)' -X '$(MODULE)/cmd.commitHash=$(COMMIT)'"

# Default target
build:
	go build $(LDFLAGS) -o $(BINARY) .

test:
	go test -v ./...

clean:
	rm -f $(BINARY)

install: build
	@echo "Installing $(BINARY) to /usr/local/bin..."
	@cp $(BINARY) /usr/local/bin/$(BINARY)
	@echo "Done. Run '$(BINARY) -v' to verify."

install-local: build
	@echo "Installing $(BINARY) to ~/.local/bin..."
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY) $(HOME)/.local/bin/$(BINARY)
	@echo "Done. Ensure $(HOME)/.local/bin is in your PATH."

fmt:
	go fmt ./...

vet:
	go vet ./...

run: build
	./$(BINARY)

help:
	@echo "Available targets:"
	@echo "  build        Build the $(BINARY) binary with build time and commit hash injected"
	@echo "  test         Run all tests"
	@echo "  clean        Remove built binary"
	@echo "  install      Install to /usr/local/bin (requires sudo if needed)"
	@echo "  install-local Install to ~/.local/bin"
	@echo "  fmt          Format Go source files"
	@echo "  vet          Run go vet"
	@echo "  run          Build and run $(BINARY)"
	@echo "  help         Show this help message"
