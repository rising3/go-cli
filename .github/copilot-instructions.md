# Copilot Instructions for go-cli

## Repository Overview

This is a Go CLI application template built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper). The project provides a scaffold for building command-line tools with configuration management, profile support, and editor integration.

- **Language**: Go 1.25.4
- **CLI Framework**: Cobra v1.10.1
- **Configuration**: Viper v1.21.0
- **Binary Name**: `mycli`

**重要**: このプロジェクトの開発原則・品質基準は `.specify/memory/constitution.md` で定義されています。すべての実装・レビューは憲章の原則に従う必要があります。

## Build Instructions

### Prerequisites

1. **Go 1.25.4** - Required version specified in `go.mod`
2. **golangci-lint v2.6.2** - Required for linting. Install with:
   ```bash
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.2
   ```
3. **PATH Configuration** - Always ensure `$(go env GOPATH)/bin` is in PATH before running lint:
   ```bash
   export PATH="$(go env GOPATH)/bin:$PATH"
   ```

### Commands

| Command | Description | Typical Duration |
|---------|-------------|------------------|
| `go mod download` | Download dependencies | ~5 seconds |
| `make build` | Build binary to `bin/mycli` | ~2 seconds |
| `make test` | Run all tests | ~5 seconds |
| `make fmt` | Format code with gofmt | ~1 second |
| `make lint` | Run golangci-lint (requires PATH setup) | ~5 seconds |
| `make all` | Run test → fmt → lint → build | ~15 seconds |
| `make clean` | Remove `bin/` directory | ~1 second |

### Recommended Workflow

```bash
# 1. Set up PATH for golangci-lint
export PATH="$(go env GOPATH)/bin:$PATH"

# 2. Install golangci-lint if not present
which golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.2

# 3. Build and validate all
make all
```

### Common Issues

- **`make: golangci-lint: No such file or directory`**: Install golangci-lint and ensure PATH includes `$(go env GOPATH)/bin`
- **Dependencies not found**: Run `go mod download` before building

## Project Structure

```
go-cli/
├── main.go                    # Entry point - calls cmd.Execute()
├── cmd/                       # CLI commands (Cobra)
│   ├── root.go               # Root command, config initialization
│   ├── configure.go          # configure subcommand
│   ├── viperutils.go         # Viper configuration utilities
│   └── *_test.go             # Command tests
├── internal/                  # Internal packages
│   ├── cmd/                  # Internal command logic
│   │   └── configure.go      # ConfigureFile implementation
│   ├── editor/               # Editor detection (EDITOR env, OS defaults)
│   ├── proc/                 # Process execution utilities
│   └── stdio/                # Standard I/O stream utilities
├── Makefile                   # Build targets
├── go.mod                     # Go module definition
└── .github/workflows/ci.yaml  # CI pipeline
```

## Key Files

- **`cmd/root.go`**: Defines `CliName`, `CliVersion`, `Config` struct, and configuration loading logic
- **`cmd/configure.go`**: Implements `configure` command for creating config files
- **`internal/cmd/configure.go`**: Core logic for config file creation and editor invocation
- **`internal/editor/editor.go`**: Editor detection with OS-specific fallbacks

## CI Pipeline

The `.github/workflows/ci.yaml` runs on PRs to `main`, `next`, and `feature/**` branches:

1. Checkout code
2. Set up Go 1.25.x
3. Install golangci-lint v2.6.2
4. `make build`
5. `make test`
6. `make fmt`
7. `make lint`

**Always run `make all` locally before pushing to ensure CI passes.**

## Testing

- Tests use `t.TempDir()` for isolation
- Tests override globals with `t.Cleanup()` for restoration
- `internal/proc.ExecCommand` can be mocked in tests
- Run specific tests: `go test -run TestName ./...`
- Run with verbose output: `go test -v ./...`

## Configuration

- Default config path: `~/.config/mycli/default.yaml`
- Profile configs: `~/.config/mycli/<profile>.yaml`
- Environment variables: `MYCLI_CONFIG`, `MYCLI_PROFILE`
- Config struct fields use `mapstructure` tags with kebab-case

## Adding New Commands

1. Create new file in `cmd/` directory (e.g., `cmd/newcmd.go`)
2. Define command using `&cobra.Command{}`
3. Register in `init()` with `rootCmd.AddCommand()`
4. Add tests in `cmd/newcmd_test.go`
5. Run `make all` to validate

## Code Style

- Use `gofmt -s -w .` for formatting
- Follow Go conventions for package naming
- Internal packages go in `internal/`
- Use `mapstructure` tags for config struct fields
