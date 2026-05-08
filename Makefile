.PHONY: build test clean install fmt vet run help

BINARY     := kwcli
MODULE     := github.com/shawn0915/kwcli
BUILDTIME  := $(shell date '+%Y-%m-%d %H:%M:%S %z')
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS    := -ldflags "-X '$(MODULE)/cmd.buildTime=$(BUILDTIME)' -X '$(MODULE)/cmd.commitHash=$(COMMIT)'"

# TSBS build configuration
TSBS_DIR    := third_party/kwdb-tsbs
TSBS_BINS   := tsbs_generate_data tsbs_generate_queries tsbs_load_kwdb tsbs_run_queries_kwdb
TSBS_BIN_DIR := bin

# Default target
build: build-tsbs
	go build $(LDFLAGS) -o $(BINARY) .

build-tsbs:
	@if [ -d "$(TSBS_DIR)/cmd" ]; then \
		mkdir -p $(TSBS_BIN_DIR); \
		for bin in $(TSBS_BINS); do \
			echo "Building $$bin..."; \
			(cd $(TSBS_DIR) && go build -o ../../$(TSBS_BIN_DIR)/$$bin ./cmd/$$bin) || exit 1; \
		done; \
	else \
		echo "Warning: $(TSBS_DIR) is empty. TSBS binaries will not be built."; \
		echo "  Run: git clone --depth=1 https://github.com/KWDB/kwdb-tsbs.git $(TSBS_DIR)"; \
	fi

test:
	go test -v ./...

clean:
	rm -f $(BINARY)
	rm -rf $(TSBS_BIN_DIR)

install: build
	@echo "Installing $(BINARY) to /usr/local/bin..."
	@cp $(BINARY) /usr/local/bin/$(BINARY)
	@if [ -d "$(TSBS_BIN_DIR)" ]; then \
		@echo "Installing TSBS binaries to /usr/local/bin..."; \
		@cp $(TSBS_BIN_DIR)/* /usr/local/bin/; \
	fi
	@echo "Done. Run '$(BINARY) -v' to verify."

install-local: build
	@echo "Installing $(BINARY) to ~/.local/bin..."
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY) $(HOME)/.local/bin/$(BINARY)
	@if [ -d "$(TSBS_BIN_DIR)" ]; then \
		@echo "Installing TSBS binaries to ~/.local/bin..."; \
		@cp $(TSBS_BIN_DIR)/* $(HOME)/.local/bin/; \
	fi
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
