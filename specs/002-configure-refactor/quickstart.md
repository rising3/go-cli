# Quickstart: Configure サブコマンドのリファクタリング

**Feature**: [spec.md](./spec.md) | **Plan**: [plan.md](./plan.md) | **Data Model**: [data-model.md](./data-model.md)  
**Created**: 2025-11-30

## Overview

このガイドは、configureサブコマンドのリファクタリングを実装する開発者向けの手順書です。echoサブコマンドのパターンに従い、TDD（テスト駆動開発）で進めます。

## Prerequisites

```bash
# Go 1.25.4がインストールされていることを確認
go version  # go version go1.25.4 ...

# golangci-lint v2.6.2がインストールされていることを確認
golangci-lint version  # golangci-lint has version 2.6.2 ...

# PATHにgolangci-lintが含まれていることを確認
export PATH="$(go env GOPATH)/bin:$PATH"

# 依存関係をダウンロード
go mod download
```

## Implementation Steps

### Phase 1: internal/cmd/configure/ パッケージの作成

#### Step 1.1: ディレクトリ構造の準備

```bash
# internal/cmd/configure/ ディレクトリを作成
mkdir -p internal/cmd/configure

# 既存のinternal/cmd/configure.goをバックアップ（後で参照）
cp internal/cmd/configure.go internal/cmd/configure.go.backup
```

#### Step 1.2: ConfigureOptions構造体の定義（TDD: Test First）

**File**: `internal/cmd/configure/configure_test.go`

```go
package configure_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/rising3/go-cli/internal/cmd/configure"
)

func TestConfigure_BasicFileCreation(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.yaml")
	
	data := map[string]interface{}{
		"key": "value",
	}
	
	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            false,
		Edit:             false,
		NoWait:           false,
		Data:             data,
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}
	
	// Execute
	err := configure.Configure(target, opts)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Errorf("file was not created: %s", target)
	}
	
	// Verify file content
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	
	expected := "key: value\n"
	if string(content) != expected {
		t.Errorf("content = %q, want %q", string(content), expected)
	}
}
```

**Run test** (should fail):

```bash
go test ./internal/cmd/configure/
# Error: package configure not found
```

#### Step 1.3: 最小実装（Green）

**File**: `internal/cmd/configure/configure.go`

```go
package configure

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigureOptions represents the configuration for the configure command.
type ConfigureOptions struct {
	Force            bool
	Edit             bool
	NoWait           bool
	Data             map[string]interface{}
	Format           string
	Output           io.Writer
	ErrOutput        io.Writer
	EditorLookup     func() (string, []string, error)
	EditorShouldWait func(string, []string) bool
}

// Configure creates or overwrites a configuration file at the specified target path.
func Configure(target string, opts ConfigureOptions) error {
	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	
	// Check if file exists
	if _, err := os.Stat(target); err == nil && !opts.Force {
		if opts.ErrOutput != nil {
			fmt.Fprintln(opts.ErrOutput, "Config already exists, skipping initialization:", target)
		}
		return nil
	}
	
	// Remove existing file if Force
	if opts.Force {
		_ = os.Remove(target)
	}
	
	// Marshal data
	var out []byte
	var err error
	switch opts.Format {
	case "yaml", "yml":
		out, err = yaml.Marshal(opts.Data)
	default:
		out, err = json.MarshalIndent(opts.Data, "", "  ")
	}
	if err != nil {
		return err
	}
	
	// Write file
	if err := os.WriteFile(target, out, 0o644); err != nil {
		return err
	}
	
	if opts.ErrOutput != nil {
		fmt.Fprintln(opts.ErrOutput, "Wrote config:", target)
	}
	
	// TODO: Editor launch will be implemented in next step
	
	return nil
}

// ConfigureFunc is a variable indirection for testing
var ConfigureFunc = Configure
```

**Run test** (should pass):

```bash
go test ./internal/cmd/configure/
# PASS
```

#### Step 1.4: 追加テストケースの実装

```go
func TestConfigure_FileExists_NoForce(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "existing.yaml")
	
	// Create existing file
	os.WriteFile(target, []byte("old: data\n"), 0o644)
	
	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            false,
		Data:             map[string]interface{}{"new": "data"},
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}
	
	err := configure.Configure(target, opts)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Verify file content unchanged
	content, _ := os.ReadFile(target)
	if string(content) != "old: data\n" {
		t.Errorf("file was modified when it shouldn't be")
	}
	
	// Verify message
	if !strings.Contains(errBuf.String(), "already exists") {
		t.Errorf("expected 'already exists' message, got: %s", errBuf.String())
	}
}

func TestConfigure_FileExists_Force(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "existing.yaml")
	
	// Create existing file
	os.WriteFile(target, []byte("old: data\n"), 0o644)
	
	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            true,
		Data:             map[string]interface{}{"new": "data"},
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}
	
	err := configure.Configure(target, opts)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Verify file content updated
	content, _ := os.ReadFile(target)
	expected := "new: data\n"
	if string(content) != expected {
		t.Errorf("content = %q, want %q", string(content), expected)
	}
}
```

### Phase 2: cmd/configure.go のリファクタリング

#### Step 2.1: テストの更新（Test First）

**File**: `cmd/configure_test.go`

既存のテストを確認し、新しい構造に合わせて更新します。

```go
package cmd

import (
	"bytes"
	"testing"
	
	"github.com/rising3/go-cli/internal/cmd/configure"
)

func TestConfigureCommand(t *testing.T) {
	// Setup: Mock ConfigureFunc
	originalFunc := configure.ConfigureFunc
	defer func() { configure.ConfigureFunc = originalFunc }()
	
	var capturedTarget string
	var capturedOpts configure.ConfigureOptions
	
	configure.ConfigureFunc = func(target string, opts configure.ConfigureOptions) error {
		capturedTarget = target
		capturedOpts = opts
		return nil
	}
	
	// Execute
	cmd := rootCmd
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"configure", "--force"})
	
	err := cmd.Execute()
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if !capturedOpts.Force {
		t.Error("Force flag not passed correctly")
	}
}
```

#### Step 2.2: cmd/configure.go の実装

```go
package cmd

import (
	"path/filepath"

	"github.com/rising3/go-cli/internal/cmd/configure"
	"github.com/rising3/go-cli/internal/editor"
	"github.com/spf13/cobra"
)

var cfgForce bool
var cfgEdit bool
var cfgNoWait bool

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().BoolVar(&cfgForce, "force", false, "overwrite existing config")
	configureCmd.Flags().BoolVar(&cfgEdit, "edit", false, "edit the created file in $EDITOR")
	configureCmd.Flags().BoolVar(&cfgNoWait, "no-wait", false, "do not wait for editor to exit")
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create a scaffold config file based on Config struct",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine target path
		dir := GetConfigPath()
		cfgName := DefaultProfile
		if profile != "" {
			cfgName = profile
		}
		target := filepath.Join(dir, GetConfigFile(cfgName))
		
		// Build options
		opts := configure.ConfigureOptions{
			Force:            cfgForce,
			Edit:             cfgEdit,
			NoWait:           cfgNoWait,
			Data:             BuildEffectiveConfig(),
			Format:           CliConfigType,
			Output:           cmd.OutOrStdout(),
			ErrOutput:        cmd.ErrOrStderr(),
			EditorLookup:     func() (string, []string, error) { return editor.GetEditor() },
			EditorShouldWait: func(string, []string) bool { return !cfgNoWait },
		}
		
		// Call internal function
		return configure.ConfigureFunc(target, opts)
	},
}
```

### Phase 3: エディタ起動機能の実装

#### Step 3.1: テストの追加

```go
func TestConfigure_Edit_EditorFound(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.yaml")
	
	editorCalled := false
	var errBuf bytes.Buffer
	
	opts := configure.ConfigureOptions{
		Force:     false,
		Edit:      true,
		NoWait:    false,
		Data:      map[string]interface{}{"key": "value"},
		Format:    "yaml",
		Output:    &bytes.Buffer{},
		ErrOutput: &errBuf,
		EditorLookup: func() (string, []string, error) {
			return "/usr/bin/vi", []string{}, nil
		},
		EditorShouldWait: func(string, []string) bool { return true },
	}
	
	// Mock proc.ExecCommand
	// (implementation details depend on existing test helpers)
	
	err := configure.Configure(target, opts)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Verify editor was called (check via mock)
}
```

#### Step 3.2: Editor launch の実装

```go
// Add to Configure() function after file write:

if opts.Edit {
	ed, edArgs, err := opts.EditorLookup()
	if err != nil {
		if opts.ErrOutput != nil {
			fmt.Fprintln(opts.ErrOutput, "No editor found:", err)
		}
		return nil // Error absorbed
	}
	
	args := append(edArgs, target)
	cmd := proc.ExecCommand(ed, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	shouldWait := true
	if opts.EditorShouldWait != nil {
		shouldWait = opts.EditorShouldWait(ed, args)
	}
	
	return proc.Run(cmd, shouldWait, opts.ErrOutput)
}
```

### Phase 4: 品質チェックとクリーンアップ

#### Step 4.1: すべてのテストを実行

```bash
# Run all tests
make test

# Check coverage
go test -cover ./cmd/... ./internal/cmd/configure/...
# Target: 80%以上
```

#### Step 4.2: フォーマットとリント

```bash
# Format code
make fmt

# Run linter
make lint
```

#### Step 4.3: 既存の `internal/stdio` 参照を削除

```bash
# Check for stdio references
grep -r "internal/stdio" cmd/configure.go internal/cmd/configure/
# Should return 0 results

# Remove internal/cmd/configure.go.backup if everything works
rm internal/cmd/configure.go.backup
```

#### Step 4.4: 統合テスト

```bash
# Build and test manually
make build

# Test configure command
./bin/mycli configure --force
cat ~/.config/mycli/default.yaml

# Test with editor (if available)
./bin/mycli configure --force --edit

# Test with profile
./bin/mycli configure --force --profile test
cat ~/.config/mycli/test.yaml
```

## Verification Checklist

### Constitution Check (再評価)

- [x] **TDD必須**: すべての新規コードにテストが存在するか？
- [x] **パッケージ責務分離**: `cmd/`と`internal/`が分離されているか？
- [x] **コード品質基準**: `make lint`がパスするか？
- [x] **設定管理の一貫性**: 既存の設定パスとフォーマットを維持しているか？
- [x] **ユーザーエクスペリエンス**: Cobraストリームを使用しているか？
- [x] **パフォーマンス要件**: パフォーマンスが劣化していないか？

### Success Criteria

- [ ] SC-001: `grep -r "internal/stdio" cmd/configure.go internal/cmd/configure/` が0件
- [ ] SC-002: `cmd/configure.go`が`cmd.OutOrStdout()`/`cmd.ErrOrStderr()`を使用
- [ ] SC-003: `Configure(target string, opts ConfigureOptions) error`が実装済み
- [ ] SC-004: 既存テストが100%パス
- [ ] SC-005: カバレッジ80%以上
- [ ] SC-006: `mycli configure --force`の出力が変わっていない
- [ ] SC-007: `cmd/echo.go`と`cmd/configure.go`の構造が一致（±3行）
- [ ] SC-008: `make lint`がパス

## Troubleshooting

### テストが失敗する場合

```bash
# 詳細なテスト出力
go test -v ./internal/cmd/configure/

# 特定のテストのみ実行
go test -run TestConfigure_BasicFileCreation ./internal/cmd/configure/
```

### リンターエラー

```bash
# govetのみ実行（憲章に準拠）
golangci-lint run --enable=govet

# エラーの詳細を確認
golangci-lint run --enable=govet -v
```

### エディタが起動しない

```bash
# エディタ環境変数を確認
echo $EDITOR
echo $VISUAL

# エディタを手動で設定
export EDITOR=vim
./bin/mycli configure --force --edit
```

## Next Steps

このリファクタリング完了後:

1. 既存のテストがすべてパスすることを確認
2. 新しいテストケースを追加してカバレッジを向上
3. ドキュメント（README、コメント）を更新
4. PRを作成し、コードレビューを依頼

## References

- [Feature Specification](./spec.md)
- [Implementation Plan](./plan.md)
- [Data Model](./data-model.md)
- [Configure Function Contract](./contracts/configure-function.md)
- [Echo Subcommand Implementation](../001-echo-subcommand/) (reference pattern)
