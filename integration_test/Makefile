BATS := bats
BATS_FORMATTER ?= pretty
BATS_VERBOSE ?=
BINARY_PATH := ../bin/mycli

# Build Bats flags
BATS_FLAGS := --formatter $(BATS_FORMATTER)
ifdef BATS_VERBOSE
    BATS_FLAGS += --verbose-run
endif

.PHONY: all test-root test-configure test-echo test-cat check-binary check-bats clean

# Default target: run all tests
all: check-binary check-bats
	@echo "Running all integration tests..."
	@$(BATS) $(BATS_FLAGS) root.bats configure.bats echo.bats cat.bats

# Run root command tests only
test-root: check-binary check-bats
	@echo "Running root command tests..."
	@$(BATS) $(BATS_FLAGS) root.bats

# Run configure command tests only
test-configure: check-binary check-bats
	@echo "Running configure command tests..."
	@$(BATS) $(BATS_FLAGS) configure.bats

# Run echo command tests only
test-echo: check-binary check-bats
	@echo "Running echo command tests..."
	@$(BATS) $(BATS_FLAGS) echo.bats

# Run cat command tests only
test-cat: check-binary check-bats
	@echo "Running cat command tests..."
	@$(BATS) $(BATS_FLAGS) cat.bats

# Check that mycli binary exists
check-binary:
	@if [ ! -f "$(BINARY_PATH)" ]; then \
		echo "Error: Binary not found at $(BINARY_PATH)"; \
		echo "Please run 'make build' from project root first."; \
		exit 1; \
	fi

# Check that bats is installed
check-bats:
	@if ! command -v $(BATS) >/dev/null 2>&1; then \
		echo "Error: bats is not installed"; \
		echo ""; \
		echo "Install bats-core:"; \
		echo "  macOS:  brew install bats-core"; \
		echo "  Linux:  sudo apt-get install bats"; \
		echo "  Manual: https://github.com/bats-core/bats-core#installation"; \
		exit 1; \
	fi

# Clean target (nothing to clean in test directory)
clean:
	@echo "Nothing to clean in integration_test/"
