.PHONY: default test clean install install-g fmt vet run help

BINARY     := kwcli
MODULE     := github.com/shawn0915/kwcli
BUILDTIME  := $(shell date '+%Y-%m-%d %H:%M:%S %z')
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
LDFLAGS    := -ldflags "-X '$(MODULE)/cmd.buildTime=$(BUILDTIME)' -X '$(MODULE)/cmd.commitHash=$(COMMIT)'"

TSBS_REPO   := https://gitcode.com/mydb/kwdb-tsbs
TSBS_DIR    := third_party/kwdb-tsbs
TSBS_BINS   := tsbs_generate_data tsbs_generate_queries tsbs_load_kwdb tsbs_run_queries_kwdb
TSBS_BIN_DIR := bin

# Default target - build kwcli and TSBS
default: check-tsbs build-tsbs
	go build $(LDFLAGS) -o $(BINARY) .

# Check if TSBS source exists, download if not
check-tsbs:
	@if [ ! -d "$(TSBS_DIR)/cmd" ] || [ -z "$(shell ls -A $(TSBS_DIR)/cmd 2>/dev/null)" ]; then \
		echo "Downloading TSBS source code..."; \
		rm -rf $(TSBS_DIR); \
		git clone --depth=1 $(TSBS_REPO) $(TSBS_DIR); \
	fi

build-tsbs:
	@if [ -d "$(TSBS_DIR)/cmd" ]; then \
		mkdir -p $(TSBS_BIN_DIR); \
		(cd $(TSBS_DIR) && go mod tidy 2>/dev/null || true); \
		for bin in $(TSBS_BINS); do \
			echo "Building $$bin..."; \
			(cd $(TSBS_DIR) && go build -o ../../$(TSBS_BIN_DIR)/$$bin ./cmd/$$bin) || exit 1; \
		done; \
	else \
		echo "Warning: $(TSBS_DIR) is empty. TSBS binaries will not be built."; \
	fi

test:
	go test -v ./...

clean:
	rm -f $(BINARY)
	rm -rf $(TSBS_BIN_DIR)
	rm -rf third_party

install:
	@if [ ! -f "$(BINARY)" ]; then \
	  echo "Error: $(BINARY) not found. Please run 'make' first."; \
	  exit 1; \
	fi
	@echo "Installing $(BINARY) to $(HOME)/.kwcli/bin..."
	@mkdir -p $(HOME)/.kwcli/bin
	@cp $(BINARY) $(HOME)/.kwcli/bin/$(BINARY)
	@if [ -d "$(TSBS_BIN_DIR)" ]; then \
	  echo "Installing TSBS binaries to $(HOME)/.kwcli/bin..."; \
	  cp $(TSBS_BIN_DIR)/* $(HOME)/.kwcli/bin/; \
	fi
	@echo "Done. Ensure $(HOME)/.kwcli/bin is in your PATH."

install-g:
	@if [ ! -f "$(BINARY)" ]; then \
	  echo "Error: $(BINARY) not found. Please run 'make' first."; \
	  exit 1; \
	fi
	@echo "Installing $(BINARY) to /usr/local/bin..."
	@cp $(BINARY) /usr/local/bin/$(BINARY)
	@if [ -d "$(TSBS_BIN_DIR)" ]; then \
	  echo "Installing TSBS binaries to /usr/local/bin..."; \
	  cp $(TSBS_BIN_DIR)/* /usr/local/bin/; \
	fi
	@echo "Done. Run '$(BINARY) -v' to verify."

fmt:
	go fmt ./...

vet:
	go vet ./...

run: default
	./$(BINARY)

help:
	@echo "Available targets:"
	@echo "  (default)    Build kwcli and TSBS binaries (same as 'make')"
	@echo "  test         Run all tests"
	@echo "  clean        Remove built binaries"
	@echo "  install      Install to ~/.kwcli/bin"
	@echo "  install-g    Install to /usr/local/bin (requires sudo if needed)"
	@echo "  fmt          Format Go source files"
	@echo "  vet          Run go vet"
	@echo "  run          Build and run $(BINARY)"
	@echo "  help         Show this help message"