# Implementation Plan: Configure サブコマンドのリファクタリング

**Branch**: `002-configure-refactor` | **Date**: 2025-11-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-configure-refactor/spec.md`

**Note**: This plan implements the refactoring of the configure subcommand to follow echo subcommand patterns and Cobra best practices.

## Summary

Refactor the `configure` subcommand to align with echo subcommand implementation patterns and Cobra best practices. The primary goals are: (1) Remove dependency on `internal/stdio` package, (2) Use Cobra's standard I/O streams (`cmd.OutOrStdout()`, `cmd.ErrOrStderr()`), (3) Separate command entry point (`cmd/configure.go`) from business logic (`internal/cmd/configure/configure.go`), (4) Define `ConfigureOptions` struct similar to `EchoOptions`, and (5) Implement test function variable pattern (`ConfigureFunc`) for testability.

## Technical Context

**Language/Version**: Go 1.25.4  
**Primary Dependencies**: Cobra v1.10.1, Viper v1.21.0, gopkg.in/yaml.v3 v3.0.1  
**Storage**: Local filesystem (~/.config/mycli/*.yaml configuration files)  
**Testing**: Go standard testing framework (`go test`), table-driven tests  
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows)  
**Project Type**: Single CLI project (Cobra-based command structure)  
**Performance Goals**: CLI startup <100ms, config file read <10ms, help display <50ms  
**Constraints**: File permissions 0644 for config files, backward compatibility with existing behavior  
**Scale/Scope**: Single refactoring feature affecting 4 files (cmd/configure.go, internal/cmd/configure/*.go, tests)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **TDD必須**: すべての実装に対してテストを先に書く計画があるか？
  - ✅ FR-011: 既存テストの継続パス、FR-012: 新規テストケース追加を要求
  - ✅ User Story 5 (P5): テストカバレッジ80%以上を目標
  - ✅ テストファースト戦略: `internal/cmd/configure/configure_test.go`と`cmd/configure_test.go`を更新/追加

- [x] **パッケージ責務分離**: `cmd/`（CLI）と`internal/`（内部ロジック）が明確に分離されているか？
  - ✅ User Story 2 (P2): `internal/cmd/configure/`パッケージの作成とロジック分離を明示
  - ✅ FR-002: `cmd/configure.go`はフラグ取得とオプション構築のみ
  - ✅ FR-003: `internal/cmd/configure.Configure()`に再利用可能なロジックを配置

- [x] **コード品質基準**: `gofmt`と`govet`による検証を通過する見込みか？
  - ✅ SC-008: `make lint`がゼロ件の警告・エラーで終了することを成功基準として定義
  - ✅ リファクタリングのため、既存コードと同等の品質を維持

- [x] **設定管理の一貫性**: Viperを使用し、`~/.config/mycli/`配下の設定ファイル構造に従っているか？
  - ✅ FR-008: 既存機能（設定ファイル作成、プロファイル対応）を維持
  - ✅ FR-006: 設定ファイルのパーミッション0644を維持
  - ✅ 既存のViper/設定管理ロジックは変更せず、I/O部分のみリファクタリング

- [x] **ユーザーエクスペリエンス**: Cobraによる一貫したCLIインターフェースを提供しているか？
  - ✅ User Story 1 (P1): Cobra標準のI/Oストリーム使用を明示
  - ✅ FR-001: `cmd.OutOrStdout()`, `cmd.ErrOrStderr()`の使用を要求
  - ✅ FR-009: echoサブコマンドと同じパターンに統一

- [x] **パフォーマンス要件**: CLI起動時間100ms以下、設定読み込み10ms以下を達成できるか？
  - ✅ リファクタリングのため、既存パフォーマンスは維持される見込み
  - ✅ `internal/stdio`削除により、間接層が減少し、パフォーマンス改善の可能性あり

**Gate Status**: ✅ **PASS** - すべての憲章原則に準拠。Phase 0研究に進行可能。

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
go-cli/
├── cmd/                           # CLI command entry points (Cobra-dependent)
│   ├── root.go                   # Root command, config initialization
│   ├── configure.go              # [REFACTORED] configure subcommand
│   ├── configure_test.go         # [UPDATED] configure command tests
│   ├── echo.go                   # [REFERENCE] echo subcommand pattern
│   └── echo_test.go              # [REFERENCE] echo command tests
│
├── internal/                      # Internal packages (framework-agnostic)
│   ├── cmd/                      # Command business logic
│   │   ├── configure/            # [NEW] configure logic package
│   │   │   ├── configure.go      # [NEW] Configure(target, opts) function
│   │   │   └── configure_test.go # [NEW] unit tests for Configure()
│   │   └── echo/                 # [REFERENCE] echo logic package
│   │       ├── echo.go           # [REFERENCE] Echo(text, opts) function
│   │       └── echo_test.go      # [REFERENCE] unit tests for Echo()
│   │
│   ├── editor/                   # Editor detection utilities
│   │   ├── editor.go             # [UNCHANGED] GetEditor() function
│   │   └── editor_test.go        # [UNCHANGED] editor tests
│   │
│   ├── proc/                     # Process execution utilities
│   │   ├── process.go            # [UNCHANGED] ExecCommand(), Run()
│   │   └── process_test.go       # [UNCHANGED] process tests
│   │
│   └── stdio/                    # [DEPRECATED] Standard I/O utilities
│       ├── stdio.go              # [TO BE DELETED] GetStreams()
│       └── stdio_test.go         # [TO BE DELETED] stdio tests
│
├── main.go                        # Entry point
├── go.mod                         # Go module definition
├── Makefile                       # Build targets (all, build, test, lint, fmt)
└── .github/workflows/ci.yaml      # CI pipeline
```

**Structure Decision**: Single CLI project (Option 1) - This is a refactoring of an existing Cobra-based CLI. The chosen structure follows established Go conventions:

- **`cmd/`**: Contains Cobra command definitions that handle flag parsing, help text, and orchestration. Each command file is lightweight and delegates to `internal/`.
- **`internal/cmd/`**: Contains framework-agnostic business logic organized by feature (configure, echo). Each package exports a primary function (e.g., `Configure()`) and an options struct.
- **`internal/editor/`, `internal/proc/`**: Shared utilities used by multiple commands.
- **`internal/stdio/`**: Deprecated package to be removed as part of this refactoring.

This structure ensures:
1. **Testability**: Business logic in `internal/` can be tested without Cobra dependencies.
2. **Reusability**: Logic can be used by other commands or future non-CLI interfaces.
3. **Constitution compliance**: Clear separation between CLI layer (`cmd/`) and domain logic (`internal/`).

## Complexity Tracking

**No violations detected.** All constitution principles are satisfied:

- ✅ TDD mandatory: Test-first approach documented in quickstart.md
- ✅ Package separation: cmd/ and internal/ clearly separated
- ✅ Code quality: make lint required for success
- ✅ Config management: Viper structure maintained
- ✅ UX consistency: Cobra streams used throughout
- ✅ Performance: Refactoring maintains existing performance

## Post-Design Constitution Re-evaluation

*Re-evaluating after Phase 1 design completion (data-model.md, contracts/, quickstart.md)*

- [x] **TDD必須**: デザインがテストファーストアプローチを反映しているか？
  - ✅ `data-model.md`: Testing Strategy セクションで80%カバレッジ目標とモックパターンを定義
  - ✅ `contracts/configure-function.md`: Testing Contract セクションで11の必須テストケースを明記
  - ✅ `quickstart.md`: Phase 1.2-1.4でテストファースト実装手順を詳述（Red → Green → Refactor）
  - ✅ すべての新規コード（ConfigureOptions, Configure()）にテストケースが対応

- [x] **パッケージ責務分離**: デザインがCLIとビジネスロジックを分離しているか？
  - ✅ `data-model.md`: Data Flow セクションで `cmd/configure.go` → `internal/cmd/configure.Configure()` の分離を明記
  - ✅ `contracts/configure-function.md`: Configure()の署名を `(target string, opts ConfigureOptions) error` として定義
  - ✅ `quickstart.md`: Phase 2で `cmd/configure.go` のリファクタリング手順を記載（フラグ → オプション → 内部関数呼び出し）
  - ✅ ConfigureOptionsにI/Oストリームを注入することで、Cobra依存を排除

- [x] **コード品質基準**: デザインが品質チェックを組み込んでいるか？
  - ✅ `quickstart.md`: Phase 4でフォーマット(`make fmt`)とリント(`make lint`)を必須化
  - ✅ Verification Checklistに `make lint` パスを明記（SC-008）
  - ✅ すべてのコードサンプルがgofmtに準拠

- [x] **設定管理の一貫性**: デザインが既存の設定構造を維持しているか？
  - ✅ `data-model.md`: ConfigureOptions.Dataフィールドで既存のViperマップ構造を保持
  - ✅ `contracts/configure-function.md`: Input Contract で `~/.config/mycli/<profile>.yaml` パスを維持
  - ✅ `quickstart.md`: BuildEffectiveConfig()の呼び出しで既存設定ロジックを再利用

- [x] **ユーザーエクスペリエンス**: デザインがCobraのベストプラクティスに従っているか？
  - ✅ `data-model.md`: ConfigureOptions.Output/ErrOutputで `cmd.OutOrStdout()` / `cmd.ErrOrStderr()` を使用
  - ✅ `quickstart.md`: Phase 2.2で `cmd.OutOrStdout()` / `cmd.ErrOrStderr()` の実装例を明記
  - ✅ SC-007: echoサブコマンドと構造が一致（±3行）

- [x] **パフォーマンス要件**: デザインがパフォーマンス目標を維持しているか？
  - ✅ `contracts/configure-function.md`: Behavior Specification で os.MkdirAll / os.WriteFile の直接呼び出し（最小オーバーヘッド）
  - ✅ `internal/stdio` 削除により間接層を排除
  - ✅ 既存のファイル操作ロジックを維持（パフォーマンス劣化なし）

**Post-Design Gate Status**: ✅ **PASS** - Phase 1デザインはすべての憲章原則を満たしている。Phase 2（タスク生成）に進行可能。
