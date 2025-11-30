# Integration Test Quickstart Guide

**Feature**: 003-bats-integration  
**Date**: 2025-11-30

## Overview

このガイドは、go-cliプロジェクトの統合テストを素早く開始するための手順を提供します。Batsを使用した統合テストフレームワークのセットアップから実行、トラブルシューティングまでをカバーします。

---

## Prerequisites

### 1. Batsのインストール

#### macOS (推奨)
```bash
brew install bats-core
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install -y bats
```

#### Linux (手動インストール)
```bash
git clone https://github.com/bats-core/bats-core.git
cd bats-core
sudo ./install.sh /usr/local
```

#### インストール確認
```bash
bats --version
# Expected output: Bats 1.10.0 (or higher)
```

### 2. アプリケーションのビルド

統合テストを実行する前に、必ずバイナリをビルドしてください：

```bash
# プロジェクトルートから
make build
```

ビルドが成功すると、`bin/mycli`が作成されます。

---

## Running Tests

### すべての統合テストを実行

プロジェクトルートから：

```bash
make integration-test
```

これにより、root、configure、echoの全コマンドのテストが実行されます。

### 個別コマンドのテストを実行

#### Rootコマンドのみ
```bash
make integration-test-root
```

#### Configureコマンドのみ
```bash
make integration-test-configure
```

#### Echoコマンドのみ
```bash
make integration-test-echo
```

### Batsを直接実行

`integration_test/`ディレクトリから：

```bash
# すべてのテスト
bats *.bats

# 特定のファイル
bats root.bats

# 特定のテストケースのみ
bats root.bats --filter "TC-ROOT-001"
```

---

## Output Modes

### デフォルト出力（Pretty形式）

```bash
make integration-test
```

出力例：
```
 ✓ TC-ROOT-001: Display help with no arguments
 ✓ TC-ROOT-002: Display help with --help flag
 ✓ TC-ROOT-003: Display help with -h flag
 ✗ TC-ROOT-004: Display version with --version flag

10 tests, 1 failure
```

### 詳細出力（Verbose形式）

```bash
BATS_VERBOSE=1 make integration-test
```

各テストケースのコマンド実行と出力が詳細に表示されます。

### TAP形式（CI/CD向け）

```bash
BATS_FORMATTER=tap make integration-test
```

機械可読なTAP（Test Anything Protocol）形式で出力されます。

### JUnit XML形式

```bash
bats --formatter junit integration_test/*.bats > test-results.xml
```

CI/CDツールでのレポート生成に使用できます。

---

## Writing New Tests

### 1. 適切なテストファイルを選択

- Rootコマンドのテスト → `integration_test/root.bats`
- Configureコマンドのテスト → `integration_test/configure.bats`
- Echoコマンドのテスト → `integration_test/echo.bats`

### 2. テストケースの基本構造

```bash
#!/usr/bin/env bats

# ヘルパーのロード
load helpers/common
load helpers/assertions
load helpers/test_env

# 各テストのセットアップ
setup() {
    setup_test_env
}

# 各テストのクリーンアップ
teardown() {
    teardown_test_env
}

# テストケース
@test "TC-XXX-NNN: Description of what is being tested" {
    # Given: 前提条件の設定
    create_test_config "default" "key: value"
    
    # When: テスト対象のアクションを実行
    run_mycli command arg1 arg2
    
    # Then: 結果を検証
    assert_success
    assert_output "expected output"
    assert_file_exists "$TEST_CONFIG_HOME/mycli/config.yaml"
}
```

### 3. 利用可能なヘルパー関数

#### `helpers/common.bash`

```bash
setup_test_env          # テスト環境を初期化
teardown_test_env       # テスト環境をクリーンアップ
run_mycli <args>        # mycliを実行（結果は$status, $output, $linesに格納）
create_test_config <profile> <content>  # テスト用設定ファイルを作成
```

#### `helpers/assertions.bash`

```bash
assert_success          # 終了ステータスが0であることを確認
assert_failure          # 終了ステータスが非0であることを確認
assert_output <expected>  # 出力が期待値と一致することを確認
assert_output_contains <substring>  # 出力に文字列が含まれることを確認
assert_line <index> <expected>  # 特定行の内容を確認
assert_file_exists <path>  # ファイルが存在することを確認
```

#### `helpers/test_env.bash`

```bash
mock_editor <behavior>  # エディタを模擬（"save", "cancel", "error"）
set_test_profile <name>  # テスト用のプロファイルを設定
```

### 4. テストケースID命名規則

```
TC-<COMMAND>-<NUMBER>: <Description>

例:
TC-ROOT-001: Display help with no arguments
TC-ECHO-015: Large output handling
TC-CONF-003: Create profile-specific config
```

---

## Test Isolation

### 自動分離

各テストは独立した一時ディレクトリで実行されます：

- `$TEST_TEMP_DIR`: 一意の一時ディレクトリ
- `$TEST_CONFIG_HOME`: テスト用の設定ホーム
- `$TEST_HOME`: テスト用のHOMEディレクトリ

これらは自動的にクリーンアップされます。

### 手動での環境変数オーバーライド

```bash
@test "Custom environment test" {
    setup_test_env
    export MYCLI_CONFIG="$TEST_CONFIG_HOME/custom"
    
    run_mycli --help
    assert_success
}
```

---

## Debugging Tests

### 特定のテストケースのみ実行

```bash
bats integration_test/root.bats --filter "TC-ROOT-001"
```

### テストの出力を確認

```bash
@test "Debug test" {
    run_mycli echo "test"
    
    # デバッグ情報を出力
    echo "Status: $status"
    echo "Output: $output"
    echo "Lines: ${lines[@]}"
    
    assert_success
}
```

### インタラクティブなデバッグ

テスト内で`set -x`を使用すると、実行されるコマンドが詳細に表示されます：

```bash
@test "Verbose debugging" {
    set -x
    setup_test_env
    run_mycli echo "test"
    assert_success
    set +x
}
```

---

## Troubleshooting

### エラー: "Binary not found at bin/mycli"

**原因**: バイナリがビルドされていません。

**解決策**:
```bash
make build
```

### エラー: "bats: command not found"

**原因**: Batsがインストールされていません。

**解決策**:
```bash
# macOS
brew install bats-core

# Linux
sudo apt-get install bats
```

### エラー: "Permission denied"

**原因**: バイナリに実行権限がありません。

**解決策**:
```bash
chmod +x bin/mycli
```

### テストが失敗する: "assert_output failed"

**原因**: 期待される出力と実際の出力が一致しません。

**デバッグ手順**:
1. `BATS_VERBOSE=1`で実行して詳細を確認
2. テスト内で`echo "$output"`を追加して実際の出力を確認
3. 改行文字や空白文字の違いに注意

### テストが非常に遅い

**原因**: 一時ディレクトリのクリーンアップに時間がかかっている可能性があります。

**解決策**:
- 個別のテストファイルを実行してボトルネックを特定
- 不要なファイル生成を減らす
- 並列実行を検討（将来の機能）

---

## CI/CD Integration

### GitHub Actionsの例

```yaml
- name: Install Bats
  run: |
    sudo apt-get update
    sudo apt-get install -y bats

- name: Build binary
  run: make build

- name: Run integration tests
  run: make integration-test

- name: Upload test results
  if: always()
  uses: actions/upload-artifact@v3
  with:
    name: integration-test-results
    path: integration_test/*.tap
```

---

## Best Practices

### ✅ DO

- テストは独立して実行可能に保つ
- `setup_test_env`と`teardown_test_env`を常に使用する
- 明確で説明的なテストケースIDと説明を使用する
- エッジケースもテストする
- CI/CDでテストを実行する

### ❌ DON'T

- 実際のユーザー設定ファイルを変更しない
- テスト間で状態を共有しない
- ハードコードされたパスを使用しない（環境変数を使用）
- 長時間実行されるテストを書かない（目標: 全体で30秒以内）
- `sleep`で待機しない（適切な完了チェックを使用）

---

## Next Steps

1. **既存のテストを実行**: `make integration-test`
2. **テスト契約を確認**: `specs/003-bats-integration/contracts/`ディレクトリ
3. **新しいテストを追加**: 上記のテンプレートを使用
4. **CI/CDに統合**: `.github/workflows/ci.yaml`を更新

---

## Resources

- **Bats Documentation**: https://bats-core.readthedocs.io/
- **Test Contracts**: `specs/003-bats-integration/contracts/`
- **Data Model**: `specs/003-bats-integration/data-model.md`
- **Research Document**: `specs/003-bats-integration/research.md`

---

## Support

問題が発生した場合は、以下を確認してください：

1. `make build`が成功していること
2. `bats --version`が正常に表示されること
3. `bin/mycli --version`が実行できること
4. テスト環境が正しくセットアップされていること

それでも問題が解決しない場合は、プロジェクトのIssueトラッカーで報告してください。
