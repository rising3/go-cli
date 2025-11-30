# Implementation Plan: Echo サブコマンド実装

**Branch**: `001-echo-subcommand` | **Date**: 2025-11-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-echo-subcommand/spec.md`

## Summary

UNIX標準の`echo`コマンドのクローンをCobraフレームワークを使用してCLIサブコマンドとして実装する。基本的なテキスト出力、改行抑制オプション（`-n`）、エスケープシーケンス解釈オプション（`-e`）、デバッグ用の`--verbose`フラグを提供し、完全なUNIX互換性を実現する。TDD（テスト駆動開発）アプローチを採用し、`cmd/`パッケージにCLIインターフェース、`internal/echo/`パッケージに内部ロジックを配置することで、関心の分離を徹底する。

**技術的アプローチ**: エスケープシーケンス処理は`internal/echo/processor.go`に独立した関数として実装し、テスト容易性を確保。Cobraの標準エラーハンドリング機構を活用し、無効なオプション時の自動ヘルプ表示を実現。UTF-8処理はGoの標準文字列操作に依存することでシンプルさを維持。

## Technical Context

**Language/Version**: Go 1.25.4（`go.mod`で指定、憲章で定義済み）  
**Primary Dependencies**: 
  - Cobra v1.10.1+（CLIフレームワーク、憲章必須要件）
  - Go標準ライブラリのみ（`fmt`, `strings`, `os`など）
  
**Storage**: N/A（ステートレスなコマンド、設定ファイル不要）  
**Testing**: Go標準テストフレームワーク（`testing`パッケージ、`go test`コマンド）  
**Target Platform**: クロスプラットフォーム（Linux、macOS、Windows対応、Go標準互換）  
**Project Type**: Single project（既存のgo-cliプロジェクトにサブコマンド追加）  
**Performance Goals**: 
  - CLI起動時間: 100ms以内（憲章要件）
  - ヘルプ表示: 50ms以内（SC-003）
  - 10,000引数処理: メモリ100MB以下（SC-004）
  
**Constraints**: 
  - UTF-8エンコーディングのみサポート（FR-014、シンプルさ優先）
  - 外部依存ライブラリ追加禁止（Go標準+Cobraのみ使用）
  - UNIX標準echoコマンドとの出力互換性必須（SC-002）
  
**Scale/Scope**: 
  - 単一サブコマンド実装（`mycli echo`）
  - 3つの主要オプション（`-n`, `-e`, `--verbose`）
  - 9種類のエスケープシーケンス対応
  - 14個の機能要件（FR-001～FR-014）
  - 8個の成功基準（SC-001～SC-008）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **TDD必須**: すべての実装に対してテストを先に書く計画があるか？
  - ✅ `cmd/echo_test.go`と`internal/echo/processor_test.go`を実装前に作成
  - ✅ Red-Green-Refactorサイクルを各User Storyで適用
  
- [x] **パッケージ責務分離**: `cmd/`（CLI）と`internal/`（内部ロジック）が明確に分離されているか？
  - ✅ `cmd/echo.go`: Cobraコマンド定義とフラグ処理のみ
  - ✅ `internal/echo/processor.go`: エスケープシーケンス処理ロジック
  
- [x] **コード品質基準**: `gofmt`と`govet`による検証を通過する見込みか？
  - ✅ 実装完了後に`make all`で検証（test → fmt → lint → build）
  
- [x] **設定管理の一貫性**: Viperを使用し、`~/.config/mycli/`配下の設定ファイル構造に従っているか？
  - ✅ N/A - echoコマンドは設定ファイル不要（ステートレス）
  
- [x] **ユーザーエクスペリエンス**: Cobraによる一貫したCLIインターフェースを提供しているか？
  - ✅ Cobraの標準パターンに従ったサブコマンド実装
  - ✅ FR-008: 簡潔なヘルプメッセージ + 2-3個の使用例
  - ✅ FR-012: 無効オプション時の自動ヘルプ表示
  
- [x] **パフォーマンス要件**: CLI起動時間100ms以下、設定読み込み10ms以下を達成できるか？
  - ✅ SC-001: 100ms以内の実行完了目標
  - ✅ SC-003: 50ms以内のヘルプ表示目標
  - ✅ N/A - 設定読み込み不要

**Phase 0 Research Gate**: ✅ PASSED - すべての憲章要件に準拠

---

### Post-Phase 1 Re-evaluation

*Constitution Check completed after Phase 1 design (research.md, data-model.md, contracts/, quickstart.md)*

- [x] **TDD必須**: 設計段階でテスト戦略が明確化されているか？
  - ✅ `data-model.md`セクション7に詳細なテスト戦略を記載
  - ✅ `ProcessEscapes()`のユニットテスト仕様定義済み
  - ✅ `GenerateOutput()`のユニットテスト仕様定義済み
  - ✅ `cmd/echo_test.go`の統合テスト仕様定義済み
  - ✅ `contracts/echo-command.md`にテストコントラクト明記
  
- [x] **パッケージ責務分離**: 設計が憲章の分離原則に準拠しているか？
  - ✅ `data-model.md`で`EchoOptions`, `ProcessEscapes()`, `GenerateOutput()`を`internal/echo`に配置
  - ✅ `contracts/echo-command.md`でCLIインターフェースを定義（`cmd/echo.go`に実装予定）
  - ✅ Data Flow Diagramで責務分離を可視化
  
- [x] **コード品質基準**: 設計が品質基準を考慮しているか？
  - ✅ `research.md`でGoの標準パターン（`strings.Builder`）を採用
  - ✅ `data-model.md`セクション5でメモリ効率を考慮した設計
  - ✅ `quickstart.md`に品質ゲート（`make all`）の実行手順を記載
  
- [x] **設定管理の一貫性**: 設定不要な設計が妥当か？
  - ✅ N/A - echoコマンドはステートレスで、設定ファイル不要
  - ✅ `data-model.md`で「永続化されるデータモデルは存在しない」と明記
  
- [x] **ユーザーエクスペリエンス**: 設計がユーザビリティを考慮しているか？
  - ✅ `contracts/echo-command.md`で明確なヘルプメッセージ形式を定義
  - ✅ `quickstart.md`で実用的なユースケース（7種類）を提供
  - ✅ エラーハンドリング（自動ヘルプ表示）を`research.md`セクション3で設計
  
- [x] **パフォーマンス要件**: 設計がパフォーマンス目標を達成可能か？
  - ✅ `data-model.md`セクション5で10,000引数・100MB以下のメモリ戦略を記載
  - ✅ `research.md`セクション2で効率的なエスケープ処理（`strings.Builder`）を選択
  - ✅ `contracts/echo-command.md`でパフォーマンス要件（100ms起動、50msヘルプ、100MB）を明記

**Phase 1 Design Gate**: ✅ PASSED - 設計は憲章のすべての原則に準拠し、実装可能な状態です

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

**Structure Decision**: Single project Go CLI application (Cobra/Viper architecture)

#### 新規作成ファイル

```text
cmd/
├── echo.go                       # NEW: Cobraコマンド定義、フラグ処理、internal/echoへの委譲
└── echo_test.go                  # NEW: コマンドレベルの統合テスト（TDD: Red phase first）

internal/
└── echo/                         # NEW: echoパッケージディレクトリ
    ├── processor.go              # NEW: エスケープシーケンス処理ロジック（framework-independent）
    └── processor_test.go         # NEW: processor.goのユニットテスト（TDD: Red phase first）
```

#### 既存ファイル（変更なし）

```text
main.go                           # UNCHANGED: cmd.Execute()呼び出し
cmd/
├── root.go                       # UNCHANGED: rootコマンド定義
├── configure.go                  # UNCHANGED: configureサブコマンド
└── viperutils.go                 # UNCHANGED: Viper設定ユーティリティ

internal/
├── cmd/configure.go              # UNCHANGED: configure実装
├── editor/editor.go              # UNCHANGED: エディタ検出
├── proc/process.go               # UNCHANGED: プロセス実行
└── stdio/stdio.go                # UNCHANGED: 標準I/O処理
```

#### テストとのマッピング

| Implementation File              | Test File                         | Test Scope                                  |
|----------------------------------|-----------------------------------|---------------------------------------------|
| `cmd/echo.go`                    | `cmd/echo_test.go`                | Cobraコマンド統合テスト（フラグパース、出力検証）|
| `internal/echo/processor.go`     | `internal/echo/processor_test.go` | エスケープシーケンス処理のユニットテスト       |

#### TDD実装順序

1. **Red Phase**: `cmd/echo_test.go` + `internal/echo/processor_test.go` 作成
2. **Green Phase**: `cmd/echo.go` + `internal/echo/processor.go` 実装
3. **Refactor Phase**: リファクタリング実施
4. **Validation**: `make all`実行（test → fmt → lint → build）

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

**Status**: ✅ No violations detected

すべての憲章要件に準拠しているため、本セクションは空です。

- TDD必須: ✅ 準拠（テスト先行アプローチ確定）
- パッケージ責務分離: ✅ 準拠（cmd/とinternal/echo/の明確な分離）
- コード品質基準: ✅ 準拠（gofmt/govetによる検証予定）
- 設定管理の一貫性: ✅ N/A（ステートレスコマンド）
- ユーザーエクスペリエンス: ✅ 準拠（Cobraの標準パターン使用）
- パフォーマンス要件: ✅ 準拠（100ms起動、50msヘルプ目標設定）
