# Research: Bats Integration Testing Framework

**Date**: 2025-11-30  
**Feature**: 003-bats-integration  
**Phase**: 0 - Research & Outline

## Executive Summary

このドキュメントは、go-cliプロジェクトに統合テストフレームワークを導入するための技術調査結果をまとめたものです。Bats-coreを選定し、テスト分離戦略、Makefile統合、CI/CD統合の具体的な実装方針を決定しました。

---

## 1. テストフレームワークの選定

### 決定: Bats-core v1.10.0+

**選定理由**:
- **業界標準**: Bash Automated Testing Systemは、CLIツールの統合テストにおけるデファクトスタンダード
- **シンプルな構文**: `@test "description" { ... }` 形式で、直感的にテストを記述可能
- **豊富なエコシステム**: bats-support、bats-assertなどの拡張ライブラリが利用可能
- **CI/CD親和性**: TAP（Test Anything Protocol）出力をサポートし、CI環境での統合が容易
- **アクティブなメンテナンス**: GitHub Actionsなどのモダンな環境で広く採用

**代替案の評価**:

| フレームワーク | 評価 | 却下理由 |
|--------------|------|---------|
| shunit2 | △ | メンテナンスが不活発、モダンなCI環境での実績が少ない |
| roundup | △ | コミュニティが小さい、ドキュメントが不足 |
| shell script only | △ | アサーションライブラリがなく、テストコードが冗長になる |
| Go testing package | × | 統合テストとしては過剰、ブラックボックステストに不向き |

**インストール方法**:

```bash
# macOS (推奨)
brew install bats-core

# Linux (Ubuntu/Debian)
sudo apt-get install bats

# 手動インストール（すべてのプラットフォーム）
git clone https://github.com/bats-core/bats-core.git
cd bats-core
sudo ./install.sh /usr/local
```

**バージョン管理**:
- 最小要件: v1.10.0
- 推奨: 最新stable版（v1.11.0以降）
- CI環境では特定バージョンを固定してインストール

---

## 2. テスト分離戦略

### 決定: 一時ディレクトリ + 環境変数オーバーライド

**実装パターン**:

```bash
# helpers/test_env.bash

setup_test_env() {
    # 一意の一時ディレクトリを作成
    export TEST_TEMP_DIR=$(mktemp -d -t mycli-test.XXXXXX)
    export TEST_CONFIG_HOME="${TEST_TEMP_DIR}/config"
    export TEST_HOME="${TEST_TEMP_DIR}/home"
    
    # ディレクトリ構造を準備
    mkdir -p "${TEST_CONFIG_HOME}/mycli"
    mkdir -p "${TEST_HOME}"
    
    # 環境変数をオーバーライド
    export MYCLI_CONFIG="${TEST_CONFIG_HOME}/mycli"
    export HOME="${TEST_HOME}"
    
    # クリーンアップのためのtrapハンドラー
    trap cleanup_test_env EXIT
}

cleanup_test_env() {
    # 一時ディレクトリを削除
    if [[ -n "${TEST_TEMP_DIR}" && -d "${TEST_TEMP_DIR}" ]]; then
        rm -rf "${TEST_TEMP_DIR}"
    fi
}

teardown_test_env() {
    cleanup_test_env
}
```

**技術的根拠**:
- **`mktemp -d`**: POSIX互換で、一意の一時ディレクトリを安全に作成
- **環境変数のオーバーライド**: アプリケーションコードを変更せず、設定パスを制御
- **`trap EXIT`**: テストが異常終了してもクリーンアップを保証
- **独立性**: 並列テスト実行時も、各プロセスが独自の一時ディレクトリを使用

**代替案の評価**:

| 戦略 | 評価 | 却下理由 |
|------|------|---------|
| バックアップ&復元 | × | レースコンディション、並列実行不可 |
| 共有テストディレクトリ | × | テスト間の干渉、クリーンアップ漏れのリスク |
| Docker コンテナ | △ | オーバーヘッドが大きい、ローカル開発での煩雑さ |
| 環境変数のみ | △ | 設定ファイルの分離ができない |

---

## 3. Makefileベースのテスト統合

### 決定: 階層的Makefile構造

**プロジェクトルート `Makefile` の拡張**:

```makefile
# 既存ターゲットに追加
.PHONY: integration-test integration-test-root integration-test-configure integration-test-echo

integration-test:
	@echo "Running integration tests..."
	@$(MAKE) -C integration_test all

integration-test-root:
	@$(MAKE) -C integration_test test-root

integration-test-configure:
	@$(MAKE) -C integration_test test-configure

integration-test-echo:
	@$(MAKE) -C integration_test test-echo
```

**`integration_test/Makefile` の新規作成**:

```makefile
BATS := bats
BATS_FLAGS := --formatter pretty
BINARY_PATH := ../bin/mycli

.PHONY: all test-root test-configure test-echo check-binary check-bats

all: check-binary check-bats
	@echo "Running all integration tests..."
	@$(BATS) $(BATS_FLAGS) root.bats configure.bats echo.bats

test-root: check-binary check-bats
	@$(BATS) $(BATS_FLAGS) root.bats

test-configure: check-binary check-bats
	@$(BATS) $(BATS_FLAGS) configure.bats

test-echo: check-binary check-bats
	@$(BATS) $(BATS_FLAGS) echo.bats

check-binary:
	@if [ ! -f "$(BINARY_PATH)" ]; then \
		echo "Error: Binary not found at $(BINARY_PATH)"; \
		echo "Please run 'make build' from project root first."; \
		exit 1; \
	fi

check-bats:
	@if ! command -v $(BATS) >/dev/null 2>&1; then \
		echo "Error: bats is not installed"; \
		echo "Please install bats-core: brew install bats-core"; \
		exit 1; \
	fi
```

**技術的根拠**:
- **階層化**: プロジェクトルートとテストディレクトリでMakefileを分離し、関心の分離を実現
- **前提条件チェック**: `check-binary`と`check-bats`で実行前に必要条件を検証
- **明確なエラーメッセージ**: ユーザーフレンドリーなガイダンスを提供
- **個別実行**: 各コマンドのテストを独立して実行可能

**代替案の評価**:

| アプローチ | 評価 | 却下理由 |
|----------|------|---------|
| 単一Makefile | △ | integration_testディレクトリの独立性が低下 |
| シェルスクリプト | △ | Makeと整合性がとれない、ビルドツールが統一されない |
| npm scripts | × | Node.js依存を追加したくない |

---

## 4. CI/CD統合のベストプラクティス

### 決定: 専用テストステージの追加

**GitHub Actions ワークフロー更新**:

```yaml
# .github/workflows/ci.yaml に追加

jobs:
  build:
    # 既存のビルドジョブ（変更なし）
    
  integration-test:
    runs-on: ubuntu-latest
    needs: build  # ビルド成功後に実行
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25.x'
      
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

**技術的根拠**:
- **専用ステージ**: ビルドと統合テストを明確に分離し、失敗箇所の特定が容易
- **needs: build**: ビルドが成功しないと統合テストは実行されない（品質ゲート）
- **if: always()**: テスト失敗時もアーティファクトをアップロードし、デバッグを容易に
- **並列実行の可能性**: 将来的にマトリックス戦略で複数OS/バージョンをテスト可能

**並列実行の検討**:

現時点では実装せず、将来の拡張として以下を検討可能：

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest]
    go-version: ['1.25.x', '1.26.x']
```

**代替案の評価**:

| アプローチ | 評価 | 却下理由 |
|----------|------|---------|
| ビルドステージに統合 | × | 失敗箇所の特定が困難、ビルドとテストの責務混在 |
| 並列実行 | △ | 現時点では過剰、必要になってから実装 |
| セルフホストランナー | △ | インフラ管理の負担増、GitHub Actionsで十分 |

---

## 5. 詳細出力モードの実装

### 決定: 環境変数 + Bats formatterオプション

**実装方法**:

```bash
# デフォルト: 簡潔な出力
make integration-test

# 詳細出力: TAP形式
BATS_FORMATTER=tap make integration-test

# 詳細出力: Pretty形式（推奨）
BATS_FORMATTER=pretty make integration-test

# さらに詳細: 個別のテストケースごとの出力
BATS_VERBOSE=1 make integration-test
```

**Makefile内での対応**:

```makefile
# integration_test/Makefile

BATS_FORMATTER ?= pretty
BATS_VERBOSE ?=

BATS_FLAGS := --formatter $(BATS_FORMATTER)
ifdef BATS_VERBOSE
    BATS_FLAGS += --verbose-run
endif

all: check-binary check-bats
	@$(BATS) $(BATS_FLAGS) root.bats configure.bats echo.bats
```

**出力形式の比較**:

| 形式 | 用途 | 詳細レベル |
|------|------|----------|
| pretty | デフォルト、ローカル開発 | 中（成功/失敗の概要） |
| tap | CI/CD、自動化 | 低（機械可読形式） |
| junit | CI/CD、統計収集 | 中（XML形式） |
| --verbose-run | デバッグ | 高（すべてのコマンド出力） |

**技術的根拠**:
- **環境変数制御**: コマンドラインから簡単に出力レベルを変更可能
- **デフォルトはpretty**: 開発者にとって最も読みやすい形式
- **CI環境での柔軟性**: TAP形式でCI/CDツールと統合可能

---

## リスクと緩和策

| リスク | 影響度 | 発生確率 | 緩和策 |
|--------|--------|---------|--------|
| Bats未インストール | 高 | 中 | `check-bats`ターゲットで事前チェック、明確なインストール手順 |
| 一時ディレクトリのクリーンアップ漏れ | 中 | 低 | `trap EXIT`で確実にクリーンアップ |
| テスト間の干渉 | 高 | 低 | 一意の一時ディレクトリで完全に分離 |
| CI環境でのBatsインストール失敗 | 中 | 低 | apt-getで確実にインストール、バージョンを固定 |
| パフォーマンス劣化（30秒超過） | 中 | 中 | テストケース数を適切に管理、並列実行を検討 |

---

## 次のステップ

Phase 0の調査結果を基に、Phase 1（Design & Contracts）に進みます：

1. **Data Model定義**: テストスイート、テストケース、ヘルパー関数の構造を定義
2. **契約作成**: 各コマンド（root, configure, echo）のテスト契約を作成
3. **Quickstart作成**: 開発者向けのクイックスタートガイドを作成

すべての技術的な疑問点は解決されており、実装に必要な情報が揃っています。
