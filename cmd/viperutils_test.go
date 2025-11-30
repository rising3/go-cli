package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func TestGetConfigFile(t *testing.T) {
	// Default profile
	got := GetConfigFile(DefaultProfile)
	want := DefaultProfile + "." + CliConfigType
	if got != want {
		t.Fatalf("GetConfigFile(DefaultProfile) = %q; want %q", got, want)
	}

	// custom profile name
	got2 := GetConfigFile("prod")
	want2 := "prod." + CliConfigType
	if got2 != want2 {
		t.Fatalf("GetConfigFile(prod) = %q; want %q", got2, want2)
	}
}

func TestInitViper_ReadInConfig(t *testing.T) {
	// create temporary home dir
	dir := t.TempDir()
	// set HOME for this test so GetConfigPath uses the temp dir
	t.Setenv("HOME", dir)

	cfgDir := filepath.Join(dir, CliConfigBase, CliName)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// write a default config file
	cfgFile := filepath.Join(cfgDir, DefaultProfile+"."+CliConfigType)
	content := "client-id: test-id\nclient-secret: test-secret\n"
	if err := os.WriteFile(cfgFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	vp := NewViper(DefaultProfile)
	if err := vp.ReadInConfig(); err != nil {
		t.Fatalf("ReadInConfig failed: %v", err)
	}

	// verify values
	if got := vp.GetString("client-id"); got != "test-id" {
		t.Fatalf("client-id = %q; want %q", got, "test-id")
	}
	if got := vp.GetString("client-secret"); got != "test-secret" {
		t.Fatalf("client-secret = %q; want %q", got, "test-secret")
	}

	// also verify that viper.Unmarshal into Config works
	var c Config
	if err := vp.Unmarshal(&c); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if c.ClientID != "test-id" || c.ClientSecret != "test-secret" {
		t.Fatalf("Unmarshaled config mismatch: %+v", c)
	}

	// ensure NewViper returns a fresh instance (no global state leakage)
	other := viper.New()
	other.SetConfigType(CliConfigType)
	other.AddConfigPath(cfgDir)
	other.SetConfigName("nonexistent")
	if err := other.ReadInConfig(); err == nil {
		t.Fatalf("expected ReadInConfig to fail for nonexistent config, but it succeeded")
	}
}

// TestBuildEffectiveConfig_HasAllFields verifies that the map returned by
// BuildEffectiveConfig contains all required keys (client-id, client-secret,
// common, hoge). This validates User Story 3: proper default value structure.
func TestBuildEffectiveConfig_HasAllFields(t *testing.T) {
	cfg := BuildEffectiveConfig()

	// Verify top-level keys exist
	if _, ok := cfg["client-id"]; !ok {
		t.Error("missing key: client-id")
	}
	if _, ok := cfg["client-secret"]; !ok {
		t.Error("missing key: client-secret")
	}
	if _, ok := cfg["common"]; !ok {
		t.Error("missing key: common")
	}
	if _, ok := cfg["hoge"]; !ok {
		t.Error("missing key: hoge")
	}

	// Verify common nested keys
	common, ok := cfg["common"].(map[string]interface{})
	if !ok {
		t.Fatal("common is not a map")
	}
	if _, ok := common["var1"]; !ok {
		t.Error("missing key: common.var1")
	}
	if _, ok := common["var2"]; !ok {
		t.Error("missing key: common.var2")
	}

	// Verify hoge nested keys
	hoge, ok := cfg["hoge"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge is not a map")
	}
	if _, ok := hoge["fuga"]; !ok {
		t.Error("missing key: hoge.fuga")
	}
	if _, ok := hoge["foo"]; !ok {
		t.Error("missing key: hoge.foo")
	}

	// Verify hoge.foo nested key
	foo, ok := hoge["foo"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge.foo is not a map")
	}
	if _, ok := foo["bar"]; !ok {
		t.Error("missing key: hoge.foo.bar")
	}
}

// TestBuildEffectiveConfig_CorrectDefaultValues verifies that all fields
// in the returned map have their correct default values.
func TestBuildEffectiveConfig_CorrectDefaultValues(t *testing.T) {
	cfg := BuildEffectiveConfig()

	// Verify top-level default values
	if got := cfg["client-id"].(string); got != "" {
		t.Errorf("client-id = %q; want empty string", got)
	}
	if got := cfg["client-secret"].(string); got != "" {
		t.Errorf("client-secret = %q; want empty string", got)
	}

	// Verify common section defaults
	common := cfg["common"].(map[string]interface{})
	if got := common["var1"].(string); got != "" {
		t.Errorf("common.var1 = %q; want empty string", got)
	}
	if got := common["var2"].(int); got != 123 {
		t.Errorf("common.var2 = %d; want 123", got)
	}

	// Verify hoge section defaults
	hoge := cfg["hoge"].(map[string]interface{})
	if got := hoge["fuga"].(string); got != "hello" {
		t.Errorf("hoge.fuga = %q; want %q", got, "hello")
	}

	// Verify hoge.foo section defaults
	foo := hoge["foo"].(map[string]interface{})
	if got := foo["bar"].(string); got != "hello" {
		t.Errorf("hoge.foo.bar = %q; want %q", got, "hello")
	}
}

// TestBuildEffectiveConfig_YAMLMarshal verifies that the map returned by
// BuildEffectiveConfig can be marshaled to YAML and unmarshaled back
// while preserving the nested structure and types.
func TestBuildEffectiveConfig_YAMLMarshal(t *testing.T) {
	cfg := BuildEffectiveConfig()

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("YAML marshal failed: %v", err)
	}

	// Unmarshal back to map
	var result map[string]interface{}
	if err := yaml.Unmarshal(yamlBytes, &result); err != nil {
		t.Fatalf("YAML unmarshal failed: %v", err)
	}

	// Verify top-level structure preserved
	if result["client-id"] != "" {
		t.Errorf("client-id after round-trip = %v; want empty string", result["client-id"])
	}
	if result["client-secret"] != "" {
		t.Errorf("client-secret after round-trip = %v; want empty string", result["client-secret"])
	}

	// Verify common section preserved
	common, ok := result["common"].(map[string]interface{})
	if !ok {
		t.Fatal("common is not a map after unmarshal")
	}
	if common["var1"] != "" {
		t.Errorf("common.var1 after round-trip = %v; want empty string", common["var1"])
	}
	if common["var2"] != 123 {
		t.Errorf("common.var2 after round-trip = %v; want 123", common["var2"])
	}

	// Verify hoge section preserved
	hoge, ok := result["hoge"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge is not a map after unmarshal")
	}
	if hoge["fuga"] != "hello" {
		t.Errorf("hoge.fuga after round-trip = %v; want hello", hoge["fuga"])
	}

	// Verify hoge.foo section preserved
	foo, ok := hoge["foo"].(map[string]interface{})
	if !ok {
		t.Fatal("hoge.foo is not a map after unmarshal")
	}
	if foo["bar"] != "hello" {
		t.Errorf("hoge.foo.bar after round-trip = %v; want hello", foo["bar"])
	}
}
