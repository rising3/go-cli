# Research: Configure設定構造のリファクタリング

**Feature**: Configure設定構造のリファクタリング  
**Branch**: `005-configure-refactor`  
**Date**: 2025-11-30

## Phase 0: Research Findings

### Overview

このフィーチャーは、既存のConfig構造体を拡張してネストされた設定項目をサポートします。すべての技術的決定は既存のコードベースとViperのベストプラクティスに基づいています。

### Research Tasks Completed

#### 1. Viper Nested Structure Mapping

**Question**: Viperでネストされた構造体をどのようにマッピングするか？

**Decision**: mapstructureタグを使用した構造体ネストがViperで自動的にサポートされる

**Rationale**:
- Viperの`Unmarshal()`は、mapstructureタグを持つネストされた構造体を自動的に処理
- YAMLの階層構造が構造体のフィールドに自然にマッピングされる
- 既存のコード（`ClientID`, `ClientSecret`）が既にこのパターンを使用している

**Alternatives considered**:
- 手動でmap[string]interface{}を走査してフィールドに設定 → mapstructure自動処理の方が保守性が高い
- フラット構造を維持してドット記法を使用（例: "hoge.foo.bar"をキーとする） → 型安全性が失われる

**References**:
- 既存コード: `cmd/root.go` の `Config` 構造体
- Viper documentation: https://github.com/spf13/viper (Unmarshaling)

#### 2. Struct Field Naming Conventions

**Question**: 構造体名とフィールド名の命名規則は？

**Decision**: 
- 構造体名: PascalCase (例: `CommonConfig`, `HogeConfig`, `FooConfig`)
- フィールド名: PascalCase (例: `Var1`, `Var2`, `Fuga`)
- mapstructureタグ: kebab-case (例: `mapstructure:"var1"`, `mapstructure:"var2"`)

**Rationale**:
- Go標準の命名規約に従う（exported名はPascalCase）
- mapstructureタグはYAML/JSONのキー名に対応（kebab-caseは設定ファイルで一般的）
- 既存のフィールド（`ClientID` → `mapstructure:"client-id"`）と一貫性を保つ

**Alternatives considered**:
- すべてsnake_caseに統一 → Go標準から外れる、既存コードとの不整合
- mapstructureタグなしでフィールド名をそのまま使用 → ケース変換の制御が困難

**References**:
- Go Code Review Comments: https://go.dev/wiki/CodeReviewComments#naming
- 既存コード: `cmd/root.go` の `ClientID string \`mapstructure:"client-id"\``

#### 3. Default Values Strategy

**Question**: デフォルト値はどこで設定するか？

**Decision**: `BuildEffectiveConfig()`関数内で構造化されたmapを返す

**Rationale**:
- 既存の`BuildEffectiveConfig()`が既にこのパターンを実装している
- `internal/cmd/configure/configure.go`は`map[string]interface{}`を受け取るため、変更不要
- YAMLシリアライズ時にネストされた構造が自動的に保持される

**Implementation**:
```go
func BuildEffectiveConfig() map[string]interface{} {
    return map[string]interface{}{
        "client-id":     "",
        "client-secret": "",
        "common": map[string]interface{}{
            "var1": "",
            "var2": 123,
        },
        "hoge": map[string]interface{}{
            "fuga": "hello",
            "foo": map[string]interface{}{
                "bar": "hello",
            },
        },
    }
}
```

**Alternatives considered**:
- Config構造体インスタンスを作成してマーシャル → 型情報が利用できるが、既存関数の戻り値型を変更する必要がある
- 各サブコマンドで個別にデフォルト値を設定 → DRY原則に違反、一貫性が失われる

**References**:
- 既存コード: `cmd/viperutils.go` の `BuildEffectiveConfig()`

#### 4. Backward Compatibility

**Question**: 既存の設定ファイルとの後方互換性は？

**Decision**: 新しいフィールドは追加のみで、既存フィールドは削除・変更しない

**Rationale**:
- Viperは存在しないフィールドを無視するため、古い設定ファイルも読み込み可能
- 新しいフィールドが欠けている場合、Goのゼロ値（"", 0など）にフォールバック
- FR-003で既存フィールドの維持を明示的に要求

**Verification**:
- 既存テストが引き続きパスすることで後方互換性を確認
- Edge Caseで古い構造の設定ファイルの動作を文書化

**Alternatives considered**:
- マイグレーションスクリプトを提供 → シンプルな追加機能には過剰
- バージョンフィールドを追加 → この段階では不要（将来的に検討可能）

**References**:
- Spec: Edge Cases section - "既存の古い構造の設定ファイルが存在する場合"

#### 5. Testing Strategy

**Question**: 新しい構造をどのようにテストするか？

**Decision**: 
1. 単体テスト: `BuildEffectiveConfig()`がマップに正しい構造を返すことを検証
2. 統合テスト: Viperの`Unmarshal()`で新しいConfig構造体に正しくマッピングされることを検証
3. 既存テスト: すべてパスすることを確認（後方互換性の証明）

**Test Cases**:
- `BuildEffectiveConfig()` returns correct nested map structure
- YAML marshaling produces expected file structure
- Viper unmarshaling populates all Config struct fields correctly
- Old config files can still be loaded (missing fields use zero values)

**Alternatives considered**:
- エンドツーエンドテストのみ → フィードバックループが遅い、問題の特定が困難
- モックを使わない実ファイル操作 → テストが遅い、環境依存のリスク

**References**:
- 憲章: "Core Principles > I. テスト駆動開発 (TDD)"
- 既存テスト: `cmd/viperutils_test.go`, `cmd/configure_test.go`

### Technology Stack (Confirmed)

すべての技術スタックは既存のプロジェクトから継承され、新しい依存関係は不要:

| Component | Version | Purpose | Status |
|-----------|---------|---------|--------|
| Go | 1.25.4 | Language | ✅ Confirmed |
| Viper | v1.21.0+ | Config management | ✅ Confirmed |
| gopkg.in/yaml.v3 | v3.0.1 | YAML serialization | ✅ Confirmed |
| Cobra | v1.10.1+ | CLI framework | ✅ Confirmed (minimal usage) |

### Implementation Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-----------|
| Viperのアンマーシャル失敗 | Low | High | 統合テストでViperのUnmarshal動作を検証 |
| 既存テストの破損 | Low | Medium | テストファースト、既存テストを先に実行 |
| YAML構造の誤り | Medium | Low | YAMLパーサーでの検証、手動確認 |
| パフォーマンス劣化 | Very Low | Low | ベンチマークテスト（必要に応じて） |

### Open Questions

**None** - すべての技術的疑問は解決済み。Phase 1（デザイン）に進行可能。

### References

1. **Viper Documentation**: https://github.com/spf13/viper
   - Unmarshaling configuration into structs
   - Nested key support

2. **Go Struct Tags**: https://go.dev/wiki/Well-known-struct-tags
   - mapstructure tag format and usage

3. **Existing Implementation**:
   - `cmd/root.go` - Config struct pattern
   - `cmd/viperutils.go` - BuildEffectiveConfig() pattern
   - `internal/cmd/configure/configure.go` - map[string]interface{} consumption

4. **Project Constitution**: `.specify/memory/constitution.md`
   - TDD principles
   - Package separation guidelines
   - Quality standards
