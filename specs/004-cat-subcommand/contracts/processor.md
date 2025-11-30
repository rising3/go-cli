# Contract: Processor Interface

**Feature**: Cat サブコマンド  
**Component**: ファイル処理層  
**Date**: 2025-11-30

## Interface Definition

```go
package cat

import "io"

// Processor は、ファイルの読み込みとストリーム処理を担当するインターフェース
type Processor interface {
    // ProcessFile は指定されたファイルを読み込み、フォーマットして出力する
    //
    // Parameters:
    //   - filename: 読み込むファイルのパス（"-"の場合はstdinから読み込み）
    //   - opts: フォーマットオプション
    //   - output: 出力先（通常はos.Stdout）
    //
    // Returns:
    //   - error: ファイルのオープン、読み込み、書き込みエラー
    //
    // Error Conditions:
    //   - ファイルが存在しない: os.ErrNotExist
    //   - ファイルがディレクトリ: syscall.EISDIR相当のエラー
    //   - 読み取り権限がない: os.ErrPermission
    //   - I/Oエラー: io.ErrUnexpectedEOF など
    ProcessFile(filename string, opts Options, output io.Writer) error
    
    // ProcessStdin はstdinから読み込み、フォーマットして出力する
    //
    // Parameters:
    //   - opts: フォーマットオプション
    //   - output: 出力先（通常はos.Stdout）
    //
    // Returns:
    //   - error: 読み込み、書き込みエラー
    ProcessStdin(opts Options, output io.Writer) error
}
```

## Default Implementation

```go
// DefaultProcessor は Processor の標準実装
type DefaultProcessor struct {
    formatter Formatter
}

// NewDefaultProcessor は DefaultProcessor を生成する
func NewDefaultProcessor(formatter Formatter) *DefaultProcessor {
    return &DefaultProcessor{
        formatter: formatter,
    }
}
```

## Usage Example

```go
formatter := NewDefaultFormatter()
processor := NewDefaultProcessor(formatter)

opts := Options{
    NumberAll: true,
    ShowEnds:  true,
}

// ファイルを処理
if err := processor.ProcessFile("example.txt", opts, os.Stdout); err != nil {
    fmt.Fprintf(os.Stderr, "cat: example.txt: %v\n", err)
}

// stdinを処理
if err := processor.ProcessStdin(opts, os.Stdout); err != nil {
    fmt.Fprintf(os.Stderr, "cat: stdin: %v\n", err)
}
```

## Contract Guarantees

### Preconditions

1. `filename`が`"-"`の場合、`ProcessFile`は自動的に`ProcessStdin`の動作を行う
2. `output`はnil以外のio.Writer
3. `opts`は有効なOptions構造体

### Postconditions

1. 成功時、すべての行がフォーマットされて`output`に書き込まれる
2. エラー時、処理途中の行も出力される（部分的成功）
3. ファイルは処理後に確実にクローズされる（deferで保証）

### Invariants

1. 処理は行単位でストリーム処理される
2. メモリ使用量はファイルサイズに依存しない（最大32KB + オーバーヘッド）
3. 行番号カウンタは単調増加する（巻き戻しなし）

## Error Handling

### Error Types

| Error Type | Description | Example |
|------------|-------------|---------|
| `os.ErrNotExist` | ファイルが存在しない | `cat: example.txt: no such file or directory` |
| `os.ErrPermission` | 読み取り権限がない | `cat: example.txt: permission denied` |
| `syscall.EISDIR` | ディレクトリを指定 | `cat: example.txt: Is a directory` |
| `io.ErrUnexpectedEOF` | ファイル読み込み中にEOF | `cat: example.txt: unexpected EOF` |

### Error Format

```go
fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
```

## Performance Characteristics

- **Time Complexity**: O(n)（nはファイルのバイト数）
- **Space Complexity**: O(1)（32KBバッファ + 定数オーバーヘッド）
- **Throughput**: 1MBファイルで100ms以下（仕様書SC-006）

## Testing Contract

### Unit Tests

```go
func TestProcessFile_Success(t *testing.T)
func TestProcessFile_NotExist(t *testing.T)
func TestProcessFile_IsDirectory(t *testing.T)
func TestProcessFile_PermissionDenied(t *testing.T)
func TestProcessStdin_Success(t *testing.T)
```

### Integration Tests

```bats
@test "cat processes file successfully" {
    echo "line1" > test.txt
    run mycli cat test.txt
    [ "$status" -eq 0 ]
    [ "$output" = "line1" ]
}

@test "cat reports error for non-existent file" {
    run mycli cat nonexistent.txt
    [ "$status" -eq 1 ]
    [[ "$stderr" =~ "no such file" ]]
}
```

## Change History

- 2025-11-30: 初版作成
