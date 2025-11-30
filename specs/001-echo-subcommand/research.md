# Phase 0: Research & Technical Decisions

**Feature**: Echo サブコマンド実装  
**Date**: 2025-11-30  
**Status**: Completed

## Research Objectives

1. UNIX `echo` コマンドの標準動作仕様を明確化
2. Goでのエスケープシーケンス処理のベストプラクティス
3. Cobraでのエラーハンドリングとヘルプメッセージのカスタマイズ方法
4. TDDアプローチにおける標準出力テストのパターン

## 1. UNIX Echo Command Specification

### Decision: POSIX準拠の基本動作 + 一般的な拡張をサポート

**Research Summary**:
- POSIX標準の`echo`は、引数をスペース区切りで出力し、デフォルトで末尾に改行を追加
- `-n`オプション（改行抑制）は事実上の標準（GNU, BSD, BusyBoxすべてでサポート）
- `-e`オプション（エスケープシーケンス解釈）もGNU, BusyBoxで標準サポート
- GNU coreutilsの`echo`は最も広く使用されているリファレンス実装

**Rationale**:
- UNIX互換性を最優先することで、既存のシェルスクリプトやドキュメントとの相互運用性を確保
- `-n`, `-e`の2つのオプションは実用上ほぼ必須の機能
- 複雑なオプション（`-E`など）は除外し、シンプルさを保つ

**Alternatives Considered**:
- **Alternative A**: POSIXのみ準拠（`-n`なし） → 却下理由: 実用性が低く、ユーザー期待に反する
- **Alternative B**: GNU coreutilsの全オプション実装 → 却下理由: 複雑すぎる（`-E`, `--version`など不要）
- **Alternative C**: Busybox echo互換 → 却下理由: GNU coreutilsの方が広く使用されている

**Reference Implementation**: GNU coreutils `echo.c` (https://github.com/coreutils/coreutils/blob/master/src/echo.c)

---

## 2. Escape Sequence Processing in Go

### Decision: カスタムパーサーで標準エスケープシーケンスを処理

**Research Summary**:
- Goの`strconv.Unquote()`は引用符付き文字列専用で、`echo`のユースケースには不適
- `strings.Replacer`は単純な置換だが、`\c`（以降の出力抑制）のような複雑なシーケンスに対応不可
- カスタムパーサーが最も柔軟で、エスケープシーケンスの仕様を完全に制御可能

**Rationale**:
- `\c`（出力抑制）や`\0nnn`（8進数）などの複雑な動作を正確に実装するには、カスタムパーサーが必要
- Goの`strings.Builder`で効率的な文字列構築が可能
- テスト可能な独立関数として実装できる（`internal/echo/processor.go`）

**Implementation Pattern**:
```go
// internal/echo/processor.go
func ProcessEscapes(input string) (output string, suppressNewline bool) {
    var builder strings.Builder
    i := 0
    for i < len(input) {
        if input[i] == '\\' && i+1 < len(input) {
            switch input[i+1] {
            case 'n': builder.WriteRune('\n'); i += 2
            case 't': builder.WriteRune('\t'); i += 2
            case '\\': builder.WriteRune('\\'); i += 2
            case 'c': return builder.String(), true // 出力抑制
            default: builder.WriteByte(input[i]); i++
            }
        } else {
            builder.WriteByte(input[i]); i++
        }
    }
    return builder.String(), false
}
```

**Alternatives Considered**:
- **Alternative A**: `strconv.Unquote()` → 却下理由: 引用符が必須で`echo`の仕様に合わない
- **Alternative B**: `strings.Replacer` → 却下理由: `\c`のような条件付き動作を実装できない
- **Alternative C**: 正規表現による置換 → 却下理由: パフォーマンスが悪く、複雑なエスケープに対応困難

**Best Practices**:
- エスケープ処理は`internal/echo/processor.go`に独立させ、Cobraから分離
- `\c`が検出された場合は即座にreturnし、それ以降の処理をスキップ
- 無効なエスケープシーケンスはリテラル文字列として扱う（GNU標準に準拠）

---

## 3. Cobra Error Handling & Help Customization

### Decision: `SilenceUsage: false`でエラー後の自動ヘルプ表示を有効化

**Research Summary**:
- Cobraの`SilenceUsage`フラグがヘルプメッセージの表示を制御
- デフォルト（`false`）では、コマンド実行エラー時に自動的にヘルプが表示される
- 無効なフラグは自動的にCobraが検出し、エラーメッセージ + ヘルプを表示

**Rationale**:
- FR-012（無効オプション時の自動ヘルプ表示）はCobraのデフォルト動作で自然に実装可能
- ユーザーフレンドリーな体験を追加コードなしで実現
- `SilenceUsage: true`にすると、ヘルプが表示されずユーザビリティが低下

**Implementation Pattern**:
```go
// cmd/echo.go
var echoCmd = &cobra.Command{
    Use:   "echo [flags] [args...]",
    Short: "Output text to stdout (UNIX echo clone)",
    Long: `Output arguments separated by spaces, with a trailing newline by default.
Supports -n (suppress newline) and -e (interpret escape sequences).`,
    Example: `  mycli echo "Hello, World!"
  mycli echo -n "No newline"
  mycli echo -e "Line1\nLine2"`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // 実装ロジック
        return nil
    },
    SilenceUsage: false, // エラー時に自動的にヘルプを表示
}
```

**Alternatives Considered**:
- **Alternative A**: `SilenceUsage: true` + カスタムエラーハンドラでヘルプ表示 → 却下理由: 不要な複雑さ
- **Alternative B**: `PreRunE`でフラグ検証 → 却下理由: Cobraが既に提供している機能の重複

**Best Practices**:
- `Use`, `Short`, `Long`, `Example`フィールドを適切に設定し、`--help`の品質を確保
- `Example`には2-3個の実用的なユースケースを含める（FR-008）
- エラーはstderrに出力されることを利用し、正常出力（stdout）と明確に分離

---

## 4. TDD Patterns for Standard Output Testing

### Decision: `bytes.Buffer`でstdout/stderrをキャプチャし、テストで検証

**Research Summary**:
- Goの`os.Stdout`は`*os.File`型でモック不可だが、`io.Writer`インターフェースにリダイレクト可能
- Cobraの`cmd.SetOut()`, `cmd.SetErr()`で出力先を変更可能
- `bytes.Buffer`は`io.Writer`を実装し、出力内容の検証に最適

**Rationale**:
- Cobraの標準機能でテスト用の出力リダイレクトが可能
- 実際のOSプロセスを起動せず、単体テストレベルで検証可能（高速・信頼性）
- TDDサイクル（Red-Green-Refactor）を高速に回せる

**Implementation Pattern**:
```go
// cmd/echo_test.go
func TestEchoCommand(t *testing.T) {
    tests := []struct {
        name           string
        args           []string
        expectedStdout string
        expectedStderr string
        expectedErr    bool
    }{
        {
            name:           "basic output",
            args:           []string{"Hello", "World"},
            expectedStdout: "Hello World\n",
            expectedStderr: "",
            expectedErr:    false,
        },
        {
            name:           "suppress newline",
            args:           []string{"-n", "Hello"},
            expectedStdout: "Hello",
            expectedStderr: "",
            expectedErr:    false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var stdout, stderr bytes.Buffer
            cmd := echoCmd
            cmd.SetOut(&stdout)
            cmd.SetErr(&stderr)
            cmd.SetArgs(tt.args)

            err := cmd.Execute()
            if (err != nil) != tt.expectedErr {
                t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
            }
            if stdout.String() != tt.expectedStdout {
                t.Errorf("stdout: expected %q, got %q", tt.expectedStdout, stdout.String())
            }
            if stderr.String() != tt.expectedStderr {
                t.Errorf("stderr: expected %q, got %q", tt.expectedStderr, stderr.String())
            }
        })
    }
}
```

**Alternatives Considered**:
- **Alternative A**: `os/exec`で実際のバイナリを実行 → 却下理由: ビルドが必要で遅い、TDDに不向き
- **Alternative B**: グローバル変数でstdoutをモック → 却下理由: Cobraの標準機能で十分
- **Alternative C**: ファイルに出力してテスト後に読み込み → 却下理由: I/Oオーバーヘッド、テスト隔離が困難

**Best Practices**:
- テストケースをテーブル駆動テストで構造化（`tests := []struct{...}`）
- 各テストケースで`cmd.SetOut()`, `cmd.SetErr()`を設定し、出力を独立させる
- `t.Run()`でサブテストを実行し、並列実行可能に（`t.Parallel()`オプション）

---

## 5. UTF-8 Encoding in Go

### Decision: Goの標準文字列処理（UTF-8ネイティブ）をそのまま使用

**Research Summary**:
- Goの文字列型（`string`）はUTF-8エンコーディングを前提に設計
- `os.Args`から取得されるコマンドライン引数は、OSがUTF-8でエンコードされていれば自動的にUTF-8
- `fmt.Fprint`, `os.Stdout.Write`などはUTF-8をそのまま出力

**Rationale**:
- FR-014（UTF-8のみサポート）は、Goの標準動作で自然に満たされる
- 追加のエンコーディング変換ライブラリ（`golang.org/x/text/encoding`など）は不要
- シンプルな実装で十分な機能を提供

**Alternatives Considered**:
- **Alternative A**: マルチエンコーディング対応（Shift_JIS, EUC-JP等） → 却下理由: 複雑性増大、ユースケース不明
- **Alternative B**: `golang.org/x/text`で明示的にUTF-8検証 → 却下理由: 過剰な実装、Goのデフォルトで十分

**Best Practices**:
- 非UTF-8バイト列が入力された場合の動作は未定義として扱う（エッジケースとしてドキュメント化）
- テストケースには日本語やEmoji（マルチバイトUTF-8文字）を含める

---

## 6. Verbose Logging Strategy

### Decision: `--verbose`フラグで条件付きログをstderrに出力

**Research Summary**:
- Cobraの`PersistentFlags()`を使用すれば、すべてのサブコマンドで`--verbose`が利用可能
- `cmd.Flags().BoolP()`でローカルフラグとして定義することも可能
- 標準ライブラリの`log`パッケージは十分な機能を提供（`log.SetOutput(os.Stderr)`でstderrに出力）

**Rationale**:
- FR-013（`--verbose`フラグ）は、Cobraのフラグ機能で簡単に実装可能
- デバッグ情報はstderrに出力し、正常出力（stdout）と分離（FR-009）
- `log.SetPrefix("[VERBOSE] ")`でログの識別を容易に

**Implementation Pattern**:
```go
// cmd/echo.go
func init() {
    echoCmd.Flags().BoolP("verbose", "v", false, "Enable verbose debug output")
}

func runEcho(cmd *cobra.Command, args []string) error {
    verbose, _ := cmd.Flags().GetBool("verbose")
    if verbose {
        log.SetOutput(cmd.ErrOrStderr()) // stderrに出力
        log.SetPrefix("[VERBOSE] ")
        log.Println("Processing echo command with args:", args)
    }
    // 実装ロジック
    return nil
}
```

**Alternatives Considered**:
- **Alternative A**: `PersistentFlags()`でグローバルフラグ化 → 却下理由: echoコマンド専用で十分
- **Alternative B**: 外部ログライブラリ（`logrus`, `zap`） → 却下理由: 過剰な依存、標準ライブラリで十分

**Best Practices**:
- `--verbose`時のログは、処理された引数、フラグの状態、エスケープシーケンス変換の詳細を含める
- ログのフォーマットは読みやすく、タイムスタンプやファイル名は不要（シンプルさ優先）

---

## Research Completion Checklist

- [x] UNIX `echo`の標準動作仕様を明確化（POSIX + GNU拡張）
- [x] エスケープシーケンス処理のベストプラクティス（カスタムパーサー）
- [x] Cobraのエラーハンドリングとヘルプカスタマイズ（`SilenceUsage: false`）
- [x] TDDパターンの標準出力テスト（`bytes.Buffer` + `cmd.SetOut()`）
- [x] UTF-8エンコーディングの扱い（Goの標準動作）
- [x] Verbose loggingの実装方法（`--verbose`フラグ + `log`パッケージ）

**Phase 0 Gate**: ✅ PASSED - すべてのNEEDS CLARIFICATIONが解決され、技術的基盤が確立されました。

---

## Next Steps

Phase 1に進み、以下のアーティファクトを作成:
1. `data-model.md`: エンティティとデータ構造の定義
2. `contracts/`: API/コマンドインターフェース仕様
3. `quickstart.md`: 実装後のクイックスタートガイド
