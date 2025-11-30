# Data Model: Echo Subcommand

**Feature**: Echo サブコマンド実装  
**Date**: 2025-11-30  
**Status**: Draft

## Overview

echoサブコマンドはステートレスなCLIコマンドであり、永続化されるデータモデルは存在しません。このドキュメントでは、実行時に扱うデータ構造とそのライフサイクルを定義します。

---

## 1. Command Options (Runtime Configuration)

### `EchoOptions`

**Purpose**: echoコマンドの実行時オプションを保持する構造体

**Lifecycle**: コマンド実行開始時に作成され、実行完了後に破棄（メモリのみ）

**Package**: `internal/echo`

```go
package echo

// EchoOptions represents runtime configuration for the echo command.
type EchoOptions struct {
    // SuppressNewline suppresses the trailing newline character when true (-n flag).
    SuppressNewline bool
    
    // InterpretEscapes enables backslash escape sequence interpretation when true (-e flag).
    InterpretEscapes bool
    
    // Verbose enables debug logging to stderr when true (--verbose flag).
    Verbose bool
    
    // Args contains the arguments to be echoed (after flag parsing).
    Args []string
}
```

**Validation Rules**:
- `SuppressNewline`, `InterpretEscapes`, `Verbose`: boolean型、デフォルト`false`
- `Args`: 空のスライスも有効（空行を出力するケース）

**Relationships**:
- Cobraの`cmd.Flags()`から値を取得し、この構造体に格納
- `Processor`関数に渡され、出力生成に使用

---

## 2. Escape Sequence Processing

### `ProcessEscapes` Function Signature

**Purpose**: エスケープシーケンスを解釈し、実際の制御文字に変換

**Package**: `internal/echo`

```go
// ProcessEscapes interprets backslash escape sequences in the input string.
// Returns the processed output string and a flag indicating if newline suppression
// was triggered by \c escape sequence.
func ProcessEscapes(input string) (output string, suppressNewline bool)
```

**Input**: 
- `input string`: エスケープシーケンスを含む可能性のある入力文字列

**Output**:
- `output string`: エスケープシーケンスが解釈された出力文字列
- `suppressNewline bool`: `\c`が検出された場合は`true`（それ以降の出力を抑制）

**Supported Escape Sequences** (FR-004):

| Escape Sequence | Description | Output Character |
|-----------------|-------------|------------------|
| `\n` | Newline | `\n` (0x0A) |
| `\t` | Horizontal tab | `\t` (0x09) |
| `\\` | Backslash | `\` (0x5C) |
| `\"` | Double quote | `"` (0x22) |
| `\a` | Alert (bell) | `\a` (0x07) |
| `\b` | Backspace | `\b` (0x08) |
| `\c` | Suppress further output | (returns immediately) |
| `\r` | Carriage return | `\r` (0x0D) |
| `\v` | Vertical tab | `\v` (0x0B) |

**Edge Cases**:
- 無効なエスケープシーケンス（例: `\z`）→ リテラル文字列として扱う（`\z`をそのまま出力）
- 末尾のバックスラッシュ（例: `Hello\`）→ リテラル文字列として扱う

---

## 3. Output Generation

### `GenerateOutput` Function Signature

**Purpose**: `EchoOptions`に基づいて最終的な出力文字列を生成

**Package**: `internal/echo`

```go
// GenerateOutput generates the final output string based on the provided options.
// Returns the output string ready to be written to stdout.
func GenerateOutput(opts EchoOptions) string
```

**Processing Logic**:
1. `opts.Args`をスペース区切りで結合（`strings.Join(opts.Args, " ")`）
2. `opts.InterpretEscapes`が`true`の場合、`ProcessEscapes()`を適用
3. `opts.SuppressNewline`が`false`かつ`ProcessEscapes()`で`\c`が検出されていない場合、末尾に`\n`を追加
4. 最終的な出力文字列をreturn

**State Transitions**:
```
Input Args → Join with spaces → (If -e) Process Escapes → (If not -n and not \c) Append \n → Output
```

---

## 4. Error Handling

### Error Types

echoコマンドは以下のエラーを処理します:

#### 4.1. Invalid Flag Error (Cobra managed)

**Trigger**: 無効なフラグが指定された場合（例: `mycli echo -x`）  
**Handling**: Cobraが自動的にエラーメッセージとヘルプを表示（FR-012）  
**Exit Code**: 1 (FR-011)

#### 4.2. No Errors Expected (Normal Operation)

**Rationale**: echoコマンドは、どのような引数でも正常に処理可能  
**Validation**: 引数の検証は不要（空文字列、大量の引数、特殊文字すべて有効）

---

## 5. Memory Management

### Performance Constraints (SC-004)

**Requirement**: 10,000個の引数を渡しても、メモリ使用量が100MB以下

**Strategy**:
- `strings.Builder`を使用した効率的な文字列構築（再割り当てを最小化）
- 引数は`[]string`で保持し、不要なコピーを避ける
- エスケープシーケンス処理はストリーミング方式（1文字ずつ処理）

**Measurement**:
- テストで10,000個の引数を生成し、`runtime.MemStats`でメモリ使用量を検証

---

## 6. Data Flow Diagram

```
┌─────────────────────┐
│ User Input          │
│ mycli echo -n -e    │
│ "Hello\nWorld"      │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Cobra Parsing       │
│ (cmd/echo.go)       │
│ - Flags: -n, -e     │
│ - Args: ["Hello\n...│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ EchoOptions         │
│ SuppressNewline=true│
│ InterpretEscapes=true│
│ Args=["Hello\nWorld"]│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ GenerateOutput()    │
│ (internal/echo)     │
│ - Join with spaces  │
│ - ProcessEscapes()  │
│ - Append \n?        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Output to stdout    │
│ "Hello              │
│ World"              │
└─────────────────────┘
```

---

## 7. Testing Strategy

### Unit Tests

**Package**: `internal/echo`

#### 7.1. `ProcessEscapes()` Tests

- 各エスケープシーケンスの正しい解釈（`\n`, `\t`, `\\`, etc.）
- `\c`の出力抑制動作
- 無効なエスケープシーケンスのリテラル扱い
- 複数のエスケープシーケンスが混在するケース

#### 7.2. `GenerateOutput()` Tests

- 基本的な引数結合（スペース区切り）
- `-n`フラグの改行抑制
- `-e`フラグのエスケープシーケンス解釈
- `-n`と`-e`の組み合わせ
- 引数なしのケース（空行出力）

### Integration Tests

**Package**: `cmd`

- Cobraコマンドの統合テスト（`cmd/echo_test.go`）
- 実際のフラグパースと出力検証
- エラーハンドリング（無効なフラグ）
- `--verbose`フラグのデバッグ出力検証

---

## 8. Constitution Compliance

### TDD必須
- ✅ `internal/echo/processor_test.go`で`ProcessEscapes()`のテストを先に作成
- ✅ `internal/echo/echo_test.go`で`GenerateOutput()`のテストを先に作成
- ✅ `cmd/echo_test.go`でコマンドレベルの統合テストを先に作成

### パッケージ責務分離
- ✅ `internal/echo`: Cobra/Viperに依存しないピュアなロジック
- ✅ `cmd/echo.go`: Cobraフレームワークとの統合のみ

### パフォーマンス要件
- ✅ SC-001: 100ms以内の実行完了（標準出力はほぼ瞬時）
- ✅ SC-004: 10,000引数でメモリ100MB以下（`strings.Builder`で効率化）

---

## Summary

echoサブコマンドは、永続化されるデータモデルを持たないステートレスなコマンドです。実行時に`EchoOptions`構造体で設定を保持し、`internal/echo`パッケージのピュアな関数（`ProcessEscapes`, `GenerateOutput`）でビジネスロジックを実装します。TDD原則に従い、すべてのロジックは単体テストでカバーされます。
