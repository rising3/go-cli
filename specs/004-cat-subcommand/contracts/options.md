# Contract: Options Structure

**Feature**: Cat サブコマンド  
**Component**: オプション管理層  
**Date**: 2025-11-30

## Structure Definition

```go
package cat

// Options は catコマンドのフォーマットオプションを保持する
type Options struct {
    // NumberAll は、すべての行に行番号を付加するかどうか (-n)
    NumberAll bool
    
    // NumberNonBlank は、空行以外に行番号を付加するかどうか (-b)
    NumberNonBlank bool
    
    // ShowEnds は、各行の末尾に $ を表示するかどうか (-E)
    ShowEnds bool
    
    // ShowTabs は、タブ文字を ^I として表示するかどうか (-T)
    ShowTabs bool
    
    // ShowNonPrinting は、制御文字を ^ 記法で表示するかどうか (-v)
    ShowNonPrinting bool
}
```

## Factory Functions

```go
// NewOptions は Cobraコマンドのフラグから Options を生成する
//
// Parameters:
//   - cmd: Cobraコマンド（フラグ情報を含む）
//
// Returns:
//   - Options: 解決済みのオプション
//   - error: フラグ取得エラー
func NewOptions(cmd *cobra.Command) (Options, error) {
    opts := Options{}
    
    // フラグを取得
    numberAll, _ := cmd.Flags().GetBool("number")
    numberNonBlank, _ := cmd.Flags().GetBool("number-nonblank")
    showEnds, _ := cmd.Flags().GetBool("show-ends")
    showTabs, _ := cmd.Flags().GetBool("show-tabs")
    showNonPrinting, _ := cmd.Flags().GetBool("show-nonprinting")
    showAll, _ := cmd.Flags().GetBool("show-all")
    
    // -A フラグの展開
    if showAll {
        showNonPrinting = true
        showEnds = true
        showTabs = true
    }
    
    // -n と -b の競合解決
    if numberAll && numberNonBlank {
        // POSIXの動作: 後に指定された方を優先
        // Cobraは最後に指定されたフラグを優先するため、特別な処理は不要
    }
    
    opts.NumberAll = numberAll
    opts.NumberNonBlank = numberNonBlank
    opts.ShowEnds = showEnds
    opts.ShowTabs = showTabs
    opts.ShowNonPrinting = showNonPrinting
    
    return opts, nil
}
```

## Usage Example

```go
var catCmd = &cobra.Command{
    Use:   "cat [flags] [file...]",
    Short: "Concatenate files and print on the standard output",
    RunE: func(cmd *cobra.Command, args []string) error {
        // オプションを生成
        opts, err := NewOptions(cmd)
        if err != nil {
            return err
        }
        
        // ファイルを処理
        processor := NewDefaultProcessor(NewDefaultFormatter())
        for _, filename := range args {
            if err := processor.ProcessFile(filename, opts, os.Stdout); err != nil {
                fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
            }
        }
        return nil
    },
}

func init() {
    catCmd.Flags().BoolP("number", "n", false, "number all output lines")
    catCmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
    catCmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
    catCmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
    catCmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
    catCmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")
}
```

## Contract Guarantees

### Preconditions

1. Cobraコマンドのフラグが正しく定義されている
2. フラグの型がboolである

### Postconditions

1. `-A`が指定された場合、`ShowNonPrinting`, `ShowEnds`, `ShowTabs`がすべて`true`になる
2. `-n`と`-b`が両方指定された場合、後に指定された方が優先される
3. 返り値の`Options`はすべてのフラグを解決済み

### Invariants

1. `NumberAll`と`NumberNonBlank`が両方`true`になることはない（後に指定された方のみ`true`）
2. `-A`フラグは他のフラグを上書きする（`-vET`と等価）

## Validation Rules

### VR-001: -n と -b の競合

```go
// POSIX準拠: 後に指定された方を優先
// 例: `cat -n -b file.txt` → -b が有効
// 例: `cat -b -n file.txt` → -n が有効
```

**Implementation Note**: Cobraのフラグ解析は最後に指定されたフラグを優先するため、特別な処理は不要。

### VR-002: -A フラグの展開

```go
if showAll {
    ShowNonPrinting = true  // -v
    ShowEnds = true         // -E
    ShowTabs = true         // -T
}
```

**Implementation Note**: `-A`が指定された場合、他のフラグの値に関係なく、上記3つのフラグが`true`に設定される。

## Testing Contract

### Unit Tests

```go
func TestNewOptions_Default(t *testing.T)
func TestNewOptions_NumberAll(t *testing.T)
func TestNewOptions_NumberNonBlank(t *testing.T)
func TestNewOptions_ShowAll(t *testing.T)
func TestNewOptions_NumberConflict_NFirst(t *testing.T)
func TestNewOptions_NumberConflict_BFirst(t *testing.T)
```

### Test Data

```go
var testCases = []struct {
    name     string
    flags    map[string]bool
    expected Options
}{
    {
        name:  "default options",
        flags: map[string]bool{},
        expected: Options{
            NumberAll:       false,
            NumberNonBlank:  false,
            ShowEnds:        false,
            ShowTabs:        false,
            ShowNonPrinting: false,
        },
    },
    {
        name: "show all",
        flags: map[string]bool{"show-all": true},
        expected: Options{
            NumberAll:       false,
            NumberNonBlank:  false,
            ShowEnds:        true,
            ShowTabs:        true,
            ShowNonPrinting: true,
        },
    },
}
```

## Edge Cases

### EC-001: -n と -b が両方指定

```go
// 後に指定された方を優先
$ mycli cat -n -b file.txt  # -b が有効
$ mycli cat -b -n file.txt  # -n が有効
```

### EC-002: -A と他のフラグの組み合わせ

```go
// -A は常に -vET を有効にする
$ mycli cat -A file.txt           # -v -E -T が有効
$ mycli cat -A -E file.txt        # -v -E -T が有効（重複しても問題なし）
```

### EC-003: すべてのフラグが無効

```go
// デフォルトの動作（オプションなし）
$ mycli cat file.txt
// Options{NumberAll: false, NumberNonBlank: false, ShowEnds: false, ShowTabs: false, ShowNonPrinting: false}
```

## Performance Characteristics

- **Time Complexity**: O(1)（フラグ数は固定）
- **Space Complexity**: O(1)（Options構造体のサイズは固定）

## Change History

- 2025-11-30: 初版作成
