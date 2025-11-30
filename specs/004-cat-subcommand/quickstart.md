# Quickstart: Cat サブコマンド開発

**Feature**: Cat サブコマンド実装  
**Date**: 2025-11-30  
**Audience**: 開発者

## Overview

このガイドは、catサブコマンドの実装を始めるための手順を提供します。TDD（テスト駆動開発）のRed-Green-Refactorサイクルに従い、段階的に実装を進めます。

## Prerequisites

### 必須ツール

| Tool | Version | Install Command |
|------|---------|-----------------|
| Go | 1.25.4 | `brew install go@1.25` |
| golangci-lint | 2.6.2 | `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \| sh -s -- -b $(go env GOPATH)/bin v2.6.2` |
| bats | 1.x | `brew install bats-core` |

### 環境設定

```bash
# PATHにGo binディレクトリを追加
export PATH="$(go env GOPATH)/bin:$PATH"

# 依存パッケージのダウンロード
cd /Users/michio/dev/rising3/go-cli
go mod download
```

### 前提知識

- Go言語の基本構文
- Cobraフレームワークの基本的な使い方
- TDDの基本原則（Red-Green-Refactor）

## Project Structure

```
go-cli/
├── cmd/
│   ├── cat.go                      # 新規作成: Cobraコマンド定義
│   └── cat_test.go                 # 新規作成: コマンドの統合テスト
├── internal/
│   └── cmd/
│       └── cat/
│           ├── options.go          # 新規作成: Options構造体とファクトリ関数
│           ├── options_test.go     # 新規作成: オプション解決のテスト
│           ├── processor.go        # 新規作成: Processorインターフェースと実装
│           ├── processor_test.go   # 新規作成: ファイル処理のテスト
│           ├── formatter.go        # 新規作成: Formatterインターフェースと実装
│           └── formatter_test.go   # 新規作成: 行フォーマットのテスト
└── integration_test/
    └── cat.bats                    # 新規作成: BATS統合テスト
```

## Development Workflow

### Phase 1: Formatter実装（最も独立したコンポーネント）

#### Step 1.1: Formatterのインターフェースとテスト作成（Red）

```bash
# ファイル作成
touch internal/cmd/cat/formatter.go
touch internal/cmd/cat/formatter_test.go
```

`internal/cmd/cat/formatter_test.go`:
```go
package cat

import "testing"

func TestFormatLine_NumberAll(t *testing.T) {
    formatter := NewDefaultFormatter()
    opts := Options{NumberAll: true}
    
    got := formatter.FormatLine("Hello", 1, false, opts)
    want := "     1  Hello"
    
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

```bash
# テスト実行（失敗することを確認）
go test ./internal/cmd/cat/
```

#### Step 1.2: Formatterの最小実装（Green）

`internal/cmd/cat/formatter.go`:
```go
package cat

import "fmt"

type Formatter interface {
    FormatLine(line string, lineNum int, isEmpty bool, opts Options) string
}

type DefaultFormatter struct {
    controlCharMap map[byte]string
}

func NewDefaultFormatter() *DefaultFormatter {
    return &DefaultFormatter{
        controlCharMap: buildControlCharMap(),
    }
}

func (f *DefaultFormatter) FormatLine(line string, lineNum int, isEmpty bool, opts Options) string {
    result := line
    
    // 行番号付加
    if opts.NumberAll && !isEmpty {
        result = fmt.Sprintf("%6d  %s", lineNum, result)
    }
    
    return result
}

func buildControlCharMap() map[byte]string {
    // TODO: 後で実装
    return map[byte]string{}
}
```

```bash
# テスト実行（成功することを確認）
go test ./internal/cmd/cat/
```

#### Step 1.3: 追加テストケースとリファクタリング（Refactor）

テストを追加し、制御文字変換、行末マーカー、タブ変換などの機能を段階的に実装します。

### Phase 2: Processor実装

#### Step 2.1: Processorのインターフェースとテスト作成（Red）

```bash
touch internal/cmd/cat/processor.go
touch internal/cmd/cat/processor_test.go
```

`internal/cmd/cat/processor_test.go`:
```go
package cat

import (
    "bytes"
    "os"
    "testing"
)

func TestProcessFile_Success(t *testing.T) {
    // 一時ファイル作成
    tmpfile, err := os.CreateTemp(t.TempDir(), "test*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer tmpfile.Close()
    
    // テストデータ書き込み
    content := "line1\nline2\n"
    if _, err := tmpfile.WriteString(content); err != nil {
        t.Fatal(err)
    }
    
    // Processorのテスト
    formatter := NewDefaultFormatter()
    processor := NewDefaultProcessor(formatter)
    
    var output bytes.Buffer
    opts := Options{}
    
    if err := processor.ProcessFile(tmpfile.Name(), opts, &output); err != nil {
        t.Fatalf("ProcessFile failed: %v", err)
    }
    
    got := output.String()
    want := "line1\nline2\n"
    
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

#### Step 2.2: Processorの最小実装（Green）

`internal/cmd/cat/processor.go`:
```go
package cat

import (
    "bufio"
    "io"
    "os"
)

type Processor interface {
    ProcessFile(filename string, opts Options, output io.Writer) error
    ProcessStdin(opts Options, output io.Writer) error
}

type DefaultProcessor struct {
    formatter Formatter
}

func NewDefaultProcessor(formatter Formatter) *DefaultProcessor {
    return &DefaultProcessor{formatter: formatter}
}

func (p *DefaultProcessor) ProcessFile(filename string, opts Options, output io.Writer) error {
    if filename == "-" {
        return p.ProcessStdin(opts, output)
    }
    
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    return p.processReader(file, opts, output)
}

func (p *DefaultProcessor) ProcessStdin(opts Options, output io.Writer) error {
    return p.processReader(os.Stdin, opts, output)
}

func (p *DefaultProcessor) processReader(reader io.Reader, opts Options, output io.Writer) error {
    scanner := bufio.NewScanner(reader)
    lineNum := 0
    
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        isEmpty := len(line) == 0
        
        formatted := p.formatter.FormatLine(line, lineNum, isEmpty, opts)
        if _, err := output.Write([]byte(formatted + "\n")); err != nil {
            return err
        }
    }
    
    return scanner.Err()
}
```

### Phase 3: Options実装

#### Step 3.1: Optionsのテスト作成（Red）

```bash
touch internal/cmd/cat/options.go
touch internal/cmd/cat/options_test.go
```

`internal/cmd/cat/options_test.go`:
```go
package cat

import (
    "testing"
    
    "github.com/spf13/cobra"
)

func TestNewOptions_Default(t *testing.T) {
    cmd := &cobra.Command{}
    cmd.Flags().BoolP("number", "n", false, "")
    cmd.Flags().BoolP("number-nonblank", "b", false, "")
    
    opts, err := NewOptions(cmd)
    if err != nil {
        t.Fatal(err)
    }
    
    if opts.NumberAll {
        t.Error("NumberAll should be false")
    }
}
```

#### Step 3.2: Optionsの実装（Green）

`internal/cmd/cat/options.go`:
```go
package cat

import "github.com/spf13/cobra"

type Options struct {
    NumberAll       bool
    NumberNonBlank  bool
    ShowEnds        bool
    ShowTabs        bool
    ShowNonPrinting bool
}

func NewOptions(cmd *cobra.Command) (Options, error) {
    opts := Options{}
    
    numberAll, _ := cmd.Flags().GetBool("number")
    numberNonBlank, _ := cmd.Flags().GetBool("number-nonblank")
    showEnds, _ := cmd.Flags().GetBool("show-ends")
    showTabs, _ := cmd.Flags().GetBool("show-tabs")
    showNonPrinting, _ := cmd.Flags().GetBool("show-nonprinting")
    showAll, _ := cmd.Flags().GetBool("show-all")
    
    if showAll {
        showNonPrinting = true
        showEnds = true
        showTabs = true
    }
    
    opts.NumberAll = numberAll
    opts.NumberNonBlank = numberNonBlank
    opts.ShowEnds = showEnds
    opts.ShowTabs = showTabs
    opts.ShowNonPrinting = showNonPrinting
    
    return opts, nil
}
```

### Phase 4: Cobraコマンド統合

#### Step 4.1: catコマンドの作成

```bash
touch cmd/cat.go
touch cmd/cat_test.go
```

`cmd/cat.go`:
```go
package cmd

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/yourusername/go-cli/internal/cmd/cat"
)

var catCmd = &cobra.Command{
    Use:   "cat [flags] [file...]",
    Short: "Concatenate files and print on the standard output",
    Long: `Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.`,
    RunE: func(cmd *cobra.Command, args []string) error {
        opts, err := cat.NewOptions(cmd)
        if err != nil {
            return err
        }
        
        formatter := cat.NewDefaultFormatter()
        processor := cat.NewDefaultProcessor(formatter)
        
        if len(args) == 0 {
            return processor.ProcessStdin(opts, os.Stdout)
        }
        
        hadError := false
        for _, filename := range args {
            if err := processor.ProcessFile(filename, opts, os.Stdout); err != nil {
                fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
                hadError = true
            }
        }
        
        if hadError {
            os.Exit(1)
        }
        
        return nil
    },
}

func init() {
    rootCmd.AddCommand(catCmd)
    
    catCmd.Flags().BoolP("number", "n", false, "number all output lines")
    catCmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
    catCmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
    catCmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
    catCmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
    catCmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")
}
```

### Phase 5: BATS統合テスト

#### Step 5.1: BATS統合テストの作成

```bash
touch integration_test/cat.bats
```

`integration_test/cat.bats`:
```bats
#!/usr/bin/env bats

load 'helpers/common'
load 'helpers/test_env'

setup() {
    setup_test_env
}

teardown() {
    teardown_test_env
}

@test "cat displays file content" {
    echo "line1" > test.txt
    echo "line2" >> test.txt
    
    run mycli cat test.txt
    assert_success
    assert_output "line1
line2"
}

@test "cat with -n flag numbers all lines" {
    echo "line1" > test.txt
    echo "" >> test.txt
    echo "line3" >> test.txt
    
    run mycli cat -n test.txt
    assert_success
    assert_output "     1  line1
     2  
     3  line3"
}
```

## Testing

### 単体テスト実行

```bash
# 全体テスト
go test ./...

# 特定パッケージのテスト
go test ./internal/cmd/cat/

# カバレッジ付きテスト
go test -cover ./internal/cmd/cat/

# 詳細出力
go test -v ./internal/cmd/cat/
```

### BATS統合テスト実行

```bash
cd integration_test
bats cat.bats

# すべてのBATSテスト
bats *.bats
```

### コード品質チェック

```bash
# フォーマット
make fmt

# Lint
make lint

# すべてのチェック
make all
```

## Common Issues

### Issue 1: golangci-lintが見つからない

```bash
# PATHを確認
echo $PATH | grep "$(go env GOPATH)/bin"

# PATHに追加
export PATH="$(go env GOPATH)/bin:$PATH"

# golangci-lintを再インストール
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.2
```

### Issue 2: テストが失敗する

```bash
# 詳細なエラーメッセージを表示
go test -v ./internal/cmd/cat/

# 特定のテストのみ実行
go test -run TestFormatLine_NumberAll ./internal/cmd/cat/
```

### Issue 3: BATSテストが失敗する

```bash
# バイナリをビルド
make build

# バイナリのパスを確認
ls -la bin/mycli

# BATS環境変数を確認
cd integration_test
cat helpers/test_env.bash
```

## Next Steps

1. **実装開始**: `internal/cmd/cat/formatter.go`から実装を開始
2. **TDDサイクル**: Red-Green-Refactorを厳守
3. **段階的実装**: Formatter → Processor → Options → Cobra統合の順で実装
4. **統合テスト**: 各段階で`make all`を実行し、品質を確認

## Resources

- **仕様書**: `specs/004-cat-subcommand/spec.md`
- **技術調査**: `specs/004-cat-subcommand/research.md`
- **データモデル**: `specs/004-cat-subcommand/data-model.md`
- **契約定義**: `specs/004-cat-subcommand/contracts/`
- **憲章**: `.specify/memory/constitution.md`

## Questions?

実装中に不明点があれば、以下のドキュメントを参照してください：
- 仕様の詳細: `spec.md`
- 技術的決定の根拠: `research.md`
- データ構造の詳細: `data-model.md`
- インターフェース定義: `contracts/*.md`
