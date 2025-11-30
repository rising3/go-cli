# Feature Specification: Configure設定構造のリファクタリング

**Feature Branch**: `005-configure-refactor`  
**Created**: 2025-11-30  
**Status**: Draft  
**Input**: User description: "mycliの設定ファイルの変更に伴いconfigureサブコマンドをリファクタリングする。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - ネストされた設定構造のサポート (Priority: P1)

開発者として、複雑な階層構造を持つ設定ファイルを生成できるようにしたい。現在の設定構造はフラットな`client-id`と`client-secret`のみだが、新しい設定では`common`や`hoge.foo`といったネストされた設定項目をサポートする必要がある。

**Why this priority**: これは今回のリファクタリングの核心部分であり、設定ファイルの構造を拡張することで、より複雑なアプリケーション設定をサポートできるようになる。この変更なしでは、他のすべての機能も意味をなさない。

**Independent Test**: `mycli configure --force`を実行し、生成された設定ファイルが期待通りの階層構造を持つことをYAMLパーサーで検証可能。

**Acceptance Scenarios**:

1. **Given** 新しいConfig構造体が定義されている、**When** `mycli configure --force`を実行、**Then** 生成されたYAMLファイルが`client-id`と`client-secret`をトップレベルに含む
2. **Given** 新しいConfig構造体が定義されている、**When** `mycli configure --force`を実行、**Then** 生成されたYAMLファイルが`common.var1`（文字列）と`common.var2`（数値）を含む
3. **Given** 新しいConfig構造体が定義されている、**When** `mycli configure --force`を実行、**Then** 生成されたYAMLファイルが`hoge.fuga`（文字列）と`hoge.foo.bar`（文字列）の2階層のネストを含む
4. **Given** 生成された設定ファイル、**When** Viper経由で読み込み、**Then** すべてのネストされた値が正しくConfig構造体にマッピングされる

---

### User Story 2 - Config構造体のフィールド定義とmapstructureタグの更新 (Priority: P2)

開発者として、Config構造体が新しい設定項目を正確に表現し、Viperのマッピングが正しく機能することを確認したい。既存のフィールドは維持しつつ、新しいネストされた構造を追加する必要がある。

**Why this priority**: Config構造体は設定ファイルとアプリケーションの間のインターフェースであり、正しく定義されていないと、設定の読み書きが失敗する。P1の実装を正しく行うための基盤。

**Independent Test**: Config構造体のインスタンスを作成し、Viperの`Unmarshal()`を使用して設定ファイルから値を読み込み、各フィールドの値が期待通りであることを検証可能。

**Acceptance Scenarios**:

1. **Given** Config構造体、**When** `ClientID`と`ClientSecret`フィールドを確認、**Then** 既存の`mapstructure:"client-id"`と`mapstructure:"client-secret"`タグが維持されている
2. **Given** Config構造体、**When** 新しい`Common`フィールドを確認、**Then** ネストされた構造体型で、`Var1 string`と`Var2 int`フィールドを持ち、`mapstructure:"common"`タグが設定されている
3. **Given** Config構造体、**When** 新しい`Hoge`フィールドを確認、**Then** ネストされた構造体型で、`Fuga string`と`Foo`（さらにネストされた構造体）フィールドを持ち、`mapstructure:"hoge"`タグが設定されている
4. **Given** Config構造体のネストされた構造体、**When** すべてのフィールドを確認、**Then** 適切なmapstructureタグ（kebab-case）が設定されている

---

### User Story 3 - デフォルト値の設定 (Priority: P3)

開発者として、生成された設定ファイルに適切なデフォルト値が含まれることを確認したい。空文字列、サンプル文字列、数値などの初期値を設定することで、ユーザーが設定ファイルの使い方を理解しやすくなる。

**Why this priority**: デフォルト値があることで、ユーザーは設定ファイルの構造を理解しやすくなり、どのような値を設定すべきかの例を見ることができる。必須ではないが、UXを大きく向上させる。

**Independent Test**: 生成された設定ファイルを読み込み、各フィールドの初期値が期待通りであることを検証可能。

**Acceptance Scenarios**:

1. **Given** `BuildEffectiveConfig()`関数、**When** 新しいConfigインスタンスを生成、**Then** `common.var1`が空文字列`""`に初期化されている
2. **Given** `BuildEffectiveConfig()`関数、**When** 新しいConfigインスタンスを生成、**Then** `common.var2`が整数`123`に初期化されている
3. **Given** `BuildEffectiveConfig()`関数、**When** 新しいConfigインスタンスを生成、**Then** `hoge.fuga`が文字列`"hello"`に初期化されている
4. **Given** `BuildEffectiveConfig()`関数、**When** 新しいConfigインスタンスを生成、**Then** `hoge.foo.bar`が文字列`"hello"`に初期化されている

---

### User Story 4 - 既存機能との互換性維持 (Priority: P4)

開発者として、既存のconfigureコマンドの機能（`--force`, `--edit`, `--profile`など）が引き続き正常に動作することを確認したい。設定構造の変更は内部実装の変更であり、ユーザーインターフェースは変わらない。

**Why this priority**: 既存のユーザーワークフローを壊さないことは重要だが、新しい設定構造のサポート（P1-P3）が正しく実装されていれば、既存機能は自動的に動作するはず。

**Independent Test**: 既存のテストケースを実行し、すべてがパスすることを確認。また、手動で各フラグの組み合わせをテスト。

**Acceptance Scenarios**:

1. **Given** リファクタリング後のconfigureコマンド、**When** `mycli configure`を実行、**Then** 既存の設定ファイルが存在する場合は上書きせずにメッセージを表示
2. **Given** リファクタリング後のconfigureコマンド、**When** `mycli configure --force`を実行、**Then** 既存の設定ファイルを上書きして新しい構造の設定ファイルを生成
3. **Given** リファクタリング後のconfigureコマンド、**When** `mycli configure --profile dev --force`を実行、**Then** `~/.config/mycli/dev.yaml`に新しい構造の設定ファイルを生成
4. **Given** リファクタリング後のconfigureコマンド、**When** `mycli configure --edit --force`を実行、**Then** 設定ファイル生成後にエディタが起動

---

### Edge Cases

- **既存の古い構造の設定ファイルが存在する場合**: 新しいバージョンのコマンドで読み込んだときに、存在しないフィールドは無視され、既存のフィールド（`client-id`, `client-secret`）のみが読み込まれる（Viperのデフォルト動作）
- **設定ファイルに不完全なネスト構造がある場合**: 例えば`common.var1`のみ設定されて`common.var2`が欠けている場合、欠けているフィールドはGoのゼロ値（0, "", nilなど）になる
- **YAMLファイルの手動編集でインデントが不正な場合**: Viperの読み込み時にYAMLパースエラーが発生し、エラーメッセージが表示される（既存の動作と同じ）
- **数値フィールドに文字列が設定された場合**: Viperのアンマーシャル時に型変換エラーが発生する可能性があるが、これは設定ファイルの誤りとして扱われる

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `cmd/root.go`の`Config`構造体を更新し、以下の新しいフィールドを追加しなければならない:
  - `Common`構造体（`Var1 string`, `Var2 int`フィールドを含む）
  - `Hoge`構造体（`Fuga string`, `Foo`構造体フィールドを含む）
  - `Foo`構造体（`Bar string`フィールドを含む）
- **FR-002**: すべてのネストされた構造体とそのフィールドに、適切な`mapstructure`タグを設定しなければならない（kebab-case形式）
- **FR-003**: 既存の`ClientID`と`ClientSecret`フィールドは削除せずに維持しなければならない
- **FR-004**: `cmd/viperutils.go`（または適切な場所）の`BuildEffectiveConfig()`関数を更新し、新しい設定構造のデフォルト値を返すようにしなければならない:
  - `common.var1`: `""`（空文字列）
  - `common.var2`: `123`
  - `hoge.fuga`: `"hello"`
  - `hoge.foo.bar`: `"hello"`
- **FR-005**: `mycli configure --force`を実行したときに、生成されるYAMLファイルが以下の構造を持たなければならない:
  ```yaml
  client-id: ""
  client-secret: ""
  
  common:
    var1: ""
    var2: 123
  
  hoge:
    fuga: "hello"
    foo:
      bar: "hello"
  ```
- **FR-006**: Viperの`Unmarshal()`が、新しいYAML構造を正しくConfig構造体にマッピングしなければならない
- **FR-007**: 既存のconfigureコマンドのすべてのフラグ（`--force`, `--edit`, `--no-wait`, `--profile`）が、リファクタリング後も同じように動作しなければならない
- **FR-008**: 既存のテスト（`cmd/configure_test.go`, `cmd/configure_wrapper_test.go`など）が、必要に応じて更新され、リファクタリング後も100%パスしなければならない
- **FR-009**: `internal/cmd/configure/configure.go`の`Configure()`関数は、`opts.Data`を受け取って、それをYAML形式でファイルに書き込む既存のロジックを変更する必要はない（構造体からmap[string]interface{}への変換は呼び出し側が行う）

### Key Entities

- **Config**: アプリケーションの設定を表現する構造体
  - トップレベルフィールド: `ClientID`, `ClientSecret`
  - ネストされた構造体: `Common`, `Hoge`
  - 2階層のネスト: `Hoge.Foo`
- **CommonConfig**: `common`セクションの設定を表現する構造体
  - `Var1`: 文字列型のサンプル設定値
  - `Var2`: 整数型のサンプル設定値
- **HogeConfig**: `hoge`セクションの設定を表現する構造体
  - `Fuga`: 文字列型の設定値
  - `Foo`: さらにネストされた構造体
- **FooConfig**: `hoge.foo`セクションの設定を表現する構造体
  - `Bar`: 文字列型の設定値

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `mycli configure --force`を実行後、生成された`~/.config/mycli/default.yaml`ファイルが、指定された7行の構造（トップレベル2項目 + `common`セクション3行 + `hoge`セクション4行）を持つ
- **SC-002**: 生成された設定ファイルをYAMLパーサー（`gopkg.in/yaml.v3`）で読み込み、7つのすべてのフィールド（`client-id`, `client-secret`, `common.var1`, `common.var2`, `hoge.fuga`, `hoge.foo.bar`）が期待通りの値を持つことを検証できる
- **SC-003**: `viper.Unmarshal(&CliConfig)`を実行後、`CliConfig.Common.Var2`が整数`123`であり、`CliConfig.Hoge.Foo.Bar`が文字列`"hello"`であることを検証できる
- **SC-004**: リファクタリング前に存在したすべてのテストケースが、`make test`で100%パスする
- **SC-005**: `make lint`がゼロ件の警告とゼロ件のエラーで終了する（終了コード0）
- **SC-006**: `mycli configure --profile prod --force`を実行したときに、`~/.config/mycli/prod.yaml`に新しい構造の設定ファイルが生成され、`default.yaml`と同じ構造を持つ
