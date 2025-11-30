# Data Model: Integration Test Suite Structure

**Date**: 2025-11-30  
**Feature**: 003-bats-integration  
**Phase**: 1 - Design & Contracts

## Overview

このドキュメントは、Bats統合テストフレームワークの構造とデータモデルを定義します。テストスイート、テストケース、ヘルパー関数、テスト環境の各エンティティとその関係性を明確にします。

---

## 1. Test Suite Structure

### Entity: TestSuite

```yaml
TestSuite:
  name: string                    # テストスイート名 ("root", "configure", "echo")
  file: string                    # ファイル名 ("root.bats", "configure.bats", "echo.bats")
  description: string             # スイートの説明
  helpers: [string]               # ロードするヘルパースクリプト
  setup_file: function            # スイート全体の初期化（オプション）
  teardown_file: function         # スイート全体のクリーンアップ（オプション）
  tests: [TestCase]               # テストケースの配列
```

**Example**:

```yaml
TestSuite:
  name: "echo"
  file: "echo.bats"
  description: "Integration tests for echo command"
  helpers:
    - "helpers/common.bash"
    - "helpers/assertions.bash"
  tests:
    - TC-ECHO-001
    - TC-ECHO-002
    - TC-ECHO-003
```

---

## 2. Test Case Structure

### Entity: TestCase

```yaml
TestCase:
  id: string                      # 一意のテストID ("TC-ROOT-001")
  description: string             # テストケースの説明
  priority: enum                  # "P1" | "P2" | "P3"
  given: string                   # 前提条件
  when: string                    # 実行するアクション
  then: string                    # 期待される結果
  setup: function                 # 個別テストの初期化（オプション）
  teardown: function              # 個別テストのクリーンアップ（オプション）
  assertions: [Assertion]         # アサーションの配列
  skip: boolean                   # テストをスキップするか（デフォルト: false）
  skip_reason: string             # スキップの理由（skipがtrueの場合）
```

**Example**:

```yaml
TestCase:
  id: "TC-ECHO-001"
  description: "Basic output with single argument"
  priority: "P1"
  given: "Single string argument is provided"
  when: "mycli echo 'Hello' is executed"
  then: "'Hello' followed by newline is output, exit status 0"
  assertions:
    - type: "status"
      expected: 0
    - type: "output"
      expected: "Hello\n"
```

---

## 3. Assertion Structure

### Entity: Assertion

```yaml
Assertion:
  type: enum                      # アサーションのタイプ
  expected: any                   # 期待される値
  actual: any                     # 実際の値（実行時に計算）
  message: string                 # 失敗時のカスタムメッセージ（オプション）
```

**Assertion Types**:

| Type | Description | Expected Format | Example |
|------|-------------|----------------|---------|
| `status` | 終了ステータスコードの検証 | integer | `0`, `1` |
| `output` | 標準出力全体の検証 | string | `"Hello\n"` |
| `output_contains` | 標準出力に特定文字列が含まれるか | string | `"Error:"` |
| `output_regex` | 標準出力が正規表現に一致するか | regex | `"^mycli version [0-9]"` |
| `line` | 特定行の内容を検証 | `{index: int, text: string}` | `{index: 0, text: "Hello"}` |
| `line_count` | 出力行数の検証 | integer | `3` |
| `file_exists` | ファイルの存在確認 | path | `"$TEST_CONFIG_HOME/mycli/default.yaml"` |
| `file_contains` | ファイル内容の検証 | `{path: string, text: string}` | `{path: "config.yaml", text: "key: value"}` |
| `env_var` | 環境変数の値を検証 | `{name: string, value: string}` | `{name: "MYCLI_CONFIG", value: "/tmp/..."}` |

**Example Assertions**:

```bash
# Bats形式での実装例
@test "TC-ECHO-001: Basic output" {
    run "$MYCLI_BINARY" echo "Hello"
    
    # Assertion: status == 0
    assert_success
    
    # Assertion: output == "Hello\n"
    assert_output "Hello"
}

@test "TC-ECHO-002: Multiple arguments" {
    run "$MYCLI_BINARY" echo Hello World
    
    # Assertion: status == 0
    assert_success
    
    # Assertion: output contains "Hello World"
    assert_output "Hello World"
}
```

---

## 4. Test Environment Structure

### Entity: TestEnvironment

```yaml
TestEnvironment:
  temp_dir: string                # 一時ディレクトリのパス
  config_home: string             # 設定ファイルのホームディレクトリ
  home_dir: string                # テスト用のHOMEディレクトリ
  binary_path: string             # テスト対象バイナリのパス
  env_vars: map[string]string     # オーバーライドする環境変数
  cleanup_required: boolean       # クリーンアップが必要か
  isolation_level: enum           # "full" | "partial" | "none"
```

**Environment Setup Flow**:

```
1. Create unique temp directory
   ├─> $TEST_TEMP_DIR = mktemp -d
   └─> Ensures no conflicts between test runs

2. Create subdirectories
   ├─> $TEST_CONFIG_HOME = $TEST_TEMP_DIR/config
   ├─> $TEST_HOME = $TEST_TEMP_DIR/home
   └─> mkdir -p both directories

3. Override environment variables
   ├─> MYCLI_CONFIG=$TEST_CONFIG_HOME/mycli
   ├─> HOME=$TEST_HOME
   └─> PATH includes project bin/

4. Register cleanup handler
   └─> trap cleanup_test_env EXIT
```

**Example**:

```yaml
TestEnvironment:
  temp_dir: "/tmp/mycli-test.abc123"
  config_home: "/tmp/mycli-test.abc123/config"
  home_dir: "/tmp/mycli-test.abc123/home"
  binary_path: "../bin/mycli"
  env_vars:
    MYCLI_CONFIG: "/tmp/mycli-test.abc123/config/mycli"
    HOME: "/tmp/mycli-test.abc123/home"
  cleanup_required: true
  isolation_level: "full"
```

---

## 5. Helper Function Structure

### Entity: HelperFunction

```yaml
HelperFunction:
  name: string                    # 関数名
  description: string             # 機能の説明
  parameters: [Parameter]         # パラメータの配列
  returns: ReturnValue            # 戻り値
  side_effects: [string]          # 副作用（環境変数設定など）
  example: string                 # 使用例
```

### Common Helper Functions

#### `helpers/common.bash`

```yaml
CommonHelpers:
  - name: setup_test_env
    description: テスト環境を初期化（一時ディレクトリ作成、環境変数設定）
    parameters: []
    returns:
      type: void
    side_effects:
      - "Sets $TEST_TEMP_DIR"
      - "Sets $TEST_CONFIG_HOME"
      - "Sets $TEST_HOME"
      - "Exports MYCLI_CONFIG"
      - "Exports HOME"
      - "Registers cleanup trap"
    example: |
      @test "Example test" {
        setup_test_env
        # test logic here
        teardown_test_env
      }

  - name: teardown_test_env
    description: テスト環境をクリーンアップ（一時ディレクトリ削除）
    parameters: []
    returns:
      type: void
    side_effects:
      - "Removes $TEST_TEMP_DIR"
    example: |
      teardown_test_env

  - name: run_mycli
    description: mycliバイナリを実行し、結果を$statusと$outputに格納
    parameters:
      - name: args
        type: variadic string
        description: mycliに渡す引数
    returns:
      type: void
    side_effects:
      - "Sets $status (exit code)"
      - "Sets $output (stdout+stderr)"
      - "Sets $lines (output as array)"
    example: |
      run_mycli echo "Hello"
      assert_success

  - name: create_test_config
    description: テスト用の設定ファイルを作成
    parameters:
      - name: profile
        type: string
        description: プロファイル名（デフォルト: "default"）
      - name: content
        type: string
        description: YAMLコンテンツ
    returns:
      type: string
      description: 作成された設定ファイルのパス
    example: |
      config_path=$(create_test_config "default" "key: value")
```

#### `helpers/assertions.bash`

```yaml
AssertionHelpers:
  - name: assert_success
    description: 終了ステータスが0であることを確認
    parameters: []
    returns:
      type: void
    example: |
      run_mycli echo "test"
      assert_success

  - name: assert_failure
    description: 終了ステータスが非0であることを確認
    parameters: []
    returns:
      type: void
    example: |
      run_mycli --invalid-flag
      assert_failure

  - name: assert_output
    description: 出力が期待値と一致することを確認
    parameters:
      - name: expected
        type: string
        description: 期待される出力
    returns:
      type: void
    example: |
      run_mycli echo "Hello"
      assert_output "Hello"

  - name: assert_output_contains
    description: 出力に特定文字列が含まれることを確認
    parameters:
      - name: substring
        type: string
        description: 含まれるべき文字列
    returns:
      type: void
    example: |
      run_mycli --help
      assert_output_contains "Usage:"

  - name: assert_line
    description: 特定行が期待値と一致することを確認
    parameters:
      - name: index
        type: integer
        description: 行番号（0-indexed）
      - name: expected
        type: string
        description: 期待される行内容
    returns:
      type: void
    example: |
      run_mycli echo -e "Line1\nLine2"
      assert_line 0 "Line1"
      assert_line 1 "Line2"

  - name: assert_file_exists
    description: ファイルが存在することを確認
    parameters:
      - name: path
        type: string
        description: ファイルパス
    returns:
      type: void
    example: |
      assert_file_exists "$TEST_CONFIG_HOME/mycli/default.yaml"
```

#### `helpers/test_env.bash`

```yaml
EnvironmentHelpers:
  - name: mock_editor
    description: エディタを模擬（configureコマンドのテスト用）
    parameters:
      - name: behavior
        type: string
        description: "save" | "cancel" | "error"
    returns:
      type: void
    side_effects:
      - "Exports EDITOR with mock script path"
    example: |
      mock_editor "save"
      run_mycli configure

  - name: set_test_profile
    description: テスト用のプロファイルを設定
    parameters:
      - name: profile
        type: string
        description: プロファイル名
    returns:
      type: void
    side_effects:
      - "Exports MYCLI_PROFILE"
    example: |
      set_test_profile "test"
      run_mycli configure
```

---

## 6. Test Suite Organization

### Directory Structure

```
integration_test/
├── Makefile                     # テスト実行の自動化
├── README.md                    # 統合テストのドキュメント
├── helpers/                     # 共有ヘルパースクリプト
│   ├── common.bash              # 基本的なテストユーティリティ
│   ├── assertions.bash          # カスタムアサーション関数
│   └── test_env.bash            # 環境設定ユーティリティ
├── root.bats                    # Rootコマンドのテスト
├── configure.bats               # Configureコマンドのテスト
└── echo.bats                    # Echoコマンドのテスト
```

### Test File Template

```bash
#!/usr/bin/env bats

# テストファイル: {command}.bats
# 説明: {Command}コマンドの統合テスト

# ヘルパーのロード
load helpers/common
load helpers/assertions
load helpers/test_env

# スイート全体のセットアップ（オプション）
setup_file() {
    # 全テストで共通の初期化処理
    export MYCLI_BINARY="../bin/mycli"
}

# 各テストのセットアップ
setup() {
    setup_test_env
}

# 各テストのクリーンアップ
teardown() {
    teardown_test_env
}

# テストケース
@test "TC-{CMD}-001: {Description}" {
    # Given: {precondition}
    # When: {action}
    # Then: {expected outcome}
    
    run_mycli {args}
    assert_success
    assert_output "{expected}"
}

# 追加のテストケース...
```

---

## 7. Relationships and Dependencies

### Entity Relationship Diagram

```
TestSuite (1) ────> (*) TestCase
                        │
                        ├──> (1) TestEnvironment
                        └──> (*) Assertion

TestCase (*) ────> (*) HelperFunction

HelperFunction (collection)
    ├── CommonHelpers
    ├── AssertionHelpers
    └── EnvironmentHelpers
```

### Dependency Flow

```
Test Execution Flow:
1. Makefile invokes bats
2. Bats loads test file (.bats)
3. Test file loads helpers
4. setup_file() runs (once per suite)
5. For each test:
   a. setup() runs
   b. Test body executes
   c. Assertions validate
   d. teardown() runs
6. teardown_file() runs (once per suite)
7. Cleanup handlers execute (trap EXIT)
```

---

## 8. Test Data Management

### Test Fixtures

テストフィクスチャは、`helpers/fixtures/`ディレクトリに配置することを推奨（将来の拡張）:

```yaml
Fixture:
  type: enum                      # "config" | "input" | "expected_output"
  name: string                    # フィクスチャ名
  path: string                    # ファイルパス
  content: string | binary        # コンテンツ
  format: string                  # "yaml" | "json" | "text" | "binary"
```

**Example Usage**:

```bash
# ヘルパー関数でフィクスチャをロード
load_fixture() {
    local fixture_name="$1"
    cat "helpers/fixtures/${fixture_name}"
}

@test "TC-CONF-001: Load config from fixture" {
    setup_test_env
    local config_content=$(load_fixture "valid_config.yaml")
    echo "$config_content" > "$TEST_CONFIG_HOME/mycli/default.yaml"
    
    run_mycli some-command
    assert_success
}
```

---

## Summary

このデータモデルは、Bats統合テストフレームワークの全体的な構造を定義しています：

- **TestSuite**: テストの最上位グループ（コマンド単位）
- **TestCase**: 個別のテストシナリオ
- **Assertion**: 期待値と実際の値の検証
- **TestEnvironment**: 分離されたテスト実行環境
- **HelperFunction**: 再利用可能なテストユーティリティ

これらのエンティティは、明確な責務と関係性を持ち、保守性の高いテストスイートを構築するための基盤となります。
