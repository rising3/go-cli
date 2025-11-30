# Integration Tests

This directory contains Bats-based integration tests for the `mycli` CLI application.

## Prerequisites

### Install Bats

**macOS (recommended)**:
```bash
brew install bats-core
```

**Linux (Ubuntu/Debian)**:
```bash
sudo apt-get update
sudo apt-get install -y bats
```

**Manual Installation** (all platforms):
```bash
git clone https://github.com/bats-core/bats-core.git
cd bats-core
sudo ./install.sh /usr/local
```

**Verify Installation**:
```bash
bats --version
# Expected: Bats 1.10.0 or higher
```

### Build the Application

Before running integration tests, you must build the binary:

```bash
# From project root
make build
```

This creates the `bin/mycli` binary that the tests will execute.

## Running Tests

### Run all integration tests

From project root:
```bash
make integration-test
```

### Run individual command tests

```bash
# Root command tests only
make integration-test-root

# Configure command tests only
make integration-test-configure

# Echo command tests only
make integration-test-echo
```

### Run tests directly with Bats

From the `integration_test/` directory:
```bash
# All tests
bats *.bats

# Specific file
bats root.bats

# Specific test case
bats root.bats --filter "TC-ROOT-001"
```

## Output Modes

### Default (Pretty) Output

```bash
make integration-test
```

Shows concise pass/fail status for each test.

### Verbose Output

```bash
BATS_VERBOSE=1 make integration-test
```

Shows detailed command execution and output for each test.

### TAP Format (for CI/CD)

```bash
BATS_FORMATTER=tap make integration-test
```

Machine-readable Test Anything Protocol format.

## Test Structure

```
integration_test/
├── Makefile              # Test execution targets
├── README.md             # This file
├── helpers/              # Shared test utilities
│   ├── common.bash       # Basic test setup/teardown
│   ├── assertions.bash   # Custom assertion functions
│   └── test_env.bash     # Environment configuration
├── root.bats             # Root command tests
├── configure.bats        # Configure command tests
└── echo.bats             # Echo command tests
```

## Writing New Tests

Each test file follows this pattern:

```bash
#!/usr/bin/env bats

# Load helper functions
load helpers/common
load helpers/assertions
load helpers/test_env

# Run before each test
setup() {
    setup_test_env
}

# Run after each test
teardown() {
    teardown_test_env
}

# Test case
@test "TC-XXX-NNN: Description" {
    # Given: Setup preconditions
    # When: Execute command
    run_mycli command args
    
    # Then: Verify results
    assert_success
    assert_output "expected output"
}
```

## Test Isolation

Each test runs in an isolated environment:

- **Unique temporary directory**: `$TEST_TEMP_DIR` (auto-created and cleaned)
- **Isolated config home**: `$TEST_CONFIG_HOME` (separate from user's actual config)
- **Environment variable overrides**: `MYCLI_CONFIG`, `HOME` point to test directories
- **Automatic cleanup**: All temporary files removed after each test

This ensures tests never interfere with your actual `mycli` configuration files.

## Troubleshooting

### Error: "Binary not found at bin/mycli"

Run `make build` from project root first.

### Error: "bats: command not found"

Install Bats using instructions above.

### Error: "Permission denied"

Ensure binary is executable:
```bash
chmod +x bin/mycli
```

### Error: "Failed to read config file: Unsupported Config Type"

This occurs when `MYCLI_CONFIG` environment variable is set to a directory instead of a file path. The test helpers correctly set `HOME` and let `mycli` use its default config location at `$HOME/.config/mycli/`. Do not export `MYCLI_CONFIG` unless you want to test with a specific config file path.

### Tests not executing (shows "Executed 0 instead of expected N tests")

This was caused by using `trap cleanup_test_env EXIT` in `setup_test_env()` which conflicted with Bats' `run` command when testing the `echo` subcommand. The fix is to rely solely on the `teardown()` function for cleanup, not trap handlers.

### Tests fail unexpectedly

Run with verbose output to see details:
```bash
BATS_VERBOSE=1 make integration-test
```

## For More Information

- **Test Contracts**: See `specs/003-bats-integration/contracts/` for test case specifications
- **Quickstart Guide**: See `specs/003-bats-integration/quickstart.md` for detailed usage
- **Bats Documentation**: https://bats-core.readthedocs.io/
