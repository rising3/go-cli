# go-cli

CobraとViperを使用したGo CLI アプリケーションのテンプレートです。

## 開発原則

このプロジェクトの開発原則と品質基準は [`.specify/memory/constitution.md`](.specify/memory/constitution.md) で定義されています。すべての実装とレビューは憲章の原則に従う必要があります。

## クイックスタート

```bash
# 依存関係のダウンロード
go mod download

# ビルドと全チェックの実行
make all

# バイナリの実行
./bin/mycli
```

## コマンド

### echo サブコマンド

UNIX互換のechoコマンド実装。テキストを標準出力に表示します。

```bash
# 基本的な出力
./bin/mycli echo "Hello, World!"

# 複数の引数（スペース区切り）
./bin/mycli echo Hello World

# 改行を抑制 (-n フラグ)
./bin/mycli echo -n "Prompt: "

# エスケープシーケンスの解釈 (-e フラグ)
./bin/mycli echo -e "Line1\nLine2\tTab"

# フラグの組み合わせ
./bin/mycli echo -n -e "No newline\twith tab"

# UTF-8サポート
./bin/mycli echo "こんにちは世界 🚀✨"

# デバッグモード
./bin/mycli echo --verbose "debug output"
```

**サポートされているエスケープシーケンス** (-e フラグ):
- `\n` - 改行
- `\t` - タブ
- `\\` - バックスラッシュ
- `\"` - ダブルクォート
- `\a` - アラート（ベル）
- `\b` - バックスペース
- `\c` - 以降の出力を抑制
- `\r` - キャリッジリターン
- `\v` - 垂直タブ

詳細な開発ガイドは [`.github/copilot-instructions.md`](.github/copilot-instructions.md) を参照してください。