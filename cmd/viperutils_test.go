package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
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
