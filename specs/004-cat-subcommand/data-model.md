# Data Model: Cat サブコマンド

**Feature**: Cat サブコマンド実装  
**Date**: 2025-11-30  
**Status**: Complete

## Overview

このドキュメントは、catサブコマンドの内部データ構造、状態遷移、バリデーションルールを定義する。

## Core Entities

### 1. Options

**Purpose**: コマンドラインフラグから抽出されたオプション設定を保持

**Structure**:
```go
type Options struct {
    // Line numbering options
    NumberAll      bool  // -n: number all lines
    NumberNonBlank bool  // -b: number non-empty lines only
    
    // Display options
    ShowEnds       bool  // -E: display $ at end of each line
    ShowTabs       bool  // -T: display TAB characters as ^I
    ShowNonPrinting bool // -v: use ^ and M- notation for control characters
}
```

**Validation Rules**:
- `-n`と`-b`が両方指定された場合、後に指定された方を優先（POSIXの動作に準拠）
- `-A`フラグが指定された場合、内部で`ShowNonPrinting=true`, `ShowEnds=true`, `ShowTabs=true`に変換

**Relationships**:
- `Processor`に渡され、ファイル処理の動作を制御
- `Formatter`に渡され、行フォーマットの方法を決定

### 2. Processor

**Purpose**: ファイルの読み込みとストリーム処理を担当

**Interface**:
```go
type Processor interface {
    ProcessFile(filename string, opts Options, output io.Writer) error
    ProcessStdin(opts Options, output io.Writer) error
}
```

**Implementation**: `DefaultProcessor`

**State**:
- 状態なし（ステートレス）
- 各メソッド呼び出しは独立

**Error Conditions**:
- ファイルが存在しない
- ファイルが読み取り不可
- ファイルがディレクトリ
- I/Oエラー（デバイスフルなど）

### 3. Formatter

**Purpose**: 行のフォーマット処理を担当（行番号付加、制御文字変換など）

**Interface**:
```go
type Formatter interface {
    FormatLine(line string, lineNum int, isEmpty bool, opts Options) string
}
```

**Implementation**: `DefaultFormatter`

**Internal State**:
- 制御文字変換マップ（事前構築された定数）

**Formatting Rules**:
1. **行番号フォーマット**: 
   - フォーマット: `"%6d  "` (6桁右揃え + 2スペース)
   - 999,999超過時: `lineNum % 1000000`で最下位6桁のみ表示
   - 空行のスキップ: `opts.NumberNonBlank`が`true`かつ`isEmpty`が`true`の場合、行番号を付加しない

2. **制御文字変換** (`opts.ShowNonPrinting`が`true`の場合):
   - ASCII 0-31（タブ9・改行10除く）: `^@`, `^A`, ..., `^_`に変換
   - ASCII 127（DEL）: `^?`に変換
   - タブ文字: `opts.ShowTabs`が`true`の場合のみ`^I`に変換

3. **行末マーカー** (`opts.ShowEnds`が`true`の場合):
   - 行末に`$`を追加

## State Transitions

### ファイル処理フロー

```
[開始] → [ファイルオープン] → [行読み込み] → [行フォーマット] → [出力] → [次の行]
                ↓                    ↓                              ↑
            [エラー]              [EOF]                            |
                ↓                    ↓                              |
            [stderrに出力]      [ファイルクローズ] ----------------+
                ↓                    ↓
            [次のファイル]         [完了]
```

### 行番号カウンタの状態

```
初期状態: lineNum = 0

各行の処理:
  1. lineNum++ (カウンタをインクリメント)
  2. 空行かつNumberNonBlank=true の場合:
     - 行番号を付加しない
     - カウンタは巻き戻さない（次の非空行で使用）
  3. それ以外:
     - 行番号を付加
  4. lineNum > 999,999 の場合:
     - lineNum % 1000000 で最下位6桁のみ表示
```

### オプション解決

```
[コマンドライン引数] → [Cobraフラグ解析]
                              ↓
                    [-n と -b の競合チェック]
                              ↓
                    [後に指定された方を優先]
                              ↓
                    [-A フラグの展開]
                              ↓
                    [Options構造体を構築]
                              ↓
                    [Processorに渡す]
```

## Data Flow

```
[ユーザー入力]
    ↓
[Cobraコマンド] (cmd/cat.go)
    ↓
[フラグ解析] → [Options構造体]
    ↓
[Processor.ProcessFile()] (internal/cmd/cat/processor.go)
    ↓
[bufio.Scanner] (32KBバッファ)
    ↓
[各行をFormatLine()に渡す] (internal/cmd/cat/formatter.go)
    ↓
[Formatter.FormatLine()]
    ├─ [行番号付加?]
    ├─ [制御文字変換?]
    └─ [行末マーカー?]
    ↓
[フォーマット済み行]
    ↓
[io.Writer (通常はos.Stdout)]
    ↓
[ユーザー出力]
```

## Validation Rules

### コマンドライン引数

| Rule | Description | Error Behavior |
|------|-------------|----------------|
| VR-001 | ファイル引数なしかつstdinがターミナル | エラー（UNIXのcat動作に準拠） |
| VR-002 | ファイルが存在しない | stderr出力後、次のファイルを処理 |
| VR-003 | ファイルがディレクトリ | stderr出力後、次のファイルを処理 |
| VR-004 | ファイルが読み取り不可 | stderr出力後、次のファイルを処理 |
| VR-005 | `-`引数が指定された場合 | stdinから読み込み |

### オプション競合

| Rule | Description | Resolution |
|------|-------------|------------|
| VR-006 | `-n`と`-b`が両方指定 | 後に指定された方を優先 |
| VR-007 | `-A`が指定 | 内部で`-v -E -T`に展開 |

### 出力フォーマット

| Rule | Description | Validation |
|------|-------------|------------|
| VR-008 | 行番号は6桁固定 | `%6d`で右揃え |
| VR-009 | 行番号が999,999超過 | `% 1000000`で最下位6桁のみ |
| VR-010 | 制御文字はASCII 0-31 + 127のみ変換 | タブ9・改行10は除外 |
| VR-011 | 行末マーカーは`$` | 改行の直前に挿入 |

## Performance Constraints

### メモリ使用量

| Component | Constraint | Rationale |
|-----------|------------|-----------|
| バッファサイズ | 32KB/ファイル | 明確化セッションで決定 |
| 制御文字マップ | ~1KB（静的） | 事前構築された定数 |
| 最大メモリ | 100MB（1GBファイル処理時） | 仕様書SC-007 |

### 処理速度

| Operation | Target | Rationale |
|-----------|--------|-----------|
| 1MBファイル処理 | 100ms以下 | 仕様書SC-006 |
| 行フォーマット | 行あたり10μs以下 | ストリーム処理のオーバーヘッド最小化 |

## Edge Cases

### EC-001: 巨大ファイル（1GB+）

**Behavior**:
- ストリーム処理により、メモリ使用量は一定（約32KB + オーバーヘッド）
- 行番号が999,999を超えた場合、最下位6桁のみ表示

**Data Model Impact**:
- `Processor`はファイルサイズに関係なく一定のメモリを使用
- 行番号カウンタは`int`型（64ビット環境で約900京まで対応）

### EC-002: 空ファイル

**Behavior**:
- 出力なし（エラーも出力しない）
- 終了コード0

**Data Model Impact**:
- `Processor`は0行を処理
- 行番号カウンタは0のまま

### EC-003: バイナリファイル

**Behavior**:
- `-v`オプションが指定された場合、制御文字を変換して出力
- 指定されていない場合、そのまま出力

**Data Model Impact**:
- `Formatter`は制御文字変換マップを参照
- バイナリデータも文字列として処理

### EC-004: 非常に長い行（1MB+）

**Behavior**:
- `bufio.Scanner`はデフォルトで最大64KBの行をサポート
- 32KBバッファを使用するため、32KB超の行はエラー

**Data Model Impact**:
- `Processor`は`bufio.Scanner`のエラーをキャッチし、stderrに出力
- 次の行の処理を継続

**Mitigation**:
- `Scanner.Buffer()`で最大トークンサイズを設定可能
- 実装時に最大行長を設定（例: 1MB）

## Implementation Checklist

- [ ] `Options`構造体を`internal/cmd/cat/options.go`に定義
- [ ] `Processor`インターフェースを`internal/cmd/cat/processor.go`に定義
- [ ] `Formatter`インターフェースを`internal/cmd/cat/formatter.go`に定義
- [ ] 制御文字変換マップを`formatter.go`に定義
- [ ] オプション競合解決ロジックを`options.go`に実装
- [ ] 各エンティティの単体テストを作成
- [ ] エッジケースのテストを作成

## References

- 仕様書: `specs/004-cat-subcommand/spec.md`
- 技術調査: `specs/004-cat-subcommand/research.md`
- 憲章: `.specify/memory/constitution.md`
