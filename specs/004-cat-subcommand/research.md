# Research & Technical Decisions: Cat サブコマンド

**Feature**: Cat サブコマンド実装  
**Date**: 2025-11-30  
**Status**: Complete

## Research Summary

このドキュメントは、仕様書の"NEEDS CLARIFICATION"項目を解決し、技術選択の根拠を記録する。

### 解決された不明点

仕様作成時には以下の技術的詳細が明確化セッションで解決されました：

1. **行番号フォーマットの桁数**: 6桁固定（999,999行超過時も6桁に収め、最下位6桁のみ表示）
2. **制御文字変換の具体的範囲**: ASCII 0-31（タブ9・改行10除く）+ ASCII 127（DEL）
3. **ストリーム処理のバッファサイズ**: 32KB（標準的、バランス重視）
4. **複数エラー時の終了コード判定**: 1つでもエラーなら終了コード1（部分的成功も失敗）
5. **行番号が999,999超過時の動作**: 最下位6桁のみ表示

すべての明確化項目が解決済みのため、追加のリサーチは不要です。

## Technology Choices

### 1. ファイル処理とストリーム処理

**Decision**: `bufio.Scanner`を使用した行単位ストリーム処理

**Rationale**:
- Go標準ライブラリの`bufio.Scanner`は行区切りのストリーム処理に最適
- デフォルトのバッファサイズ（64KB）から32KBに調整することで、メモリ効率と性能のバランスを確保
- `Scanner.Buffer()`メソッドで明示的にバッファサイズを設定可能
- 行番号付加やフォーマット処理を行単位で実行できるため、巨大ファイルでもメモリ使用量を抑制

**Alternatives Considered**:
- `io.Copy()`: バイト単位のコピーは行番号付加などの行単位処理に不向き
- `ioutil.ReadFile()`: ファイル全体をメモリに読み込むため、巨大ファイルで不適切
- `bufio.Reader.ReadString('\n')`: 手動でエラーハンドリングが必要で冗長

### 2. 行番号フォーマット

**Decision**: `fmt.Sprintf("%6d  ", lineNum % 1000000)`による6桁固定フォーマット

**Rationale**:
- 明確化セッションで決定された6桁固定フォーマットに準拠
- `%6d`で右揃え6桁を保証
- 999,999を超えた場合は`% 1000000`で最下位6桁のみ取得（1,000,000 → 0、1,234,567 → 234,567）
- フォーマット文字列は事前定義された定数として保持し、性能を最適化

**Alternatives Considered**:
- 動的幅調整: 実装が複雑で、出力フォーマットの一貫性が損なわれる
- エラー時の処理中断: 実用上、百万行を超えるファイルも処理できる必要がある

### 3. 制御文字変換

**Decision**: ルックアップテーブルによる高速変換

**Rationale**:
- ASCII 0-31（タブ9・改行10除く）と127（DEL）の変換マップを事前構築
- マップキー: `byte`（制御文字のASCII値）
- マップ値: `string`（変換後の文字列、例: 7 → "^G"、127 → "^?"）
- `strings.Builder`で効率的な文字列結合を実現

**Implementation**:
```go
var controlCharMap = map[byte]string{
    0: "^@", 1: "^A", 2: "^B", 3: "^C", 4: "^D", 5: "^E", 6: "^F", 7: "^G",
    8: "^H", /* 9: tab - skip */, /* 10: newline - skip */, 11: "^K", 12: "^L",
    13: "^M", 14: "^N", 15: "^O", 16: "^P", 17: "^Q", 18: "^R", 19: "^S",
    20: "^T", 21: "^U", 22: "^V", 23: "^W", 24: "^X", 25: "^Y", 26: "^Z",
    27: "^[", 28: "^\\", 29: "^]", 30: "^^", 31: "^_",
    127: "^?",
}
```

**Alternatives Considered**:
- 都度計算: `if c < 32 { return "^" + string('A' + c - 1) }` は可読性が低く、バグの温床
- 正規表現: オーバーヘッドが大きく、バイト単位処理には不向き

### 4. オプション処理

**Decision**: Cobraの`Flags()`とブールフラグを使用

**Rationale**:
- Cobraの標準的なフラグ機構を活用
- `-n`, `-b`, `-E`, `-T`, `-v`, `-A`を個別のブールフラグとして定義
- `-A`フラグが指定された場合、内部で`-v`, `-E`, `-T`を自動的に有効化
- `-n`と`-b`の競合は、後に指定された方を優先（POSIXの標準動作に準拠）

**Implementation**:
```go
cmd.Flags().BoolP("number", "n", false, "number all output lines")
cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")
```

**Alternatives Considered**:
- 手動パース: Cobraの利点を活かせず、バグの原因になりやすい
- カスタムフラグ型: 過度に複雑で保守性が低下

### 5. エラーハンドリング

**Decision**: 継続処理と終了コード管理

**Rationale**:
- 明確化セッションで決定: 1つでもエラーがあれば最終的に終了コード1
- ファイルごとにエラーをキャッチし、エラーメッセージをstderrに出力
- エラーが発生しても次のファイルの処理を継続（UNIXのcatと同じ動作）
- `hadError`フラグで最終的な終了コードを判定

**Implementation Pattern**:
```go
hadError := false
for _, filename := range args {
    if err := processFile(filename, opts); err != nil {
        fmt.Fprintf(os.Stderr, "cat: %s: %v\n", filename, err)
        hadError = true
        continue
    }
}
if hadError {
    os.Exit(1)
}
```

**Alternatives Considered**:
- 最初のエラーで即座に終了: ユーザビリティが低く、複数ファイル処理時に不便
- エラーを無視: データ完全性の観点から不適切

### 6. テスト戦略

**Decision**: 3層テストアプローチ

**Rationale**:
- **単体テスト** (`*_test.go`): 各関数の独立したロジックを検証（TDD必須原則に準拠）
- **統合テスト** (`cmd/cat_test.go`): Cobraコマンドのフラグ解析とエンドツーエンドフローを検証
- **BATS統合テスト** (`integration_test/cat.bats`): 実際のファイルI/Oと複数ファイル処理を検証

**Coverage Goals**:
- 単体テスト: 内部ロジックの100%カバレッジ
- 統合テスト: すべてのユーザーストーリーの受入シナリオをカバー
- BATS: エッジケースと実際のファイルシステム操作を検証

**Test Utilities**:
- `t.TempDir()`: テストごとに一時ディレクトリを作成し、テスト分離を保証
- `t.Cleanup()`: グローバル変数のモック後に元の状態を復元
- `bytes.Buffer`: 標準出力/エラー出力のキャプチャに使用

## Best Practices Applied

### Go標準慣習

1. **エラーハンドリング**: すべてのエラーを明示的にチェックし、適切に処理
2. **パッケージ構成**: `cmd/`（CLI層）と`internal/cmd/cat/`（ロジック層）の明確な分離
3. **命名規則**: Go標準の命名規則に従う（例: `processFile`, `formatLineNumber`）
4. **ドキュメント**: すべてのエクスポート関数にGoDocコメントを付与

### 憲章準拠

1. **TDD必須**: テストを先に書き、Red-Green-Refactorサイクルを厳守
2. **パッケージ責務分離**: CLI層とロジック層を明確に分離し、Cobraへの依存を最小化
3. **コード品質基準**: `gofmt -s`と`golangci-lint run --enable=govet`をパス
4. **ユーザーエクスペリエンス**: 明確なヘルプメッセージ、適切なエラーメッセージ、UNIX標準との互換性

### パフォーマンス最適化

1. **ストリーム処理**: 32KBバッファで巨大ファイルをメモリ効率的に処理
2. **事前計算**: フォーマット文字列や変換マップを定数/変数として事前定義
3. **不要なアロケーション回避**: `strings.Builder`で効率的な文字列結合

## Implementation Checkpoints

- [x] すべての"NEEDS CLARIFICATION"項目を解決
- [x] 技術選択の根拠を文書化
- [x] ベストプラクティスの適用方針を明確化
- [x] テスト戦略の定義完了

## Next Steps

Phase 1（Design & Contracts）に進み、以下を作成：
1. `data-model.md`: エンティティと状態遷移の詳細設計
2. `contracts/`: 関数シグネチャとインターフェース定義
3. `quickstart.md`: 開発者向けクイックスタートガイド
