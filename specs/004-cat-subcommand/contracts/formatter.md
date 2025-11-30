# Contract: Formatter Interface

**Feature**: Cat サブコマンド  
**Component**: 行フォーマット層  
**Date**: 2025-11-30

## Interface Definition

```go
package cat

// Formatter は、行のフォーマット処理を担当するインターフェース
type Formatter interface {
    // FormatLine は1行をフォーマットして返す
    //
    // Parameters:
    //   - line: フォーマット対象の行（改行文字を含まない）
    //   - lineNum: 現在の行番号（1から開始）
    //   - isEmpty: 行が空行かどうか
    //   - opts: フォーマットオプション
    //
    // Returns:
    //   - string: フォーマット済みの行（改行文字を含まない）
    //
    // Formatting Rules:
    //   1. 行番号付加（opts.NumberAll または opts.NumberNonBlank が true の場合）
    //   2. 制御文字変換（opts.ShowNonPrinting が true の場合）
    //   3. タブ文字変換（opts.ShowTabs が true の場合）
    //   4. 行末マーカー（opts.ShowEnds が true の場合）
    FormatLine(line string, lineNum int, isEmpty bool, opts Options) string
}
```

## Default Implementation

```go
// DefaultFormatter は Formatter の標準実装
type DefaultFormatter struct {
    controlCharMap map[byte]string
}

// NewDefaultFormatter は DefaultFormatter を生成する
func NewDefaultFormatter() *DefaultFormatter {
    return &DefaultFormatter{
        controlCharMap: buildControlCharMap(),
    }
}

// buildControlCharMap は制御文字変換マップを構築する
func buildControlCharMap() map[byte]string {
    return map[byte]string{
        0: "^@", 1: "^A", 2: "^B", 3: "^C", 4: "^D", 5: "^E", 6: "^F", 7: "^G",
        8: "^H", /* 9: tab - skip */, /* 10: newline - skip */, 11: "^K", 12: "^L",
        13: "^M", 14: "^N", 15: "^O", 16: "^P", 17: "^Q", 18: "^R", 19: "^S",
        20: "^T", 21: "^U", 22: "^V", 23: "^W", 24: "^X", 25: "^Y", 26: "^Z",
        27: "^[", 28: "^\\", 29: "^]", 30: "^^", 31: "^_",
        127: "^?",
    }
}
```

## Usage Example

```go
formatter := NewDefaultFormatter()

opts := Options{
    NumberAll: true,
    ShowEnds:  true,
}

// 行をフォーマット
formatted := formatter.FormatLine("Hello, World!", 1, false, opts)
// Output: "     1  Hello, World!$"

// 空行をフォーマット（opts.NumberNonBlank = true の場合）
opts.NumberNonBlank = true
opts.NumberAll = false
formatted = formatter.FormatLine("", 2, true, opts)
// Output: "$" (行番号なし)

// 制御文字を含む行をフォーマット
opts.ShowNonPrinting = true
formatted = formatter.FormatLine("Hello\x07World", 3, false, opts)
// Output: "     3  Hello^GWorld$"
```

## Contract Guarantees

### Preconditions

1. `line`は改行文字を含まない
2. `lineNum`は正の整数（1以上）
3. `isEmpty`は`line`が空文字列の場合に`true`
4. `opts`は有効なOptions構造体

### Postconditions

1. 返り値は改行文字を含まない
2. 行番号は6桁右揃えでフォーマットされる
3. 999,999を超える行番号は最下位6桁のみ表示される

### Invariants

1. 同じ入力に対しては常に同じ出力を返す（副作用なし）
2. メモリアロケーションは最小限に抑えられる（`strings.Builder`使用）
3. 制御文字マップは初期化時に一度だけ構築される

## Formatting Rules

### 1. 行番号フォーマット

```go
// 6桁右揃え + 2スペース
format := "%6d  "

// 999,999を超える場合
lineNum = lineNum % 1000000

// 例:
//     1 → "     1  "
//   999 → "   999  "
// 1,234,567 → "234567  "
```

### 2. 空行の行番号スキップ

```go
if opts.NumberNonBlank && isEmpty {
    // 行番号を付加しない
    // 次の非空行で行番号カウンタを使用
}
```

### 3. 制御文字変換

```go
// ASCII 0-31（タブ9・改行10除く）と127を変換
if opts.ShowNonPrinting {
    for _, b := range []byte(line) {
        if converted, ok := controlCharMap[b]; ok {
            // 変換
        }
    }
}
```

### 4. タブ文字変換

```go
if opts.ShowTabs {
    // '\t' → "^I"
    line = strings.ReplaceAll(line, "\t", "^I")
}
```

### 5. 行末マーカー

```go
if opts.ShowEnds {
    // 行末に "$" を追加
    line = line + "$"
}
```

## Performance Characteristics

- **Time Complexity**: O(n)（nは行の長さ）
- **Space Complexity**: O(n)（フォーマット済み行のサイズ）
- **Target**: 行あたり10μs以下

## Testing Contract

### Unit Tests

```go
func TestFormatLine_NumberAll(t *testing.T)
func TestFormatLine_NumberNonBlank_EmptyLine(t *testing.T)
func TestFormatLine_NumberNonBlank_NonEmptyLine(t *testing.T)
func TestFormatLine_ShowEnds(t *testing.T)
func TestFormatLine_ShowTabs(t *testing.T)
func TestFormatLine_ShowNonPrinting(t *testing.T)
func TestFormatLine_LineNumberOverflow(t *testing.T)
func TestFormatLine_AllOptions(t *testing.T)
```

### Test Data

```go
var testCases = []struct {
    name     string
    line     string
    lineNum  int
    isEmpty  bool
    opts     Options
    expected string
}{
    {
        name:     "simple line with number",
        line:     "Hello",
        lineNum:  1,
        isEmpty:  false,
        opts:     Options{NumberAll: true},
        expected: "     1  Hello",
    },
    {
        name:     "line number overflow",
        line:     "Test",
        lineNum:  1234567,
        isEmpty:  false,
        opts:     Options{NumberAll: true},
        expected: "234567  Test",
    },
    {
        name:     "control characters",
        line:     "Test\x07\x1b",
        lineNum:  1,
        isEmpty:  false,
        opts:     Options{ShowNonPrinting: true},
        expected: "Test^G^[",
    },
}
```

## Edge Cases

### EC-001: 999,999を超える行番号

```go
Input:  lineNum = 1000000
Output: "000000  " (最下位6桁のみ)

Input:  lineNum = 1234567
Output: "234567  "
```

### EC-002: 空行（opts.NumberNonBlank = true）

```go
Input:  line = "", isEmpty = true, opts.NumberNonBlank = true
Output: "" (行番号なし)

Input:  line = "", isEmpty = true, opts.NumberAll = true
Output: "     1  " (行番号あり)
```

### EC-003: 制御文字とタブの混在

```go
Input:  line = "Hello\tWorld\x07", opts.ShowTabs = true, opts.ShowNonPrinting = true
Output: "Hello^IWorld^G"
```

### EC-004: 非常に長い行（1MB+）

```go
// Formatterはメモリ効率のため、strings.Builderを使用
// 1MB行のフォーマットでも100ms以下で処理可能
```

## Change History

- 2025-11-30: 初版作成
