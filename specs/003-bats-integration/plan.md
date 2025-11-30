# Implementation Plan: Bats Integration Testing Framework

**Branch**: `003-bats-integration` | **Date**: 2025-11-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-bats-integration/spec.md`

## Summary

統合テストフレームワークをBatsを使用して実装します。この機能により、開発者はビルド済みの`bin/mycli`バイナリに対して、実環境に近い形で全コマンド（root, configure, echo）の動作を検証できます。テストはサブコマンドごとに分離され、個別実行または全体実行が可能で、CI/CDパイプラインにも統合されます。テスト環境は一時ディレクトリによる分離戦略を採用し、開発者の設定ファイルとの競合を防ぎます。

## Technical Context

**Language/Version**: Bash shell scripts (POSIX compatible), Bats (Bash Automated Testing System)  
**Primary Dependencies**: 
- Bats-core (test framework) - 推奨 v1.10.0+
- GNU Make (build automation)
- bash 4.0+ or zsh (shell environment)
**Storage**: Temporary directories (`mktemp -d` pattern), no persistent storage  
**Testing**: Bats test framework for integration tests, existing Go unit tests remain unchanged  
**Target Platform**: macOS, Linux (any UNIX-like OS with bash/zsh)  
**Project Type**: CLI application - integration test layer addition  
**Performance Goals**: Complete test suite execution under 30 seconds  
**Constraints**: 
- Tests must not modify user's actual configuration files
- Each test run must be isolated (unique temp directories)
- Tests must work from any working directory
**Scale/Scope**: 
- 3 bats test files (root.bats, configure.bats, echo.bats)
- Estimated 10-20 test cases per command
- Shared helper scripts for common operations

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **TDD必須**: Batsテスト自体がテストコードであり、既存の動作を検証する。新しいヘルパー関数にはテストを含める
- [x] **パッケージ責務分離**: 統合テストは外部から実行するため、既存の`cmd/`と`internal/`分離には影響しない
- [x] **コード品質基準**: Bashスクリプトには`shellcheck`による検証を推奨（オプション）、Goコードの品質基準には影響なし
- [x] **設定管理の一貫性**: テストは一時ディレクトリを使用し、`~/.config/mycli/`の実ファイルには触れない。環境変数で設定パスをオーバーライド
- [x] **ユーザーエクスペリエンス**: テストは既存のCobraインターフェースの動作を検証するものであり、UXには影響しない
- [x] **パフォーマンス要件**: 統合テスト全体で30秒以内（個別のCLI起動時間は既存の100ms要件を維持）

**特記事項**: この機能は既存実装に新しいコードを追加せず、外部からの統合テストレイヤーを追加するため、憲章のすべての原則に準拠しています。

## Project Structure

### Documentation (this feature)

```text
specs/003-bats-integration/
├── spec.md              # Feature specification (created)
├── plan.md              # This file (implementation plan)
├── research.md          # Phase 0: Technology research and decision rationale
├── data-model.md        # Phase 1: Test suite structure and test case data model
├── quickstart.md        # Phase 1: Quick start guide for running integration tests
├── contracts/           # Phase 1: Test behavior contracts
│   ├── root-tests.md    # Root command test scenarios
│   ├── configure-tests.md # Configure command test scenarios
│   └── echo-tests.md    # Echo command test scenarios
└── tasks.md             # Phase 2: Task breakdown (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
go-cli/                              # Project root
├── integration_test/                # NEW: Integration test directory
│   ├── Makefile                     # NEW: Test execution targets
│   ├── helpers/                     # NEW: Shared test utilities
│   │   ├── common.bash              # NEW: Common setup/teardown functions
│   │   ├── assertions.bash          # NEW: Custom assertion functions
│   │   └── test_env.bash            # NEW: Environment setup utilities
│   ├── root.bats                    # NEW: Root command integration tests
│   ├── configure.bats               # NEW: Configure command integration tests
│   ├── echo.bats                    # NEW: Echo command integration tests
│   └── README.md                    # NEW: Integration test documentation
├── Makefile                         # MODIFIED: Add integration-test targets
├── .github/
│   └── workflows/
│       └── ci.yaml                  # MODIFIED: Add integration test stage
├── bin/
│   └── mycli                        # EXISTING: Binary to test (built by make build)
├── cmd/                             # EXISTING: No changes
├── internal/                        # EXISTING: No changes
└── [other existing files]           # EXISTING: No changes
```

**Structure Decision**: 

統合テストは既存のGoコード構造とは独立した`integration_test/`ディレクトリに配置します。この決定により：

1. **明確な分離**: 統合テスト（ブラックボックス）と単体テスト（ホワイトボックス）の責務が明確
2. **独立した実行**: Batsフレームワークによる実行は、Go tool chainと干渉しない
3. **ツール選択の柔軟性**: 将来的に他のテストツールを追加する際も、ディレクトリ構造が拡張可能
4. **CI/CD統合の容易さ**: 専用のMakefileにより、テスト実行を独立してスケジュール可能

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

該当なし - すべての憲章原則に準拠しています。

---

## Phase 0: Research & Outline

### Research Topics

1. **Batsフレームワークの選定と評価**
   - Bats-core vs alternatives (shunit2, roundup)
   - インストール方法とバージョン管理
   - macOS/Linux互換性

2. **テスト分離戦略の実装パターン**
   - 一時ディレクトリ生成（`mktemp -d`）
   - 環境変数のオーバーライド手法
   - クリーンアップの確実性（trapハンドラー）

3. **Makefileベースのテスト統合**
   - 複数batsファイルの実行戦略
   - 個別テスト実行のターゲット設計
   - 失敗時の出力フォーマット

4. **CI/CD統合のベストプラクティス**
   - GitHub Actions上でのBatsインストール
   - テスト結果の可視化
   - 並列実行の可能性

5. **詳細出力モードの実装**
   - Batsの`-t`（tap形式）と`--formatter pretty`オプション
   - 環境変数による出力制御
   - CI環境での適切な出力形式

### Deliverables

`research.md` ファイルに以下を含める：

- **決定事項**: 各研究トピックに対する選択とその理由
- **技術的根拠**: ベンチマークやコミュニティでの採用状況
- **代替案の評価**: 検討した他の手法と却下理由
- **リスクと緩和策**: 予見される問題と対処法

---

## Phase 1: Design & Contracts

### Data Model

`data-model.md` に以下の構造を定義：

#### Test Suite Structure

```yaml
TestSuite:
  name: string                    # "root", "configure", "echo"
  file: string                    # "root.bats", "configure.bats", "echo.bats"
  helpers: [string]               # ["helpers/common.bash", "helpers/assertions.bash"]
  setup: function                 # Test setup logic
  teardown: function              # Test cleanup logic
  tests: [TestCase]

TestCase:
  description: string             # Test case description
  given: string                   # Precondition
  when: string                    # Action
  then: string                    # Expected outcome
  assertions: [Assertion]

Assertion:
  type: enum                      # "status", "output", "file_exists", "env_var"
  expected: any                   # Expected value
  actual: any                     # Actual value (computed)
```

#### Test Environment Model

```yaml
TestEnvironment:
  temp_dir: string                # Unique temp directory for this test
  config_home: string             # Overridden config directory
  env_vars:
    MYCLI_CONFIG: string          # Points to temp config
    HOME: string                  # May be overridden for isolation
  binary_path: string             # Path to bin/mycli
  cleanup_required: boolean       # Flag for teardown
```

#### Helper Functions Model

```yaml
HelperFunction:
  name: string                    # Function name
  description: string             # What it does
  parameters: [Parameter]
  returns: type                   # Return type (status code, string, etc.)
  example: string                 # Usage example

CommonHelpers:
  - setup_test_env()              # Create temp dirs, set env vars
  - teardown_test_env()           # Clean up temp dirs
  - run_mycli(args...)            # Execute binary with args
  - assert_success()              # Check $status == 0
  - assert_failure()              # Check $status != 0
  - assert_output(expected)       # Check $output matches
  - assert_line(index, expected)  # Check specific line in output
```

### API Contracts

`contracts/` ディレクトリに以下を作成：

#### `contracts/root-tests.md`

```markdown
# Root Command Test Contract

## Test Cases

### TC-ROOT-001: Display Help
- **Given**: No arguments provided
- **When**: `mycli` is executed
- **Then**: Help message is displayed, exit status 0

### TC-ROOT-002: Display Version
- **Given**: `--version` flag is provided
- **When**: `mycli --version` is executed
- **Then**: Version string is displayed, exit status 0

### TC-ROOT-003: Invalid Flag
- **Given**: Invalid flag is provided
- **When**: `mycli --invalid-flag` is executed
- **Then**: Error message is displayed, exit status 1

[Add more test cases as needed]
```

#### `contracts/configure-tests.md`

```markdown
# Configure Command Test Contract

## Test Cases

### TC-CONF-001: Create New Config
- **Given**: No existing config file
- **When**: `mycli configure` is executed and editor saves valid config
- **Then**: Config file created at expected path, exit status 0

### TC-CONF-002: Edit Existing Config
- **Given**: Existing config file
- **When**: `mycli configure` is executed
- **Then**: Editor opens with existing content, changes are saved

### TC-CONF-003: Profile Support
- **Given**: `--profile test` flag is provided
- **When**: `mycli configure --profile test` is executed
- **Then**: Profile-specific config file is created

[Add more test cases as needed]
```

#### `contracts/echo-tests.md`

```markdown
# Echo Command Test Contract

## Test Cases

### TC-ECHO-001: Basic Output
- **Given**: String argument provided
- **When**: `mycli echo "Hello"` is executed
- **Then**: "Hello" followed by newline is output, exit status 0

### TC-ECHO-002: Multiple Arguments
- **Given**: Multiple string arguments provided
- **When**: `mycli echo Hello World` is executed
- **Then**: "Hello World" followed by newline is output

### TC-ECHO-003: No Trailing Newline (-n flag)
- **Given**: `-n` flag is provided
- **When**: `mycli echo -n "Hello"` is executed
- **Then**: "Hello" without newline is output

### TC-ECHO-004: Escape Sequences (-e flag)
- **Given**: `-e` flag and escape sequences provided
- **When**: `mycli echo -e "Line1\nLine2"` is executed
- **Then**: Two lines are output with actual newline

[Add more test cases as needed]
```

### Quickstart Guide

`quickstart.md` に以下を記載：

```markdown
# Integration Test Quickstart

## Prerequisites

1. Install Bats:
   ```bash
   # macOS
   brew install bats-core
   
   # Linux (manual install)
   git clone https://github.com/bats-core/bats-core.git
   cd bats-core
   sudo ./install.sh /usr/local
   ```

2. Build the application:
   ```bash
   make build
   ```

## Running Tests

### Run all integration tests
```bash
make integration-test
```

### Run specific command tests
```bash
# From project root
make integration-test-root
make integration-test-configure
make integration-test-echo

# Or directly with bats
bats integration_test/root.bats
```

### Verbose output
```bash
BATS_VERBOSE=1 make integration-test
```

## Writing New Tests

1. Choose appropriate test file (root.bats, configure.bats, echo.bats)
2. Load helper functions:
   ```bash
   load helpers/common
   load helpers/assertions
   ```
3. Write test case:
   ```bash
   @test "TC-XXX-NNN: Description" {
     setup_test_env
     run_mycli arg1 arg2
     assert_success
     assert_output "expected output"
     teardown_test_env
   }
   ```

## Troubleshooting

- **Binary not found**: Run `make build` first
- **Tests fail**: Check `bin/mycli --version` works
- **Permission errors**: Ensure `bin/mycli` is executable
```

### Agent Context Update

After Phase 1 completion, run:

```bash
.specify/scripts/bash/update-agent-context.sh copilot
```

This will update `.github/copilot-instructions.md` with integration test information.

---

## Phase 2: Task Breakdown

⚠️ **この段階は `/speckit.tasks` コマンドで実行されます**

Phase 2のタスク分解は、Phase 0とPhase 1の完了後に別のコマンドで実行されます。

期待されるタスクカテゴリ：

1. **Infrastructure Setup Tasks**
   - Batsインストールドキュメント作成
   - `integration_test/`ディレクトリ構造作成
   - Makefileターゲット追加

2. **Helper Script Implementation Tasks**
   - `helpers/common.bash`実装
   - `helpers/assertions.bash`実装
   - `helpers/test_env.bash`実装

3. **Test Implementation Tasks**
   - `root.bats`実装（優先度P1）
   - `configure.bats`実装（優先度P2）
   - `echo.bats`実装（優先度P2）

4. **CI/CD Integration Tasks**
   - `.github/workflows/ci.yaml`更新
   - Bats installation stepを追加
   - Integration test stageを追加

5. **Documentation Tasks**
   - `integration_test/README.md`作成
   - プロジェクトルート README 更新

詳細は `/speckit.tasks` 実行時に生成されます。
