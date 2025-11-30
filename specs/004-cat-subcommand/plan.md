# Implementation Plan: Cat サブコマンド実装

**Branch**: `004-cat-subcommand` | **Date**: 2025-11-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/004-cat-subcommand/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

UNIX標準のcatコマンドクローンとして、ファイル内容の表示、複数ファイルの連結、標準入力からの読み込み、行番号表示(-n, -b)、非表示文字の可視化(-E, -T, -v, -A)の機能を実装する。32KBバッファによるストリーム処理で1GB以上のファイルをメモリ効率的に処理し、すべてのオプションの組み合わせをサポートする。Cobraフレームワークでサブコマンドとして実装し、TDDアプローチでBATSによるintegration testも含む。

## Technical Context

**Language/Version**: Go 1.25.4  
**Primary Dependencies**: Cobra v1.10.1+（CLI フレームワーク）、Viper v1.21.0+（設定管理）  
**Storage**: ファイルシステム（読み取り専用）、標準入力/出力  
**Testing**: Go標準テスト（go test）、BATS v1.x（integration test）  
**Target Platform**: macOS、Linux（クロスプラットフォーム対応）  
**Project Type**: CLI ツール（単一プロジェクト）  
**Performance Goals**: 1MB以下のファイルを100ms以内に処理、1GBファイルをメモリ100MB以下で処理  
**Constraints**: 32KBバッファサイズ、行番号6桁固定フォーマット、ASCII制御文字変換（0-31, 127）  
**Scale/Scope**: 8つのユーザーストーリー、18の機能要件、6つのCLIオプション

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **TDD必須**: ✅ すべての実装に対してテストを先に書く計画がある（SC-007で確認）。cmd/cat_test.go、internal/cmd/cat/processor_test.go等を先行実装
- [x] **パッケージ責務分離**: ✅ `cmd/cat.go`（Cobraコマンドラッパーのみ）と`internal/cmd/cat/`（ファイル処理、行番号付加、制御文字変換等の内部ロジック）が明確に分離されている
- [x] **コード品質基準**: ✅ `make all`により`gofmt -s`と`golangci-lint run --enable=govet`を実行予定。既存コードベースと同じ品質基準を適用
- [x] **設定管理の一貫性**: ✅ catコマンドは設定ファイルに依存しない（純粋なファイル処理コマンド）。Viperの設定読み込みは不要
- [x] **ユーザーエクスペリエンス**: ✅ Cobraで`mycli cat`サブコマンドとして実装。`--help`で簡潔な説明と使用例を提供（FR-016）
- [x] **パフォーマンス要件**: ✅ 1MBファイルを100ms以内で処理（SC-001）。設定読み込みは不要のため該当せず。32KBバッファによるストリーム処理で1GBファイルをメモリ100MB以下で処理（SC-002）

**評価**: すべての憲章原則に準拠。違反なし。

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
├── main.go                           # エントリーポイント（変更なし）
├── cmd/
│   ├── root.go                       # ルートコマンド（変更なし）
│   ├── cat.go                        # [新規] Cobraコマンド定義
│   ├── cat_test.go                   # [新規] catコマンドのCLI統合テスト
│   └── (既存コマンド...)
├── internal/
│   └── cmd/
│       └── cat/                      # [新規] catコマンド内部実装
│           ├── processor.go          # ファイル処理・ストリーム処理
│           ├── processor_test.go     # プロセッサーの単体テスト
│           ├── formatter.go          # 行番号・制御文字フォーマット
│           ├── formatter_test.go     # フォーマッターの単体テスト
│           └── options.go            # オプション構造体とバリデーション
├── integration_test/                 # BATSテストディレクトリ（既存）
│   ├── cat.bats                      # [新規] catコマンドのBATSテスト
│   └── helpers/                      # テストヘルパー（既存）
│       ├── common.bash
│       └── assertions.bash
├── bin/
│   └── mycli                         # ビルド成果物
├── Makefile                          # ビルド定義（変更なし）
└── go.mod                            # 依存関係（変更なし）
```

**Structure Decision**: Go標準の単一プロジェクト構造を採用。既存のプロジェクト構造に準拠し、catサブコマンドを`cmd/cat.go`（CLI層）と`internal/cmd/cat/`（ロジック層）に分離して実装。BATSテストは既存の`integration_test/`ディレクトリに追加。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

**評価結果**: 憲章違反なし。このセクションは該当なし。

## Phase 0: Research & Outline ✅ COMPLETED

**Status**: Complete  
**Date**: 2025-11-30

### Research Document

**File**: `specs/004-cat-subcommand/research.md`

**Summary**:
- すべての"NEEDS CLARIFICATION"項目が明確化セッションで解決済み
- 技術選択の根拠を文書化（`bufio.Scanner`、32KBバッファ、制御文字ルックアップテーブル）
- ベストプラクティスの適用方針を明確化（Go標準慣習、憲章準拠、パフォーマンス最適化）
- 3層テストアプローチを定義（単体テスト、統合テスト、BATSテスト）

**Key Decisions**:
1. **ファイル処理**: `bufio.Scanner`による行単位ストリーム処理（32KBバッファ）
2. **行番号フォーマット**: `fmt.Sprintf("%6d  ", lineNum % 1000000)`で6桁固定
3. **制御文字変換**: ルックアップテーブルによる高速変換（ASCII 0-31 + 127）
4. **エラーハンドリング**: 継続処理と終了コード管理（1つでもエラーがあれば終了コード1）
5. **テスト戦略**: 単体テスト（100%カバレッジ目標）、統合テスト（全ユーザーストーリー）、BATS（エッジケース）

## Phase 1: Design & Contracts ✅ COMPLETED

**Status**: Complete  
**Date**: 2025-11-30

### Data Model

**File**: `specs/004-cat-subcommand/data-model.md`

**Core Entities**:
1. **Options**: コマンドラインフラグから抽出されたオプション設定を保持
   - `NumberAll`, `NumberNonBlank`, `ShowEnds`, `ShowTabs`, `ShowNonPrinting`
   - オプション競合解決ロジック（`-n`と`-b`の後方優先、`-A`フラグの展開）

2. **Processor**: ファイルの読み込みとストリーム処理を担当
   - `ProcessFile(filename, opts, output)`: ファイルを処理
   - `ProcessStdin(opts, output)`: 標準入力を処理
   - ステートレス、各メソッド呼び出しは独立

3. **Formatter**: 行のフォーマット処理を担当
   - `FormatLine(line, lineNum, isEmpty, opts)`: 行をフォーマット
   - 制御文字変換マップを事前構築

**State Transitions**:
- ファイル処理フロー: 開始 → ファイルオープン → 行読み込み → 行フォーマット → 出力 → 次の行
- 行番号カウンタ: 単調増加（巻き戻しなし）、999,999超過時は最下位6桁のみ表示
- オプション解決: コマンドライン引数 → Cobraフラグ解析 → `-A`展開 → Options構造体

### Contracts

**Files**: `specs/004-cat-subcommand/contracts/*.md`

1. **processor.md**: Processorインターフェースの契約定義
   - 前提条件、事後条件、不変条件を明確化
   - エラーハンドリングの詳細（`os.ErrNotExist`, `os.ErrPermission`など）
   - パフォーマンス特性（O(n)時間、O(1)空間、100ms/1MBファイル）

2. **formatter.md**: Formatterインターフェースの契約定義
   - フォーマットルール（行番号6桁、制御文字変換、タブ変換、行末マーカー）
   - エッジケース（999,999超過、空行、制御文字とタブの混在）
   - パフォーマンス特性（O(n)時間、行あたり10μs以下）

3. **options.md**: Options構造体の契約定義
   - ファクトリ関数`NewOptions(cmd)`による生成
   - オプション競合解決ロジック（`-n`と`-b`の後方優先、`-A`の展開）
   - バリデーションルール（VR-001からVR-011）

### Quickstart

**File**: `specs/004-cat-subcommand/quickstart.md`

**Content**:
- 必須ツールと環境設定（Go 1.25.4、golangci-lint 2.6.2、bats 1.x）
- プロジェクト構造の説明
- TDDワークフロー（Red-Green-Refactor）の具体的な手順
  - Phase 1: Formatter実装
  - Phase 2: Processor実装
  - Phase 3: Options実装
  - Phase 4: Cobraコマンド統合
  - Phase 5: BATS統合テスト
- テスト実行方法とトラブルシューティング

### Agent Context Update

**Status**: ✅ Completed

Copilot agent context file (`.github/agents/copilot-instructions.md`) updated with:
- Language: Go 1.25.4
- Framework: Cobra v1.10.1+, Viper v1.21.0+
- Database: ファイルシステム（読み取り専用）、標準入力/出力

### Constitution Re-Check

**Status**: ✅ Passed

Phase 1設計後の憲章再評価:
- [x] **TDD必須**: quickstart.mdでRed-Green-Refactorの具体的手順を定義
- [x] **パッケージ責務分離**: data-model.mdで`cmd/cat.go`と`internal/cmd/cat/`の明確な分離を確認
- [x] **コード品質基準**: quickstart.mdで`make all`による品質チェック手順を明記
- [x] **設定管理の一貫性**: catコマンドは設定ファイル不要（変更なし）
- [x] **ユーザーエクスペリエンス**: contracts/options.mdでCobraフラグの標準的な使用方法を定義
- [x] **パフォーマンス要件**: data-model.mdで32KBバッファとストリーム処理による性能保証を明記

**評価**: すべての憲章原則に準拠。違反なし。

## Phase 2: Task Breakdown 🔄 PENDING

**Status**: Not Started  
**Next Action**: Execute `/speckit.tasks` command to generate `tasks.md`

**Expected Output**:
- `specs/004-cat-subcommand/tasks.md`: 優先順位付きタスクリスト
- タスクごとにTDDサイクルの具体的な手順（Red-Green-Refactor）
- 依存関係とマイルストーンの定義

## Implementation Phase 🔄 PENDING

**Status**: Not Started  
**Prerequisites**: Phase 2（Task Breakdown）完了後

**Approach**: TDD（Red-Green-Refactor）サイクルに従い、tasks.mdの順序で実装

**Milestones**:
1. Formatter実装（最も独立したコンポーネント）
2. Processor実装（Formatterに依存）
3. Options実装（フラグ解析）
4. Cobraコマンド統合（すべてのコンポーネントを統合）
5. BATS統合テスト（エンドツーエンド検証）
