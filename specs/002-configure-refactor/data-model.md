# Data Model: Configure サブコマンドのリファクタリング

**Feature**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md)  
**Created**: 2025-11-30

## Overview

このドキュメントは、configureサブコマンドのリファクタリングで使用されるデータ構造を定義します。echoサブコマンドの`EchoOptions`パターンに倣い、`ConfigureOptions`構造体を設計します。

## Core Entities

### ConfigureOptions

**Purpose**: configureコマンドの実行に必要なすべてのオプションと依存関係を集約する構造体

**Location**: `internal/cmd/configure/configure.go`

**Fields**:

| Field | Type | Purpose | Required | Default |
|-------|------|---------|----------|---------|
| `Force` | `bool` | 既存の設定ファイルを強制上書きするかどうか | Yes | `false` |
| `Edit` | `bool` | 設定ファイル作成後にエディタで開くかどうか | Yes | `false` |
| `NoWait` | `bool` | エディタをバックグラウンドで起動するかどうか (`Edit`がtrueの時のみ有効) | Yes | `false` |
| `Data` | `map[string]interface{}` | 設定ファイルに書き込むデータ | Yes | - |
| `Format` | `string` | 出力フォーマット ("yaml" または "json") | Yes | `"yaml"` |
| `Output` | `io.Writer` | 標準出力用のストリーム（メッセージ出力に使用） | Yes | - |
| `ErrOutput` | `io.Writer` | エラー出力用のストリーム（エラーメッセージ・デバッグ情報に使用） | Yes | - |
| `EditorLookup` | `func() (string, []string, error)` | エディタを検出する関数（`internal/editor.GetEditor`を注入） | Yes | - |
| `EditorShouldWait` | `func(string, []string) bool` | エディタプロセスの待機要否を判定する関数 | Yes | - |

**Usage Pattern**:

```go
opts := configure.ConfigureOptions{
    Force:            force,
    Edit:             edit,
    NoWait:           noWait,
    Data:             BuildEffectiveConfig(),
    Format:           CliConfigType,
    Output:           cmd.OutOrStdout(),
    ErrOutput:        cmd.ErrOrStderr(),
    EditorLookup:     func() (string, []string, error) { return editor.GetEditor() },
    EditorShouldWait: func(ed string, args []string) bool { return !noWait },
}
```

**Validation**:

- `Data`はnilであってはならない（空のmapは許可）
- `Format`は "yaml", "yml", "json" のいずれか
- `Output`と`ErrOutput`はnilであってはならない
- `EditorLookup`と`EditorShouldWait`は`Edit`がtrueの場合のみ使用されるが、常に設定されるべき（テスト容易性のため）

### Function Signatures

**Configure Function**:

```go
// Configure creates or overwrites a configuration file at the specified target path.
// It marshals the opts.Data according to opts.Format and writes it to the file.
// If opts.Edit is true, it launches the configured editor after file creation.
//
// Parameters:
//   - target: Absolute path to the configuration file
//   - opts: Configuration options including data, format, and I/O streams
//
// Returns:
//   - error: Returns error if file creation fails or editor launch fails
//            (unless editor detection fails, in which case error is logged and nil is returned)
func Configure(target string, opts ConfigureOptions) error
```

**ConfigureFunc Variable** (for testing):

```go
// ConfigureFunc is a variable indirection so callers (cmd package tests)
// can replace the implementation with a stub. By default it points to Configure.
var ConfigureFunc = Configure
```

## Data Flow

```
cmd/configure.go (RunE)
  ↓
1. Extract flags: force, edit, noWait, profile
  ↓
2. Determine target path: filepath.Join(GetConfigPath(), GetConfigFile(profile))
  ↓
3. Build ConfigureOptions:
   - Data: BuildEffectiveConfig()
   - Format: CliConfigType (from root.go)
   - Output: cmd.OutOrStdout()
   - ErrOutput: cmd.ErrOrStderr()
   - EditorLookup: func() { return editor.GetEditor() }
   - EditorShouldWait: func(ed, args) { return !noWait }
  ↓
4. Call internal/cmd/configure.ConfigureFunc(target, opts)
  ↓
5. internal/cmd/configure.Configure():
   a. Create parent directory (os.MkdirAll)
   b. Check if file exists (os.Stat)
      - If exists and !Force: write message to ErrOutput, return nil
      - If Force: remove existing file (os.Remove)
   c. Marshal data (yaml.Marshal or json.MarshalIndent)
   d. Write file with os.OpenFile(target, flags, 0o644)
   e. Write success message to ErrOutput
   f. If opts.Edit: launch editor
      - Call opts.EditorLookup()
      - If error: write message to ErrOutput, return nil (error absorbed)
      - Create exec.Cmd with os.Stdin/Stdout/Stderr
      - Determine wait based on opts.EditorShouldWait()
      - Run editor via internal/proc.Run()
  ↓
6. Return error (if any) to cmd/configure.go RunE
```

## File Structure Changes

### New Files

- `internal/cmd/configure/configure.go`: Core business logic
- `internal/cmd/configure/configure_test.go`: Unit tests for Configure function

### Modified Files

- `cmd/configure.go`: Simplified to entry point + option building
- `cmd/configure_test.go`: Updated to use new structure
- `cmd/configure_wrapper_test.go`: Updated if necessary

### Deleted Files

No files will be deleted. `internal/cmd/configure.go` (existing) will be moved/replaced by the new structure.

## Dependencies

### Removed

- `github.com/rising3/go-cli/internal/stdio` - No longer imported

### Retained

- `github.com/spf13/cobra` - For CLI structure (`cmd/`)
- `github.com/spf13/viper` - For config management (`cmd/root.go`)
- `gopkg.in/yaml.v3` - For YAML marshaling
- `encoding/json` - For JSON marshaling (Go standard library)
- `os` - For file operations
- `os/exec` - For editor process launch
- `io` - For io.Writer interface
- `github.com/rising3/go-cli/internal/editor` - For editor detection
- `github.com/rising3/go-cli/internal/proc` - For process execution

## Backwards Compatibility

**File Permissions**: 設定ファイルは引き続き0644 (rw-r--r--)で作成される

**File Location**: `~/.config/mycli/default.yaml` および `~/.config/mycli/<profile>.yaml` は変更なし

**Command Behavior**: すべてのフラグ（`--force`, `--edit`, `--no-wait`, `--profile`）は既存と同じ動作を維持

**Output Format**: YAML/JSONフォーマットは既存実装と同一

**Error Handling**: エディタ検出失敗時のエラー吸収動作を維持（設定ファイル作成は成功とみなす）

## Testing Strategy

### Unit Tests (internal/cmd/configure/)

- `Configure()`関数の各ブランチをテスト:
  - ファイルが存在しない → 作成成功
  - ファイルが存在、Force=false → スキップ
  - ファイルが存在、Force=true → 上書き
  - ディレクトリ作成失敗 → エラー
  - ファイル書き込み失敗 → エラー
  - YAMLフォーマット → 正しいマーシャリング
  - JSONフォーマット → 正しいマーシャリング
  - Edit=true、エディタ検出成功 → エディタ起動
  - Edit=true、エディタ検出失敗 → エラー吸収
  - NoWait=true → バックグラウンド起動

- Mocking:
  - `Output`, `ErrOutput`: bytes.Buffer
  - `EditorLookup`: カスタム関数
  - `EditorShouldWait`: カスタム関数
  - `proc.ExecCommand`: テスト用モック（既存パターン踏襲）

### Integration Tests (cmd/)

- Cobraコマンドとしてのフラグ処理:
  - `--force`フラグの動作
  - `--edit`フラグの動作
  - `--no-wait`フラグの動作
  - `--profile`フラグの動作
  - フラグの組み合わせ
  - `cmd.SetOut()`, `cmd.SetErr()`でストリームをキャプチャ

### Coverage Goal

- 80%以上のカバレッジ（SC-005）
- 特にエッジケース（ファイル存在、権限エラー、エディタエラー）を網羅
