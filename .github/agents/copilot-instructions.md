# go-cli Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-11-30

## Active Technologies
- Go 1.25.4 + Cobra v1.10.1, Viper v1.21.0, gopkg.in/yaml.v3 v3.0.1 (002-configure-refactor)
- Local filesystem (~/.config/mycli/*.yaml configuration files) (002-configure-refactor)
- Bash shell scripts (POSIX compatible), Bats (Bash Automated Testing System) (003-bats-integration)
- Temporary directories (`mktemp -d` pattern), no persistent storage (003-bats-integration)
- Go 1.25.4 + Cobra v1.10.1+（CLI フレームワーク）、Viper v1.21.0+（設定管理） (004-cat-subcommand)
- ファイルシステム（読み取り専用）、標準入力/出力 (004-cat-subcommand)

- Go 1.25.4（`go.mod`で指定、憲章で定義済み） (001-echo-subcommand)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.25.4（`go.mod`で指定、憲章で定義済み）

## Code Style

Go 1.25.4（`go.mod`で指定、憲章で定義済み）: Follow standard conventions

## Recent Changes
- 004-cat-subcommand: Added Go 1.25.4 + Cobra v1.10.1+（CLI フレームワーク）、Viper v1.21.0+（設定管理）
- 003-bats-integration: Added Bash shell scripts (POSIX compatible), Bats (Bash Automated Testing System)
- 002-configure-refactor: Added Go 1.25.4 + Cobra v1.10.1, Viper v1.21.0, gopkg.in/yaml.v3 v3.0.1


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
