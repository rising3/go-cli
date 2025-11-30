# Implementation Plan: Configure設定構造のリファクタリング

**Branch**: `005-configure-refactor` | **Date**: 2025-11-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/005-configure-refactor/spec.md`

**Note**: This plan implements the expansion of the Config struct to support nested configuration structure.

## Summary

Extend the mycli configuration file structure to support nested settings. The primary goals are: (1) Add nested Config struct fields (`Common`, `Hoge`, `Hoge.Foo`), (2) Update `BuildEffectiveConfig()` to populate default values for new fields, (3) Maintain backward compatibility with existing `client-id` and `client-secret` fields, (4) Ensure all mapstructure tags use kebab-case for Viper integration, and (5) Verify tests pass with the new configuration structure.

## Technical Context

**Language/Version**: Go 1.25.4  
**Primary Dependencies**: Cobra v1.10.1, Viper v1.21.0, gopkg.in/yaml.v3 v3.0.1  
**Storage**: Local filesystem (~/.config/mycli/*.yaml configuration files)  
**Testing**: Go standard testing framework (`go test`), table-driven tests  
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows)  
**Project Type**: Single CLI project (Cobra-based command structure)  
**Performance Goals**: CLI startup <100ms, config file read <10ms, help display <50ms  
**Constraints**: File permissions 0644 for config files, backward compatibility with existing config files  
**Scale/Scope**: Config struct expansion affecting 2 files (cmd/root.go, cmd/viperutils.go) + tests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **TDD必須**: すべての実装に対してテストを先に書く計画があるか？
  - ✅ FR-008: 既存テストの継続パスと必要に応じた更新を要求
  - ✅ SC-004: リファクタリング前のすべてのテストケースが`make test`で100%パスすることを成功基準として定義
  - ✅ User Story 2 (P2): Config構造体のViperマッピングテストを明示

- [x] **パッケージ責務分離**: `cmd/`（CLI）と`internal/`（内部ロジック）が明確に分離されているか？
  - ✅ 変更対象は`cmd/root.go`と`cmd/viperutils.go`のみ（データ構造定義）
  - ✅ FR-009: `internal/cmd/configure/configure.go`の既存ロジックは変更不要
  - ✅ 既存のパッケージ分離構造を維持

- [x] **コード品質基準**: `gofmt`と`govet`による検証を通過する見込みか？
  - ✅ SC-005: `make lint`がゼロ件の警告・エラーで終了することを成功基準として定義
  - ✅ 構造体とmapstructureタグの追加のみで、既存コード品質を維持

- [x] **設定管理の一貫性**: Viperを使用し、`~/.config/mycli/`配下の設定ファイル構造に従っているか？
  - ✅ FR-001, FR-002: Config構造体にmapstructureタグを使用してViper統合を維持
  - ✅ FR-005, FR-006: 既存のViper読み込みロジックで新構造を自動マッピング
  - ✅ 設定ファイルパス（`~/.config/mycli/`）は変更なし

- [x] **ユーザーエクスペリエンス**: Cobraによる一貫したCLIインターフェースを提供しているか？
  - ✅ FR-007: 既存の全フラグ（`--force`, `--edit`, `--no-wait`, `--profile`）の動作を維持
  - ✅ User Story 4 (P4): 既存機能との互換性維持を明示
  - ✅ CLI動作に変更なし（内部データ構造の拡張のみ）

- [x] **パフォーマンス要件**: CLI起動時間100ms以下、設定読み込み10ms以下を達成できるか？
  - ✅ 構造体フィールドの追加のみで、パフォーマンスへの影響は最小限
  - ✅ Viperのアンマーシャル処理は変更なし（既存の効率的な処理を維持）
  - ✅ YAMLファイルサイズの増加は微小（7フィールド程度）

**Gate Status**: ✅ **PASS** - すべての憲章原則に準拠。Phase 0研究に進行可能。

## Project Structure

### Documentation (this feature)

```text
specs/005-configure-refactor/
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
│   ├── root.go                   # [TO UPDATE] Config struct definition
│   ├── viperutils.go             # [TO UPDATE] BuildEffectiveConfig() function
│   ├── configure.go              # [UNCHANGED] configure subcommand
│   ├── configure_test.go         # [POTENTIALLY UPDATE] configure command tests
│   └── viperutils_test.go        # [POTENTIALLY UPDATE] viperutils tests
│
├── internal/                      # Internal packages (framework-agnostic)
│   └── cmd/                      # Command business logic
│       └── configure/            # configure logic package
│           ├── configure.go      # [UNCHANGED] Configure(target, opts) function
│           └── configure_test.go # [UNCHANGED] unit tests for Configure()
│
├── main.go                        # Entry point
├── go.mod                         # Go module definition
├── Makefile                       # Build targets (all, build, test, lint, fmt)
└── .github/workflows/ci.yaml      # CI pipeline
```

**Structure Decision**: Single CLI project (Option 1) - This is an expansion of the existing Config struct to support nested configuration. The chosen approach:

- **`cmd/root.go`**: Define new nested structs (`CommonConfig`, `HogeConfig`, `FooConfig`) and add fields to existing `Config` struct
- **`cmd/viperutils.go`**: Update `BuildEffectiveConfig()` to return map with new nested structure
- **`cmd/configure.go`**: No changes required (already uses `BuildEffectiveConfig()`)
- **`internal/cmd/configure/`**: No changes required (receives config as `map[string]interface{}`)

This structure ensures:
1. **Minimal Change Surface**: Only 2 files require updates (root.go, viperutils.go)
2. **Backward Compatibility**: Existing fields (`ClientID`, `ClientSecret`) remain unchanged
3. **Constitution Compliance**: Package separation maintained (data structures in cmd/, logic in internal/)
4. **Testability**: Config unmarshaling can be tested with Viper integration tests

## Complexity Tracking

**No violations detected.** All constitution principles are satisfied:

- ✅ TDD mandatory: Test updates planned for existing test files
- ✅ Package separation: Changes limited to cmd/ (data structures only)
- ✅ Code quality: make lint required for success
- ✅ Config management: Viper structure extended with mapstructure tags
- ✅ UX consistency: CLI behavior unchanged
- ✅ Performance: Minimal impact from additional struct fields

## Post-Design Constitution Re-evaluation

*Re-evaluating after Phase 1 design completion (data-model.md, contracts/, quickstart.md)*

- [x] **TDD必須**: デザインがテストファーストアプローチを反映しているか？
  - ✅ `quickstart.md`: Phase 1.2でテストを先に書き、Phase 1.1で実装する順序を明記
  - ✅ `quickstart.md`: Phase 2.1でBuildEffectiveConfigのテストを先に書き、Phase 2.2で実装
  - ✅ `contracts/config-struct.md`: Testing Contract セクションで5つの必須テストケースを定義
  - ✅ `contracts/build-effective-config.md`: Testing Contract セクションで6つの必須テストケースを定義
  - ✅ 各コントラクトにテストヘルパー関数の例を含む

- [x] **パッケージ責務分離**: デザインがCLIとビジネスロジックを分離しているか？
  - ✅ `data-model.md`: Config構造体は`cmd/root.go`に配置（データ定義のみ）
  - ✅ `data-model.md`: BuildEffectiveConfigは`cmd/viperutils.go`に配置（純粋関数）
  - ✅ `quickstart.md`: `internal/cmd/configure/configure.go`は変更不要と明記
  - ✅ 変更対象は`cmd/`パッケージのデータ構造のみで、ロジックの分離を維持

- [x] **コード品質基準**: デザインが品質チェックを組み込んでいているか？
  - ✅ `quickstart.md`: Phase 4で`make fmt`, `make lint`, `make test`を必須化
  - ✅ `quickstart.md`: Phase 6でビルドパイプライン全体（`make all`）を実行
  - ✅ Success Criteria Verification でSC-004 (make test)とSC-005 (make lint)を確認

- [x] **設定管理の一貫性**: デザインが既存の設定構造を維持しているか？
  - ✅ `data-model.md`: mapstructureタグでViperとの統合を維持
  - ✅ `contracts/config-struct.md`: Backward Compatibility セクションで後方互換性を保証
  - ✅ `contracts/build-effective-config.md`: YAML出力形式を明確に定義
  - ✅ `quickstart.md`: 既存のViper読み込みロジックを変更しないことを確認

- [x] **ユーザーエクスペリエンス**: デザインがCobraのベストプラクティスに従っているか？
  - ✅ `quickstart.md`: Phase 5でCLI動作の手動検証手順を提供
  - ✅ CLI動作（フラグ、サブコマンドなど）に変更なし
  - ✅ 生成される設定ファイルの構造が明確で理解しやすい

- [x] **パフォーマンス要件**: デザインがパフォーマンス目標を維持しているか？
  - ✅ `contracts/config-struct.md`: Performance Characteristics で O(n) アンマーシャル時間を明記
  - ✅ `contracts/build-effective-config.md`: Performance Characteristics で O(1) 実行時間を保証
  - ✅ `data-model.md`: 構造体サイズは7フィールド程度で、メモリ影響は微小
  - ✅ YAMLファイルサイズは10行程度で、読み込み時間への影響は最小限

**Post-Design Gate Status**: ✅ **PASS** - Phase 1デザインはすべての憲章原則を満たしている。Phase 2（タスク生成）に進行可能。
