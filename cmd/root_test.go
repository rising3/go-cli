package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestConfigUnmarshal_NewStructure verifies that Viper correctly unmarshals
// a complete nested YAML configuration into the Config struct.
// This test validates User Story 1: support for nested configuration structure.
func TestConfigUnmarshal_NewStructure(t *testing.T) {
	// Create temporary home directory
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfgDir := filepath.Join(dir, CliConfigBase, CliName)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write a complete nested config file
	cfgFile := filepath.Join(cfgDir, DefaultProfile+"."+CliConfigType)
	content := `client-id: test-client-id
client-secret: test-client-secret
common:
  var1: value1
  var2: 42
hoge:
  fuga: fuga-value
  foo:
    bar: bar-value
`
	if err := os.WriteFile(cfgFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Initialize Viper and read config
	vp := NewViper(DefaultProfile)
	if err := vp.ReadInConfig(); err != nil {
		t.Fatalf("ReadInConfig failed: %v", err)
	}

	// Unmarshal into Config struct
	var c Config
	if err := vp.Unmarshal(&c); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify top-level fields
	if c.ClientID != "test-client-id" {
		t.Errorf("ClientID = %q; want %q", c.ClientID, "test-client-id")
	}
	if c.ClientSecret != "test-client-secret" {
		t.Errorf("ClientSecret = %q; want %q", c.ClientSecret, "test-client-secret")
	}

	// Verify common nested fields
	if c.Common.Var1 != "value1" {
		t.Errorf("Common.Var1 = %q; want %q", c.Common.Var1, "value1")
	}
	if c.Common.Var2 != 42 {
		t.Errorf("Common.Var2 = %d; want %d", c.Common.Var2, 42)
	}

	// Verify hoge nested fields
	if c.Hoge.Fuga != "fuga-value" {
		t.Errorf("Hoge.Fuga = %q; want %q", c.Hoge.Fuga, "fuga-value")
	}

	// Verify hoge.foo deeply nested field
	if c.Hoge.Foo.Bar != "bar-value" {
		t.Errorf("Hoge.Foo.Bar = %q; want %q", c.Hoge.Foo.Bar, "bar-value")
	}
}

// TestConfigUnmarshal_BackwardCompatibility verifies that old configuration files
// (with only client-id and client-secret) still load correctly without errors.
// Missing nested fields should be initialized to their zero values.
// This test validates backward compatibility requirement.
func TestConfigUnmarshal_BackwardCompatibility(t *testing.T) {
	// Create temporary home directory
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfgDir := filepath.Join(dir, CliConfigBase, CliName)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write an old-style config file (only top-level fields)
	cfgFile := filepath.Join(cfgDir, DefaultProfile+"."+CliConfigType)
	content := `client-id: old-client-id
client-secret: old-client-secret
`
	if err := os.WriteFile(cfgFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Initialize Viper and read config
	vp := NewViper(DefaultProfile)
	if err := vp.ReadInConfig(); err != nil {
		t.Fatalf("ReadInConfig failed: %v", err)
	}

	// Unmarshal into Config struct - should succeed without errors
	var c Config
	if err := vp.Unmarshal(&c); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify top-level fields are loaded correctly
	if c.ClientID != "old-client-id" {
		t.Errorf("ClientID = %q; want %q", c.ClientID, "old-client-id")
	}
	if c.ClientSecret != "old-client-secret" {
		t.Errorf("ClientSecret = %q; want %q", c.ClientSecret, "old-client-secret")
	}

	// Verify nested fields have zero values (backward compatibility)
	if c.Common.Var1 != "" {
		t.Errorf("Common.Var1 = %q; want empty string", c.Common.Var1)
	}
	if c.Common.Var2 != 0 {
		t.Errorf("Common.Var2 = %d; want 0", c.Common.Var2)
	}
	if c.Hoge.Fuga != "" {
		t.Errorf("Hoge.Fuga = %q; want empty string", c.Hoge.Fuga)
	}
	if c.Hoge.Foo.Bar != "" {
		t.Errorf("Hoge.Foo.Bar = %q; want empty string", c.Hoge.Foo.Bar)
	}
}

// TestConfigUnmarshal_PartialStructure verifies that configuration files
// with only partial nested sections (e.g., only common) load correctly.
// This test validates User Story 2: proper struct validation and field handling.
func TestConfigUnmarshal_PartialStructure(t *testing.T) {
	// Create temporary home directory
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	cfgDir := filepath.Join(dir, CliConfigBase, CliName)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write a partial config file (only common section)
	cfgFile := filepath.Join(cfgDir, DefaultProfile+"."+CliConfigType)
	content := `client-id: partial-client-id
client-secret: partial-client-secret
common:
  var1: partial-value
  var2: 99
`
	if err := os.WriteFile(cfgFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	// Initialize Viper and read config
	vp := NewViper(DefaultProfile)
	if err := vp.ReadInConfig(); err != nil {
		t.Fatalf("ReadInConfig failed: %v", err)
	}

	// Unmarshal into Config struct - should succeed
	var c Config
	if err := vp.Unmarshal(&c); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify top-level fields
	if c.ClientID != "partial-client-id" {
		t.Errorf("ClientID = %q; want %q", c.ClientID, "partial-client-id")
	}
	if c.ClientSecret != "partial-client-secret" {
		t.Errorf("ClientSecret = %q; want %q", c.ClientSecret, "partial-client-secret")
	}

	// Verify common section is loaded
	if c.Common.Var1 != "partial-value" {
		t.Errorf("Common.Var1 = %q; want %q", c.Common.Var1, "partial-value")
	}
	if c.Common.Var2 != 99 {
		t.Errorf("Common.Var2 = %d; want %d", c.Common.Var2, 99)
	}

	// Verify hoge section has zero values (not present in config)
	if c.Hoge.Fuga != "" {
		t.Errorf("Hoge.Fuga = %q; want empty string", c.Hoge.Fuga)
	}
	if c.Hoge.Foo.Bar != "" {
		t.Errorf("Hoge.Foo.Bar = %q; want empty string", c.Hoge.Foo.Bar)
	}
}
