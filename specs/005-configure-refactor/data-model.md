# Data Model: Configure設定構造のリファクタリング

**Feature**: Configure設定構造のリファクタリング  
**Branch**: `005-configure-refactor`  
**Date**: 2025-11-30

## Overview

このドキュメントは、ネストされた設定構造をサポートするためのConfig構造体とその関連エンティティのデータモデルを定義します。

## Core Entities

### 1. Config (Root Configuration)

**Purpose**: アプリケーションの全体設定を表現するルート構造体

**Location**: `cmd/root.go`

**Structure**:
```go
type Config struct {
    ClientID     string       `mapstructure:"client-id"`
    ClientSecret string       `mapstructure:"client-secret"`
    Common       CommonConfig `mapstructure:"common"`
    Hoge         HogeConfig   `mapstructure:"hoge"`
}
```

**Fields**:
- `ClientID` (string): クライアント識別子（既存フィールド、維持）
- `ClientSecret` (string): クライアント秘密鍵（既存フィールド、維持）
- `Common` (CommonConfig): 共通設定セクション（新規追加）
- `Hoge` (HogeConfig): アプリケーション固有設定セクション（新規追加）

**Relationships**:
- Contains one `CommonConfig` instance
- Contains one `HogeConfig` instance

**Validation Rules**:
- None (すべてのフィールドはオプショナル)
- 欠けているフィールドはGoのゼロ値にフォールバック

**State Transitions**: N/A (immutable after loading)

---

### 2. CommonConfig

**Purpose**: 共通設定項目をグループ化

**Location**: `cmd/root.go` (Config構造体と同じファイル)

**Structure**:
```go
type CommonConfig struct {
    Var1 string `mapstructure:"var1"`
    Var2 int    `mapstructure:"var2"`
}
```

**Fields**:
- `Var1` (string): 文字列型の汎用設定値（デフォルト: ""）
- `Var2` (int): 整数型の汎用設定値（デフォルト: 123）

**Relationships**:
- Embedded in `Config` struct

**Validation Rules**:
- `Var1`: 任意の文字列（バリデーションなし）
- `Var2`: 任意の整数（バリデーションなし）

**State Transitions**: N/A

---

### 3. HogeConfig

**Purpose**: アプリケーション固有の設定セクション（2階層のネストをデモ）

**Location**: `cmd/root.go`

**Structure**:
```go
type HogeConfig struct {
    Fuga string    `mapstructure:"fuga"`
    Foo  FooConfig `mapstructure:"foo"`
}
```

**Fields**:
- `Fuga` (string): 1階層目の文字列設定値（デフォルト: "hello"）
- `Foo` (FooConfig): さらにネストされた設定セクション

**Relationships**:
- Embedded in `Config` struct
- Contains one `FooConfig` instance

**Validation Rules**:
- `Fuga`: 任意の文字列（バリデーションなし）

**State Transitions**: N/A

---

### 4. FooConfig

**Purpose**: 2階層のネスト設定を表現

**Location**: `cmd/root.go`

**Structure**:
```go
type FooConfig struct {
    Bar string `mapstructure:"bar"`
}
```

**Fields**:
- `Bar` (string): 2階層目の文字列設定値（デフォルト: "hello"）

**Relationships**:
- Embedded in `HogeConfig` struct

**Validation Rules**:
- `Bar`: 任意の文字列（バリデーションなし）

**State Transitions**: N/A

---

## Data Flow

### Configuration Loading (Viper → Config)

```
1. YAML File (~/.config/mycli/default.yaml)
   ↓ Viper.ReadInConfig()
2. Viper Internal Map
   ↓ Viper.Unmarshal(&CliConfig)
3. Config Struct Instance (CliConfig)
   - ClientID: ""
   - ClientSecret: ""
   - Common:
       Var1: ""
       Var2: 123
   - Hoge:
       Fuga: "hello"
       Foo:
           Bar: "hello"
```

### Configuration Generation (BuildEffectiveConfig → YAML)

```
1. BuildEffectiveConfig() returns map[string]interface{}
   {
     "client-id": "",
     "client-secret": "",
     "common": {
       "var1": "",
       "var2": 123
     },
     "hoge": {
       "fuga": "hello",
       "foo": {
         "bar": "hello"
       }
     }
   }
   ↓ yaml.Marshal()
2. YAML Bytes
   ↓ os.WriteFile()
3. YAML File (~/.config/mycli/default.yaml)
```

## File Structure

### Generated YAML File

**Path**: `~/.config/mycli/default.yaml` (or `<profile>.yaml`)

**Format**:
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

**Characteristics**:
- Total: 10 lines (including empty lines for readability)
- Top-level: 2 fields (`client-id`, `client-secret`)
- 1st-level nest: 2 sections (`common`, `hoge`)
- 2nd-level nest: 1 section (`hoge.foo`)
- Mixed types: strings and integers

## Implementation Notes

### Struct Definition Pattern

**Location**: `cmd/root.go` (after existing `Config` struct)

```go
type Config struct {
    ClientID     string       `mapstructure:"client-id"`
    ClientSecret string       `mapstructure:"client-secret"`
    Common       CommonConfig `mapstructure:"common"`
    Hoge         HogeConfig   `mapstructure:"hoge"`
}

type CommonConfig struct {
    Var1 string `mapstructure:"var1"`
    Var2 int    `mapstructure:"var2"`
}

type HogeConfig struct {
    Fuga string    `mapstructure:"fuga"`
    Foo  FooConfig `mapstructure:"foo"`
}

type FooConfig struct {
    Bar string `mapstructure:"bar"`
}
```

**Key Decisions**:
- All nested structs defined in the same file for cohesion
- mapstructure tags use kebab-case (consistent with existing fields)
- Field names use PascalCase (Go exported field convention)
- No validation tags (validation not in scope for this feature)

### BuildEffectiveConfig Update

**Location**: `cmd/viperutils.go`

**Before**:
```go
func BuildEffectiveConfig() map[string]interface{} {
    return map[string]interface{}{
        "client-id":     CliConfig.ClientID,
        "client-secret": CliConfig.ClientSecret,
    }
}
```

**After**:
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

**Rationale**:
- Returns literal default values instead of reading from CliConfig
- Easier to understand and maintain
- No dependency on global state during config generation

## Testing Strategy

### Unit Tests

**File**: `cmd/viperutils_test.go`

**Test Cases**:
1. `TestBuildEffectiveConfig_HasAllFields`
   - Verify returned map contains all expected keys
   - Check nested structure exists (common, hoge, hoge.foo)

2. `TestBuildEffectiveConfig_CorrectTypes`
   - Verify common.var2 is int (123)
   - Verify all string fields are strings

3. `TestBuildEffectiveConfig_CorrectDefaultValues`
   - Verify client-id is ""
   - Verify common.var2 is 123
   - Verify hoge.fuga is "hello"
   - Verify hoge.foo.bar is "hello"

### Integration Tests

**File**: `cmd/root_test.go` (or new test file)

**Test Cases**:
1. `TestConfigUnmarshal_NewStructure`
   - Create YAML with new structure
   - Unmarshal into Config struct
   - Verify all fields populated correctly

2. `TestConfigUnmarshal_BackwardCompatibility`
   - Create YAML with only old fields (client-id, client-secret)
   - Unmarshal into Config struct
   - Verify old fields loaded, new fields are zero values

3. `TestConfigUnmarshal_PartialStructure`
   - Create YAML with only common section
   - Unmarshal into Config struct
   - Verify common fields loaded, hoge fields are zero values

### End-to-End Tests

**File**: `cmd/configure_test.go`

**Test Cases**:
1. `TestConfigureCommand_GeneratesNewStructure`
   - Execute `mycli configure --force`
   - Read generated YAML file
   - Parse and verify structure matches expected output

## Constraints

### Backward Compatibility

- **Existing Fields**: `ClientID` and `ClientSecret` must remain unchanged
- **Old Config Files**: New code must load old config files (missing fields → zero values)
- **Viper Behavior**: Relies on Viper's built-in handling of missing keys

### Performance

- **Struct Size**: 4 structs, ~7 fields total (negligible memory impact)
- **Unmarshal Time**: No measurable impact (same Viper algorithm)
- **File Size**: ~10 lines of YAML (well under performance thresholds)

### Type Safety

- **Compile-Time Checks**: All fields have explicit types
- **mapstructure Tags**: Ensure correct YAML key → struct field mapping
- **No Interface{}**: Avoid type assertions in application code

## Dependencies

### Internal

- `cmd/root.go`: Config struct definition
- `cmd/viperutils.go`: BuildEffectiveConfig() implementation
- Viper's `Unmarshal()`: Automatic struct population

### External

- `github.com/spf13/viper`: Configuration management
- `gopkg.in/yaml.v3`: YAML marshaling (used by internal/cmd/configure)

## Future Considerations

### Potential Extensions

1. **Validation Tags**: Add struct tags for field validation (e.g., `validate:"required"`)
2. **Environment Variable Overrides**: Explicit mapping for nested fields (e.g., `MYCLI_COMMON_VAR2`)
3. **Config Schema**: Generate JSON schema for documentation/validation
4. **Migration Tool**: Helper to upgrade old config files to new structure

### Refactoring Opportunities

1. **Type Safety**: Consider using typed constructors instead of map literals in BuildEffectiveConfig()
2. **Defaults Package**: Move default values to a separate package for centralized management
3. **Config Validation**: Add a Validate() method to Config struct

## References

- **Spec**: [spec.md](./spec.md) - Feature requirements
- **Research**: [research.md](./research.md) - Technology decisions
- **Viper Docs**: https://github.com/spf13/viper#unmarshaling
- **Constitution**: `.specify/memory/constitution.md` - Quality standards
