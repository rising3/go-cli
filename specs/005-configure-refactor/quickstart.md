# Quickstart: Configure設定構造のリファクタリング

**Feature**: Configure設定構造のリファクタリング  
**Branch**: `005-configure-refactor`  
**Date**: 2025-11-30

## Overview

このドキュメントは、Config構造体の拡張を実装するための段階的な手順を提供します。テストファースト開発(TDD)のアプローチに従い、各ステップで検証可能な成果物を作成します。

## Prerequisites

- Go 1.25.4がインストールされている
- golangci-lint v2.6.2がインストールされている
- `$(go env GOPATH)/bin`がPATHに含まれている
- 作業ブランチ `005-configure-refactor` にチェックアウト済み

**確認コマンド**:
```bash
go version  # Go 1.25.4を確認
golangci-lint --version  # v2.6.2を確認
git branch --show-current  # 005-configure-refactorを確認
```

## Implementation Phases

### Phase 1: Config Struct Definition (10 minutes)

#### Step 1.1: Define Nested Structs

**File**: `cmd/root.go`

**Action**: Add new struct definitions after the existing Config struct

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

**Verification**:
```bash
go build ./cmd
# Should compile without errors
```

#### Step 1.2: Write Config Struct Test

**File**: `cmd/root_test.go` (create if doesn't exist)

**Action**: Add test to verify struct can be unmarshaled from YAML

```go
package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigUnmarshal_NewStructure(t *testing.T) {
	yamlContent := `
client-id: "test-id"
client-secret: "test-secret"
common:
  var1: "value1"
  var2: 456
hoge:
  fuga: "world"
  foo:
    bar: "baz"
`

	vp := viper.New()
	vp.SetConfigType("yaml")
	err := vp.ReadConfig(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var cfg Config
	err = vp.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Verify top-level fields
	if cfg.ClientID != "test-id" {
		t.Errorf("ClientID = %q, want %q", cfg.ClientID, "test-id")
	}
	if cfg.ClientSecret != "test-secret" {
		t.Errorf("ClientSecret = %q, want %q", cfg.ClientSecret, "test-secret")
	}

	// Verify common section
	if cfg.Common.Var1 != "value1" {
		t.Errorf("Common.Var1 = %q, want %q", cfg.Common.Var1, "value1")
	}
	if cfg.Common.Var2 != 456 {
		t.Errorf("Common.Var2 = %d, want %d", cfg.Common.Var2, 456)
	}

	// Verify hoge section (2 levels)
	if cfg.Hoge.Fuga != "world" {
		t.Errorf("Hoge.Fuga = %q, want %q", cfg.Hoge.Fuga, "world")
	}
	if cfg.Hoge.Foo.Bar != "baz" {
		t.Errorf("Hoge.Foo.Bar = %q, want %q", cfg.Hoge.Foo.Bar, "baz")
	}
}

func TestConfigUnmarshal_BackwardCompatibility(t *testing.T) {
	yamlContent := `
client-id: "old-id"
client-secret: "old-secret"
`

	vp := viper.New()
	vp.SetConfigType("yaml")
	err := vp.ReadConfig(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	var cfg Config
	err = vp.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Old fields should load
	if cfg.ClientID != "old-id" {
		t.Errorf("ClientID = %q, want %q", cfg.ClientID, "old-id")
	}
	if cfg.ClientSecret != "old-secret" {
		t.Errorf("ClientSecret = %q, want %q", cfg.ClientSecret, "old-secret")
	}

	// New fields should be zero values
	if cfg.Common.Var1 != "" {
		t.Errorf("Common.Var1 = %q, want empty string", cfg.Common.Var1)
	}
	if cfg.Common.Var2 != 0 {
		t.Errorf("Common.Var2 = %d, want 0", cfg.Common.Var2)
	}
	if cfg.Hoge.Fuga != "" {
		t.Errorf("Hoge.Fuga = %q, want empty string", cfg.Hoge.Fuga)
	}
	if cfg.Hoge.Foo.Bar != "" {
		t.Errorf("Hoge.Foo.Bar = %q, want empty string", cfg.Hoge.Foo.Bar)
	}
}
```

**Verification**:
```bash
go test -v ./cmd -run TestConfigUnmarshal
# Both tests should PASS
```

### Phase 2: Update BuildEffectiveConfig (15 minutes)

#### Step 2.1: Write BuildEffectiveConfig Test

**File**: `cmd/viperutils_test.go`

**Action**: Add tests for updated BuildEffectiveConfig function

```go
package cmd

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBuildEffectiveConfig_HasAllFields(t *testing.T) {
	cfg := BuildEffectiveConfig()

	// Check top-level keys
	if _, ok := cfg["client-id"]; !ok {
		t.Error("Missing key: client-id")
	}
	if _, ok := cfg["client-secret"]; !ok {
		t.Error("Missing key: client-secret")
	}
	if _, ok := cfg["common"]; !ok {
		t.Error("Missing key: common")
	}
	if _, ok := cfg["hoge"]; !ok {
		t.Error("Missing key: hoge")
	}

	// Check common nested keys
	common, ok := cfg["common"].(map[string]interface{})
	if !ok {
		t.Fatal("common is not a map")
	}
	if _, ok := common["var1"]; !ok {
		t.Error("Missing key: common.var1")
	}
	if _, ok := common["var2"]; !ok {
		t.Error("Missing key: common.var2")
	}

	// Check hoge nested keys (2 levels)
	hoge, ok := cfg["hoge"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge is not a map")
	}
	if _, ok := hoge["fuga"]; !ok {
		t.Error("Missing key: hoge.fuga")
	}
	foo, ok := hoge["foo"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge.foo is not a map")
	}
	if _, ok := foo["bar"]; !ok {
		t.Error("Missing key: hoge.foo.bar")
	}
}

func TestBuildEffectiveConfig_CorrectDefaultValues(t *testing.T) {
	cfg := BuildEffectiveConfig()

	// Top-level defaults
	if cfg["client-id"] != "" {
		t.Errorf("client-id = %q, want empty string", cfg["client-id"])
	}
	if cfg["client-secret"] != "" {
		t.Errorf("client-secret = %q, want empty string", cfg["client-secret"])
	}

	// Common defaults
	common := cfg["common"].(map[string]interface{})
	if common["var1"] != "" {
		t.Errorf("common.var1 = %q, want empty string", common["var1"])
	}
	if common["var2"] != 123 {
		t.Errorf("common.var2 = %v, want 123", common["var2"])
	}

	// Hoge defaults (2 levels)
	hoge := cfg["hoge"].(map[string]interface{})
	if hoge["fuga"] != "hello" {
		t.Errorf("hoge.fuga = %q, want \"hello\"", hoge["fuga"])
	}
	foo := hoge["foo"].(map[string]interface{})
	if foo["bar"] != "hello" {
		t.Errorf("hoge.foo.bar = %q, want \"hello\"", foo["bar"])
	}
}

func TestBuildEffectiveConfig_YAMLMarshal(t *testing.T) {
	cfg := BuildEffectiveConfig()

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal to YAML: %v", err)
	}

	// Unmarshal back to verify structure
	var result map[string]interface{}
	err = yaml.Unmarshal(yamlBytes, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Verify structure is preserved
	if _, ok := result["common"]; !ok {
		t.Error("common section missing after YAML round-trip")
	}
	if _, ok := result["hoge"]; !ok {
		t.Error("hoge section missing after YAML round-trip")
	}
}
```

**Verification** (tests will FAIL initially):
```bash
go test -v ./cmd -run TestBuildEffectiveConfig
# Tests should FAIL (Red phase)
```

#### Step 2.2: Implement BuildEffectiveConfig

**File**: `cmd/viperutils.go`

**Action**: Update BuildEffectiveConfig to return new structure

Replace the existing function:

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

**Verification** (tests should now PASS):
```bash
go test -v ./cmd -run TestBuildEffectiveConfig
# All tests should PASS (Green phase)
```

### Phase 3: Integration Testing (10 minutes)

#### Step 3.1: Test Full Workflow

**File**: `cmd/configure_test.go` (update if exists)

**Action**: Add or update integration test

```go
func TestConfigureCommand_GeneratesNewStructure(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()
	
	// Override config path for testing
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", oldHome)
	})

	// Execute configure command
	cmd := rootCmd
	configDir := filepath.Join(tmpDir, CliConfigBase, CliName)
	configFile := filepath.Join(configDir, GetConfigFile(DefaultProfile))
	
	// Create configure command with --force flag
	configureCmd.Flags().Set("force", "true")
	err := configureCmd.RunE(configureCmd, []string{})
	if err != nil {
		t.Fatalf("Configure command failed: %v", err)
	}

	// Read generated file
	yamlBytes, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	// Parse YAML
	var result map[string]interface{}
	err = yaml.Unmarshal(yamlBytes, &result)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	// Verify structure
	if _, ok := result["client-id"]; !ok {
		t.Error("Generated YAML missing client-id")
	}
	if _, ok := result["common"]; !ok {
		t.Error("Generated YAML missing common section")
	}
	if _, ok := result["hoge"]; !ok {
		t.Error("Generated YAML missing hoge section")
	}

	// Verify nested structure
	common, ok := result["common"].(map[string]interface{})
	if !ok {
		t.Fatal("common is not a map")
	}
	if common["var2"] != 123 {
		t.Errorf("common.var2 = %v, want 123", common["var2"])
	}

	hoge, ok := result["hoge"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge is not a map")
	}
	foo, ok := hoge["foo"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge.foo is not a map")
	}
	if foo["bar"] != "hello" {
		t.Errorf("hoge.foo.bar = %q, want \"hello\"", foo["bar"])
	}
}
```

**Verification**:
```bash
go test -v ./cmd -run TestConfigureCommand_GeneratesNewStructure
# Test should PASS
```

### Phase 4: Quality Checks (5 minutes)

#### Step 4.1: Run All Tests

```bash
make test
# All existing and new tests should PASS
```

#### Step 4.2: Format Code

```bash
make fmt
# Code should be formatted
```

#### Step 4.3: Lint Code

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
make lint
# Should pass with 0 warnings/errors
```

#### Step 4.4: Build Binary

```bash
make build
# Binary should build successfully in bin/mycli
```

### Phase 5: Manual Verification (5 minutes)

#### Step 5.1: Generate Config File

```bash
./bin/mycli configure --force
```

**Expected Output** (stderr):
```
Wrote config: /Users/<username>/.config/mycli/default.yaml
```

#### Step 5.2: Inspect Generated File

```bash
cat ~/.config/mycli/default.yaml
```

**Expected Content**:
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

#### Step 5.3: Test with Profile

```bash
./bin/mycli configure --profile dev --force
cat ~/.config/mycli/dev.yaml
```

**Expected**: Same structure as default.yaml

#### Step 5.4: Test Editor Integration (Optional)

```bash
export EDITOR=cat
./bin/mycli configure --force --edit
```

**Expected**: Config file content displayed via cat

### Phase 6: Cleanup and Commit (5 minutes)

#### Step 6.1: Run Full Build

```bash
make all
# Should complete without errors
```

#### Step 6.2: Check Git Status

```bash
git status
```

**Expected Modified Files**:
- `cmd/root.go` (Config struct updated)
- `cmd/viperutils.go` (BuildEffectiveConfig updated)
- `cmd/root_test.go` (new tests)
- `cmd/viperutils_test.go` (new tests)
- `cmd/configure_test.go` (updated tests, if applicable)

#### Step 6.3: Commit Changes

```bash
git add cmd/root.go cmd/viperutils.go cmd/root_test.go cmd/viperutils_test.go cmd/configure_test.go
git commit -m "feat: expand Config struct to support nested configuration

- Add CommonConfig, HogeConfig, FooConfig structs
- Update BuildEffectiveConfig() to return nested structure
- Add integration tests for Viper unmarshaling
- Maintain backward compatibility with old config files

Refs: specs/005-configure-refactor/spec.md"
```

## Troubleshooting

### Issue: Tests Fail with "mapstructure: cannot find field"

**Cause**: mapstructure tag mismatch

**Solution**: Verify all mapstructure tags match YAML keys exactly (kebab-case)

### Issue: Viper Unmarshal Returns Zero Values

**Cause**: YAML key names don't match mapstructure tags

**Solution**: Check YAML indentation and key names

### Issue: Type Mismatch in Tests

**Cause**: YAML parser returns different type than expected

**Solution**: Use type assertion with error checking:
```go
common, ok := cfg["common"].(map[string]interface{})
if !ok {
    t.Fatal("common is not a map")
}
```

### Issue: Generated YAML Has Wrong Structure

**Cause**: BuildEffectiveConfig returns incorrect map structure

**Solution**: Verify nested map literals in BuildEffectiveConfig

## Success Criteria Verification

After completing all phases, verify:

- [ ] **SC-001**: Generated YAML has 10 lines (including blank lines)
- [ ] **SC-002**: All 7 fields present and parseable
- [ ] **SC-003**: `CliConfig.Common.Var2 == 123` and `CliConfig.Hoge.Foo.Bar == "hello"` after unmarshal
- [ ] **SC-004**: `make test` passes 100%
- [ ] **SC-005**: `make lint` exits with code 0
- [ ] **SC-006**: Profile config generation works

## Time Estimate

| Phase | Time | Cumulative |
|-------|------|------------|
| Phase 1: Struct Definition | 10 min | 10 min |
| Phase 2: BuildEffectiveConfig | 15 min | 25 min |
| Phase 3: Integration Testing | 10 min | 35 min |
| Phase 4: Quality Checks | 5 min | 40 min |
| Phase 5: Manual Verification | 5 min | 45 min |
| Phase 6: Cleanup & Commit | 5 min | 50 min |

**Total**: ~50 minutes

## Next Steps

After completing this implementation:

1. Run `/speckit.tasks` to generate detailed task breakdown
2. Create PR with changes
3. Verify CI pipeline passes
4. Request code review
5. Merge to main branch

## References

- **Spec**: [spec.md](./spec.md)
- **Data Model**: [data-model.md](./data-model.md)
- **Contracts**: [contracts/](./contracts/)
- **Research**: [research.md](./research.md)
- **Constitution**: `.specify/memory/constitution.md`
