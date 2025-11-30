package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigureCreatesFile(t *testing.T) {
	// prepare temp HOME
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// backup and restore globals
	oldCliConfig := CliConfig
	oldCfgForce := cfgForce
	oldCfgOpen := cfgEdit
	oldProfile := profile
	t.Cleanup(func() {
		CliConfig = oldCliConfig
		cfgForce = oldCfgForce
		cfgEdit = oldCfgOpen
		profile = oldProfile
	})

	// set config values to scaffold
	CliConfig.ClientID = "test-id"
	CliConfig.ClientSecret = "test-secret"

	// ensure profile empty
	profile = ""
	cfgForce = false
	cfgEdit = false

	// run command
	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed: %v", err)
	}

	// check file exists
	cfgPath := filepath.Join(GetConfigPath(), GetConfigFile(DefaultProfile))
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("expected config file to exist at %s: %v", cfgPath, err)
	}

	// verify contents
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	s := string(b)
	if !contains(s, "client-id: test-id") {
		t.Fatalf("config content missing client-id; got:\n%s", s)
	}
	if !contains(s, "client-secret: test-secret") {
		t.Fatalf("config content missing client-secret; got:\n%s", s)
	}
}

func TestConfigureOpenInvokesEditor(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// backup and restore globals
	oldCliConfig := CliConfig
	oldCfgForce := cfgForce
	oldCfgOpen := cfgEdit
	oldProfile := profile
	oldEditor := os.Getenv("EDITOR")
	t.Cleanup(func() {
		CliConfig = oldCliConfig
		cfgForce = oldCfgForce
		cfgEdit = oldCfgOpen
		profile = oldProfile
		if err := os.Setenv("EDITOR", oldEditor); err != nil {
			t.Fatalf("failed to restore EDITOR: %v", err)
		}
	})

	CliConfig.ClientID = "x"
	CliConfig.ClientSecret = "y"
	cfgForce = false
	cfgEdit = true

	// set EDITOR to a no-op that exists on PATH; use 'true'
	t.Setenv("EDITOR", "true")

	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed: %v", err)
	}

	// file should exist
	cfgPath := filepath.Join(GetConfigPath(), GetConfigFile(DefaultProfile))
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("expected config file to exist at %s: %v", cfgPath, err)
	}
}

// contains is a tiny helper that does substring check but avoids importing strings
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}

func TestConfigureProfileCreatesProfileFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// backup/restore globals
	oldCliConfig := CliConfig
	oldCfgForce := cfgForce
	oldCfgOpen := cfgEdit
	oldProfile := profile
	t.Cleanup(func() {
		CliConfig = oldCliConfig
		cfgForce = oldCfgForce
		cfgEdit = oldCfgOpen
		profile = oldProfile
	})

	CliConfig.ClientID = "p-id"
	CliConfig.ClientSecret = "p-secret"
	profile = "prod"
	cfgForce = false
	cfgEdit = false

	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed for profile: %v", err)
	}

	// check file exists for profile
	cfgPath := filepath.Join(GetConfigPath(), GetConfigFile("prod"))
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("expected profile config file to exist at %s: %v", cfgPath, err)
	}
}

func TestConfigureForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// backup/restore globals
	oldCliConfig := CliConfig
	oldCfgForce := cfgForce
	oldCfgOpen := cfgEdit
	oldProfile := profile
	t.Cleanup(func() {
		CliConfig = oldCliConfig
		cfgForce = oldCfgForce
		cfgEdit = oldCfgOpen
		profile = oldProfile
	})

	CliConfig.ClientID = "f-id"
	CliConfig.ClientSecret = "f-secret"
	profile = ""

	cfgDir := GetConfigPath()
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	target := filepath.Join(cfgDir, GetConfigFile(DefaultProfile))
	// create existing file
	if err := os.WriteFile(target, []byte("client-id: old\n"), 0o644); err != nil {
		t.Fatalf("write existing: %v", err)
	}

	cfgForce = false
	cfgEdit = false
	// when not forcing, RunE should succeed but not overwrite the existing file
	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed when file exists and not forcing: %v", err)
	}

	// verify file still contains old client-id (was not overwritten)
	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read existing after skip: %v", err)
	}
	if !contains(string(b), "client-id: old") {
		t.Fatalf("expected existing file to remain with client-id old; got:\n%s", string(b))
	}

	// now force overwrite
	cfgForce = true
	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed with --force: %v", err)
	}

	// verify file now contains new client-id
	b, err = os.ReadFile(target)
	if err != nil {
		t.Fatalf("read after force: %v", err)
	}
	if !contains(string(b), "client-id: f-id") {
		t.Fatalf("expected file to be overwritten with client-id f-id; got:\n%s", string(b))
	}
}
