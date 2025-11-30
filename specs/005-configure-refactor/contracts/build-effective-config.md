# Contract: BuildEffectiveConfig Function

**Feature**: Configure設定構造のリファクタリング  
**Component**: `cmd/viperutils.go` - BuildEffectiveConfig function  
**Date**: 2025-11-30

## Purpose

Define the contract for the updated BuildEffectiveConfig() function that returns a map containing the new nested configuration structure with default values.

## Function Signature

```go
func BuildEffectiveConfig() map[string]interface{}
```

### Parameters

None

### Return Value

- **Type**: `map[string]interface{}`
- **Description**: A nested map representing the complete configuration structure with default values
- **Usage**: Passed to `internal/cmd/configure.Configure()` for YAML file generation

## Output Contract

### Map Structure

The function returns a map with the following structure:

```go
map[string]interface{}{
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
```

### Key-Value Specifications

| Key Path | Type | Default Value | Description |
|----------|------|---------------|-------------|
| `client-id` | string | `""` | Client identifier (existing) |
| `client-secret` | string | `""` | Client secret (existing) |
| `common` | map | (nested) | Common settings section |
| `common.var1` | string | `""` | Common string variable |
| `common.var2` | int | `123` | Common integer variable |
| `hoge` | map | (nested) | Application-specific section |
| `hoge.fuga` | string | `"hello"` | First-level string value |
| `hoge.foo` | map | (nested) | Second-level nested section |
| `hoge.foo.bar` | string | `"hello"` | Second-level string value |

### Type Guarantees

- **Top-level keys**: Always strings
- **String values**: Always string type (not nil)
- **Integer values**: Always int type (not nil)
- **Nested maps**: Always `map[string]interface{}` type
- **No nil values**: All leaf values are concrete types

## Behavior Specification

### 1. Consistent Output

**Given**: Function is called  
**When**: No parameters passed  
**Then**: 
- Returns identical map structure every time
- All default values are consistent
- Map is newly allocated (not shared reference)

**Verification**:
```go
map1 := BuildEffectiveConfig()
map2 := BuildEffectiveConfig()
// map1 and map2 have same structure and values
// but are different map instances
```

### 2. Nested Structure Preservation

**Given**: Function returns map  
**When**: Map is marshaled to YAML  
**Then**: 
- YAML output preserves nested structure
- Indentation reflects nesting levels
- Type information is preserved (int vs string)

**Verification**:
```go
cfg := BuildEffectiveConfig()
yamlBytes, _ := yaml.Marshal(cfg)
// Result:
// client-id: ""
// client-secret: ""
// common:
//   var1: ""
//   var2: 123
// hoge:
//   fuga: "hello"
//   foo:
//     bar: "hello"
```

### 3. Independence from Global State

**Given**: CliConfig global variable has custom values  
**When**: BuildEffectiveConfig() is called  
**Then**: 
- Returned map contains hardcoded defaults, not CliConfig values
- Function does not read from CliConfig
- Function is pure (no side effects)

**Rationale**: Simplifies testing and ensures consistent default generation

## Integration Contract

### Usage in configure Command

```go
// cmd/configure.go
opts := configure.ConfigureOptions{
    // ... other options ...
    Data: BuildEffectiveConfig(),
}
return configure.ConfigureFunc(target, opts)
```

### Usage in internal/cmd/configure

```go
// internal/cmd/configure/configure.go
func Configure(target string, opts ConfigureOptions) error {
    // opts.Data is the map returned by BuildEffectiveConfig()
    out, err := yaml.Marshal(opts.Data)
    // ... write to file ...
}
```

### YAML Output

When `BuildEffectiveConfig()` output is marshaled to YAML and written to file, the result is:

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
- 10 lines total (including blank lines for readability)
- YAML v3 formatting (via gopkg.in/yaml.v3)
- Indentation: 2 spaces per level

## Error Handling

### No Errors Possible

This function:
- Takes no parameters (no invalid input)
- Returns hardcoded literals (no runtime failures)
- Performs no I/O (no file/network errors)
- Does not panic

Therefore, no error handling is required by callers.

## Testing Contract

### Required Test Cases

1. **Map Structure Verification**
   - Input: None
   - Expected: Map contains all required keys
   - Verification: Check each key exists

2. **Default Value Verification**
   - Input: None
   - Expected: All values match specification
   - Verification: Assert each value equals expected default

3. **Type Verification**
   - Input: None
   - Expected: Values have correct types
   - Verification: Type assertions succeed

4. **Nested Map Verification**
   - Input: None
   - Expected: Nested maps are accessible
   - Verification: `map["common"]["var1"]` succeeds

5. **YAML Marshaling Verification**
   - Input: None
   - Expected: Map marshals to valid YAML
   - Verification: Unmarshal result equals original structure

6. **Independence Verification**
   - Input: None
   - Expected: Multiple calls return independent maps
   - Verification: Modifying one map doesn't affect another

### Test Example

```go
func TestBuildEffectiveConfig(t *testing.T) {
    cfg := BuildEffectiveConfig()
    
    // Top-level fields
    assert.Equal(t, "", cfg["client-id"])
    assert.Equal(t, "", cfg["client-secret"])
    
    // Nested common section
    common, ok := cfg["common"].(map[string]interface{})
    assert.True(t, ok)
    assert.Equal(t, "", common["var1"])
    assert.Equal(t, 123, common["var2"])
    
    // Nested hoge section (2 levels)
    hoge, ok := cfg["hoge"].(map[string]interface{})
    assert.True(t, ok)
    assert.Equal(t, "hello", hoge["fuga"])
    
    foo, ok := hoge["foo"].(map[string]interface{})
    assert.True(t, ok)
    assert.Equal(t, "hello", foo["bar"])
}
```

## Performance Characteristics

- **Time Complexity**: O(1) - returns literal map
- **Space Complexity**: O(1) - fixed-size map (7 leaf values)
- **Allocation**: Single map allocation + nested map allocations
- **Execution Time**: < 1μs (literal return)

## Backward Compatibility

### Changes from Previous Version

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

### Compatibility Impact

1. **Signature**: Unchanged (return type still `map[string]interface{}`)
2. **Existing Keys**: `client-id` and `client-secret` still present
3. **New Keys**: Additional keys added (backward compatible)
4. **Behavior Change**: Now returns literals instead of reading CliConfig

**Breaking Change**: Callers expecting CliConfig values will now get defaults. However:
- Current usage in `cmd/configure.go` wants defaults (for scaffolding)
- No other callers exist in codebase
- Therefore, no actual breaking change in practice

## Dependencies

### Internal

- Used by: `cmd/configure.go` (RunE function)
- Depends on: None (pure function)

### External

- None (standard library only)

## Future Extensions

### Potential Enhancements

1. **Profile-Specific Defaults**: Accept profile parameter for different defaults
   ```go
   func BuildEffectiveConfig(profile string) map[string]interface{}
   ```

2. **Merge with Current Config**: Combine defaults with existing CliConfig
   ```go
   func BuildEffectiveConfig() map[string]interface{} {
       defaults := getDefaults()
       current := configToMap(CliConfig)
       return merge(defaults, current)
   }
   ```

3. **Schema-Driven Generation**: Load defaults from JSON schema
   ```go
   func BuildEffectiveConfig() map[string]interface{} {
       return loadFromSchema("config.schema.json")
   }
   ```

### Deprecation Path

If function needs to change signature:
1. Create new function with different name (e.g., `BuildConfigWithDefaults()`)
2. Mark current function as deprecated
3. Migrate callers gradually
4. Remove after 2 major versions

## References

- **Data Model**: [data-model.md](../data-model.md)
- **Config Struct Contract**: [config-struct.md](./config-struct.md)
- **Research**: [research.md](../research.md)
- **Spec**: [spec.md](../spec.md)
- **YAML v3**: https://pkg.go.dev/gopkg.in/yaml.v3
