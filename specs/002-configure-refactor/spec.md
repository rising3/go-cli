# Feature Specification: Configure サブコマンドのリファクタリング

**Feature Branch**: `002-configure-refactor`  
**Created**: 2025-11-30  
**Status**: Draft  
**Input**: User description: "echoサブコマンドの実装を参考にconfigureサブコマンドをリファクタリングする。internal/sdtio/sdtio.goを利用せず、cobraが提供する標準出力を利用し、ベストプラクティスに沿った実装に変更する。"

## Clarifications

### Session 2025-11-30

- Q: Configure関数のシグネチャとtargetパラメータの扱い → A: Option B - `Configure(target string, opts ConfigureOptions)`で、targetを別パラメータにする（既存実装と同じ）
- Q: エディタ起動時のストリーム管理 → A: Option A - `os.Stdin`, `os.Stdout`, `os.Stderr`を直接使用（エディタは常に実際の端末と対話）
- Q: エディタ検出失敗時のエラーハンドリング → A: Option A - エラーを吸収（stderrにメッセージ出力後、nilを返す）- 既存動作を維持
- Q: ファイル書き込み時のパーミッション指定 → A: Option A - 0644 (rw-r--r--) を維持 - 所有者のみ書き込み可、全員読み取り可
- Q: テスト用の関数変数パターン → A: Option A - `var ConfigureFunc = Configure`パターンを採用（echo同様、テストでモック可能）

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Cobra標準のI/Oストリームを使用した実装 (Priority: P1)

開発者として、configureサブコマンドの実装がCobraのベストプラクティスに従っていることを確認したい。`cmd.OutOrStdout()`, `cmd.ErrOrStderr()`, `cmd.InOrStdin()` などのCobraが提供する標準ストリームAPIを使用することで、テストの容易性と一貫性が向上する。

**Why this priority**: 現在の実装は独自の`internal/stdio`パッケージを使用しているが、Cobraフレームワークが提供する標準的なストリーム管理機能を使用することで、フレームワークとの統合性が高まり、テストも簡単になる。echoサブコマンドで既に実装されているパターンを踏襲することで、コードベース全体の一貫性も向上する。

**Independent Test**: configureコマンドの単体テストで、`cmd.SetOut()`, `cmd.SetErr()`, `cmd.SetIn()`を使用してストリームをモックし、出力が正しく行われることを検証可能。

**Acceptance Scenarios**:

1. **Given** configureコマンドが実装されている、**When** `cmd.OutOrStdout()`を使用して標準出力にメッセージを書き込む、**Then** テストで`cmd.SetOut()`によって設定されたバッファに出力される
2. **Given** configureコマンドが実装されている、**When** `cmd.ErrOrStderr()`を使用してエラーメッセージを書き込む、**Then** テストで`cmd.SetErr()`によって設定されたバッファに出力される
3. **Given** echoサブコマンドの実装パターン、**When** configureサブコマンドを同じパターンでリファクタリング、**Then** 両サブコマンドが一貫したストリーム管理を行う

---

### User Story 2 - internal/cmd/configure/パッケージの作成とロジックの分離 (Priority: P2)

開発者として、コマンドのエントリーポイント（`cmd/configure.go`）とビジネスロジック（`internal/cmd/configure/configure.go`）を明確に分離したい。echoサブコマンドと同様に、`internal/cmd/configure/`パッケージを作成して、再利用可能なコア機能を実装する。

**Why this priority**: ロジックの分離により、テストが容易になり、将来的に他のコマンドから同じ機能を再利用できる。echoサブコマンドで採用されている構造をconfigureサブコマンドにも適用することで、プロジェクト全体の構造が統一される。

**Independent Test**: `internal/cmd/configure/configure.go`の関数を、Cobraコマンドに依存せずに単体テストで直接テスト可能。

**Acceptance Scenarios**:

1. **Given** `internal/cmd/configure/`パッケージが作成されている、**When** `Configure(opts ConfigureOptions)`関数を実装、**Then** コマンドラインフラグに依存しないロジックとして機能する
2. **Given** `cmd/configure.go`がエントリーポイント、**When** フラグの値を取得して`ConfigureOptions`を構築、**Then** `internal/cmd/configure.Configure()`を呼び出してビジネスロジックを実行
3. **Given** `ConfigureOptions`構造体、**When** `Output`, `ErrOutput`フィールドに`io.Writer`を設定、**Then** 任意のストリーム（バッファ、ファイル、標準出力）に出力可能

---

### User Story 3 - ConfigureOptions構造体の設計 (Priority: P3)

開発者として、configureコマンドのすべてのオプションを構造化された形で管理したい。echoサブコマンドの`EchoOptions`と同様に、`ConfigureOptions`構造体を定義して、フラグの値、I/Oストリーム、その他の設定を集約する。

**Why this priority**: 構造体を使用することで、関数のシグネチャがシンプルになり、新しいオプションの追加が容易になる。テスト時にもオプションの組み合わせを簡単に構築できる。

**Independent Test**: 様々な`ConfigureOptions`の組み合わせを作成して、`Configure()`関数の動作を検証可能。

**Acceptance Scenarios**:

1. **Given** `ConfigureOptions`構造体が定義されている、**When** `Force`, `Edit`, `NoWait`などのフラグフィールドを含む、**Then** コマンドの動作を制御する設定として機能する
2. **Given** `ConfigureOptions`構造体、**When** `Output`, `ErrOutput`フィールドを含む、**Then** 出力先をテストで自由に変更可能
3. **Given** `ConfigureOptions`構造体、**When** `EditorLookup`, `EditorShouldWait`などの関数フィールドを含む、**Then** エディタの動作をテストでモック可能

---

### User Story 4 - internal/stdioパッケージへの依存を削除 (Priority: P4)

開発者として、configureサブコマンドが独自の`internal/stdio`パッケージに依存しないようにしたい。Cobraが提供する`io.Writer`ベースのインターフェースを使用することで、標準的なGoのI/Oパターンに従う。

**Why this priority**: `internal/stdio`パッケージは独自の抽象化レイヤーだが、Cobraと標準ライブラリだけで十分に対応可能。不要な抽象化を削除することで、コードがシンプルになり、保守性が向上する。echoサブコマンドでは既に`internal/stdio`を使用していないため、統一性のためにもconfigureサブコマンドから削除すべき。

**Independent Test**: `internal/stdio`パッケージをインポートせずに、configureコマンドが正常に動作することをテストで確認可能。

**Acceptance Scenarios**:

1. **Given** リファクタリング後のconfigureコマンド、**When** `internal/stdio`パッケージへの参照を検索、**Then** `cmd/configure.go`および`internal/cmd/configure/`に参照が存在しない
2. **Given** ファイル出力が必要な場合、**When** 標準ライブラリの`os.OpenFile()`や`os.Create()`を直接使用、**Then** `stdio.OpenWriter()`などのラッパー関数を使用しない
3. **Given** エディタ起動時のストリーム管理、**When** `cmd.Stdin`, `cmd.Stdout`, `cmd.Stderr`に直接`io.Writer`/`io.Reader`を設定、**Then** `stdio.BindCommand()`などのヘルパー関数を使用しない

---

### User Story 5 - テストの改善とカバレッジの向上 (Priority: P5)

開発者として、configureサブコマンドのテストカバレッジを向上させ、echoサブコマンドと同レベルの品質を確保したい。リファクタリングによってロジックが分離されることで、より細かい粒度でのテストが可能になる。

**Why this priority**: リファクタリングの主な目的の一つはテスト容易性の向上。ロジックが分離され、依存関係が明確になることで、単体テストが書きやすくなり、エッジケースも網羅しやすくなる。

**Independent Test**: リファクタリング後の各関数に対して個別の単体テストを作成し、`make test`でパスすることを確認。

**Acceptance Scenarios**:

1. **Given** リファクタリング後のコード、**When** `go test -cover ./cmd/... ./internal/cmd/configure/...`を実行、**Then** カバレッジが80%以上になる
2. **Given** `internal/cmd/configure/configure.go`の`Configure()`関数、**When** 様々な`ConfigureOptions`の組み合わせでテスト、**Then** すべてのエッジケース（ファイル存在、強制上書き、エディタ起動など）がカバーされる
3. **Given** `cmd/configure.go`のRunE関数、**When** Cobraのテストヘルパー（`ExecuteC()`など）を使用、**Then** コマンドライン引数とフラグの組み合わせが正しく処理されることを検証

---

### Edge Cases

- **設定ファイルが既に存在する場合**: `--force`フラグなしでは上書きせず、メッセージを表示（現在の動作を維持）
- **エディタが見つからない場合**: エラーメッセージをstderrに出力し、エディタ起動をスキップして、nilを返す（設定ファイル作成は成功しているため、コマンド全体を失敗させない - 現在の動作を維持）
- **ファイル作成権限がない場合**: 標準的な`os.OpenFile()`のエラーをそのまま返す
- **`--edit`と`--no-wait`を同時に指定**: エディタをバックグラウンドで起動し、すぐに終了（現在の動作を維持）
- **出力先が指定されていない場合**: デフォルトの設定ファイルパスを使用（現在の動作を維持）
- **プロファイル指定時**: `--profile`フラグで指定されたプロファイル名に基づいて設定ファイル名を決定（現在の動作を維持）

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `cmd/configure.go`は、Cobraの`cmd.OutOrStdout()`と`cmd.ErrOrStderr()`を使用してメッセージを出力しなければならない
- **FR-002**: `cmd/configure.go`の`RunE`関数は、フラグの値を取得して`ConfigureOptions`構造体を構築し、targetパスを決定してから、`internal/cmd/configure.Configure(target, opts)`を呼び出さなければならない
- **FR-003**: `internal/cmd/configure/`パッケージに、再利用可能な`Configure(target string, opts ConfigureOptions) error`関数を実装しなければならない（targetパラメータは別に指定）
- **FR-004**: `ConfigureOptions`構造体は、少なくとも以下のフィールドを含まなければならない:
  - `Force bool`: 既存の設定ファイルを強制上書き
  - `Edit bool`: 設定ファイル作成後にエディタで開く
  - `NoWait bool`: エディタをバックグラウンドで起動
  - `Data map[string]interface{}`: 設定データ
  - `Format string`: 出力フォーマット（yaml/json）
  - `Output io.Writer`: 標準出力用ストリーム
  - `ErrOutput io.Writer`: エラー出力用ストリーム
  - `EditorLookup func() (string, []string, error)`: エディタ検索関数
  - `EditorShouldWait func(string, []string) bool`: エディタ待機判定関数
- **FR-005**: `internal/stdio`パッケージへの依存を完全に削除しなければならない（`cmd/configure.go`および`internal/cmd/configure/`から）
- **FR-006**: 設定ファイルの作成は、`os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)`を使用して、パーミッション0644 (rw-r--r--) を明示的に指定しなければならない（既存の動作を維持）
- **FR-007**: エディタ起動時のストリーム管理は、`exec.Cmd`の`Stdin`, `Stdout`, `Stderr`フィールドに`os.Stdin`, `os.Stdout`, `os.Stderr`を直接設定しなければならない（エディタは実際の端末と対話する必要があるため、Cobraのストリームではなく、OS標準ストリームを使用）
- **FR-008**: 既存の機能（設定ファイル作成、強制上書き、エディタ起動、プロファイル対応など）は、リファクタリング後も同じように動作しなければならない
- **FR-009**: `cmd/configure.go`の`RunE`関数は、echoサブコマンドの`RunE`関数と同様のパターン（フラグ取得→オプション構築→内部関数呼び出し）に従わなければならない
- **FR-010**: `internal/cmd/configure/configure.go`に、テスト用の関数変数`var ConfigureFunc = Configure`を定義し、echoサブコマンドの`EchoFunc`と同様のパターンでテストでのモック化を可能にしなければならない
- **FR-011**: すべての既存テスト（`cmd/configure_test.go`, `cmd/configure_wrapper_test.go`, `internal/cmd/configure_test.go`）は、リファクタリング後も引き続きパスしなければならない（必要に応じてテストを更新）
- **FR-012**: リファクタリングによって、新しいテストケースを追加し、エッジケースのカバレッジを向上させなければならない

### Key Entities *(include if feature involves data)*

- **ConfigureOptions**: configureコマンドの実行オプションを集約する構造体
  - フラグ値（Force, Edit, NoWait）
  - 設定データ（Data, Format）
  - I/Oストリーム（Output, ErrOutput）
  - 外部依存関数（EditorLookup, EditorShouldWait）

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: リファクタリング後、`grep -r "internal/stdio" cmd/configure.go internal/cmd/configure/`がゼロ件の結果を返す
- **SC-002**: `cmd/configure.go`の`RunE`関数内で、`fmt.Fprintln()`や`fmt.Fprintf()`の呼び出しが0件であり、すべての出力が`cmd.Print*()`, `cmd.OutOrStdout()`, `cmd.ErrOrStderr()`を経由している
- **SC-003**: `internal/cmd/configure/configure.go`に、`Configure(target string, opts ConfigureOptions) error`関数が実装され、`*cobra.Command`型への依存がゼロである
- **SC-004**: リファクタリング前に存在したすべてのテストケース（`cmd/configure_test.go`, `cmd/configure_wrapper_test.go`, `internal/cmd/configure_test.go`）が、`make test`で100%パスする
- **SC-005**: リファクタリング後、`go test -cover ./cmd/... ./internal/cmd/configure/...`でテストカバレッジが80%以上になる
- **SC-006**: リファクタリング前後で、`mycli configure --force && cat ~/.config/mycli/default.yaml`の出力が同一であり、すべてのフラグ（`--force`, `--edit`, `--no-wait`, `--profile`）が同じ動作をする
- **SC-007**: `cmd/echo.go`と`cmd/configure.go`の`RunE`関数の構造（フラグ取得→オプション構築→内部関数呼び出しのパターン）が同じ行数±3行以内で一致している
- **SC-008**: `make lint`がゼロ件の警告とゼロ件のエラーで終了する（終了コード0）
