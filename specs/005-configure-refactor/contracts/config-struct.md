# Contract: Config Struct

**Feature**: Configure設定構造のリファクタリング  
**Component**: `cmd/root.go` - Config struct and nested types  
**Date**: 2025-11-30

## Purpose

Define the contract for the expanded Config struct that supports nested configuration structure while maintaining backward compatibility with existing fields.

## Interface Definition

### Config Struct

```go
type Config struct {
    ClientID     string       `mapstructure:"client-id"`
    ClientSecret string       `mapstructure:"client-secret"`
    Common       CommonConfig `mapstructure:"common"`
    Hoge         HogeConfig   `mapstructure:"hoge"`
}
```

### CommonConfig Struct

```go
type CommonConfig struct {
    Var1 string `mapstructure:"var1"`
    Var2 int    `mapstructure:"var2"`
}
```

### HogeConfig Struct

```go
type HogeConfig struct {
    Fuga string    `mapstructure:"fuga"`
    Foo  FooConfig `mapstructure:"foo"`
}
```

### FooConfig Struct

```go
type FooConfig struct {
    Bar string `mapstructure:"bar"`
}
```

## Input Contract

### YAML Structure

The Config struct is designed to unmarshal from the following YAML structure:

```yaml
client-id: "value"
client-secret: "value"

common:
  var1: "value"
  var2: 123

hoge:
  fuga: "value"
  foo:
    bar: "value"
```

### Viper Integration

- **Unmarshal Method**: `viper.Unmarshal(&CliConfig)` populates all fields
- **Key Mapping**: mapstructure tags define YAML key → struct field mapping
- **Missing Keys**: Fields with missing YAML keys use Go zero values

## Output Contract

### Field Access

After unmarshaling, fields are accessed via dot notation:

```go
// Top-level fields
clientID := CliConfig.ClientID
clientSecret := CliConfig.ClientSecret

// Nested fields (1st level)
var1 := CliConfig.Common.Var1
var2 := CliConfig.Common.Var2
fuga := CliConfig.Hoge.Fuga

// Nested fields (2nd level)
bar := CliConfig.Hoge.Foo.Bar
```

### Type Guarantees

- `ClientID`, `ClientSecret`: Always string type
- `Common.Var1`, `Hoge.Fuga`, `Hoge.Foo.Bar`: Always string type
- `Common.Var2`: Always int type
- All fields: Non-nullable (use zero values for missing data)

## Behavior Specification

### 1. Complete Configuration

**Given**: YAML contains all fields  
**When**: `viper.Unmarshal(&CliConfig)` is called  
**Then**: 
- All fields populated with values from YAML
- No zero values (unless explicitly set in YAML)

**Example**:
```yaml
client-id: "abc123"
client-secret: "secret456"
common:
  var1: "test"
  var2: 999
hoge:
  fuga: "world"
  foo:
    bar: "baz"
```

Result:
```go
CliConfig.ClientID == "abc123"
CliConfig.ClientSecret == "secret456"
CliConfig.Common.Var1 == "test"
CliConfig.Common.Var2 == 999
CliConfig.Hoge.Fuga == "world"
CliConfig.Hoge.Foo.Bar == "baz"
```

### 2. Backward Compatibility (Old Config)

**Given**: YAML contains only old fields  
**When**: `viper.Unmarshal(&CliConfig)` is called  
**Then**: 
- Old fields populated from YAML
- New fields use zero values (empty strings, 0)

**Example**:
```yaml
client-id: "abc123"
client-secret: "secret456"
```

Result:
```go
CliConfig.ClientID == "abc123"
CliConfig.ClientSecret == "secret456"
CliConfig.Common.Var1 == ""
CliConfig.Common.Var2 == 0
CliConfig.Hoge.Fuga == ""
CliConfig.Hoge.Foo.Bar == ""
```

### 3. Partial Configuration

**Given**: YAML contains some new fields but not all  
**When**: `viper.Unmarshal(&CliConfig)` is called  
**Then**: 
- Present fields populated from YAML
- Missing fields use zero values

**Example**:
```yaml
client-id: "abc123"
common:
  var2: 456
```

Result:
```go
CliConfig.ClientID == "abc123"
CliConfig.ClientSecret == ""
CliConfig.Common.Var1 == ""
CliConfig.Common.Var2 == 456
CliConfig.Hoge.Fuga == ""
CliConfig.Hoge.Foo.Bar == ""
```

### 4. Type Mismatch Handling

**Given**: YAML contains wrong type for a field  
**When**: `viper.Unmarshal(&CliConfig)` is called  
**Then**: 
- Viper/mapstructure attempts type conversion
- If conversion fails, field uses zero value
- No error returned (silent failure per Viper behavior)

**Example**:
```yaml
common:
  var2: "not a number"
```

Result:
```go
CliConfig.Common.Var2 == 0
```

## Error Handling

### Unmarshal Errors

The Config struct itself does not validate data. Validation errors are handled by the caller:

```go
if err := viper.Unmarshal(&CliConfig); err != nil {
    // Handle error (e.g., log to stderr)
    fmt.Fprintln(os.Stderr, "Failed to parse configuration:", err)
}
```

### Common Error Cases

1. **Malformed YAML**: Viper.ReadInConfig() fails before Unmarshal
2. **Type Conversion Failure**: Silent failure with zero value (Viper default)
3. **Missing Required Fields**: Not applicable (all fields are optional)

## Testing Contract

### Required Test Cases

1. **Full Structure Unmarshal**
   - Input: Complete YAML with all fields
   - Expected: All Config fields match YAML values

2. **Old Config Compatibility**
   - Input: YAML with only `client-id` and `client-secret`
   - Expected: Old fields loaded, new fields are zero values

3. **Partial Nested Structure**
   - Input: YAML with `common` but not `hoge`
   - Expected: Common fields loaded, Hoge fields are zero values

4. **Empty YAML**
   - Input: Empty YAML file
   - Expected: All fields are zero values

5. **Nested Type Verification**
   - Input: YAML with nested structures
   - Expected: `CliConfig.Common` is CommonConfig type
   - Expected: `CliConfig.Hoge.Foo` is FooConfig type

### Test Helper

```go
func createTestConfig(yaml string) (*Config, error) {
    vp := viper.New()
    vp.SetConfigType("yaml")
    vp.ReadConfig(strings.NewReader(yaml))
    
    var cfg Config
    err := vp.Unmarshal(&cfg)
    return &cfg, err
}
```

## Dependencies

### Internal

- `cmd/viperutils.go`: Uses Config struct for BuildEffectiveConfig()
- `cmd/root.go`: Defines global CliConfig variable

### External

- `github.com/spf13/viper`: Configuration unmarshaling
- `github.com/mitchellh/mapstructure`: Underlying struct mapping (via Viper)

## Performance Characteristics

- **Unmarshal Time**: O(n) where n = number of YAML keys (~7 keys)
- **Memory**: 4 structs, ~100 bytes total (negligible)
- **Allocation**: No dynamic allocation after unmarshal (all fields are value types or nested structs)

## Backward Compatibility Guarantees

1. **Existing Fields Unchanged**: `ClientID` and `ClientSecret` signatures identical
2. **Zero Value Safety**: Missing fields safely default to zero values
3. **No Breaking Changes**: Existing code accessing old fields continues to work
4. **Viper Behavior**: Relies on Viper's stable unmarshaling behavior

## Future Extensions

### Potential Additions

1. **Validation Tags**: Add `validate:"required"` or `validate:"min=1"` tags
2. **Default Tags**: Add `default:"value"` tags for explicit defaults
3. **Environment Variable Mapping**: Explicit `envconfig` tags for nested fields
4. **Custom Unmarshal**: Implement `UnmarshalYAML()` for complex validation

### Deprecation Path

If fields need to be removed in the future:
1. Mark field as deprecated in comments
2. Add deprecation warning when field is used
3. Remove after 2 major versions

## References

- **Data Model**: [data-model.md](../data-model.md)
- **Research**: [research.md](../research.md)
- **Spec**: [spec.md](../spec.md)
- **Viper Unmarshaling**: https://github.com/spf13/viper#unmarshaling
- **mapstructure**: https://github.com/mitchellh/mapstructure
