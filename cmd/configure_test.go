package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rising3/go-cli/internal/cmd/configure"
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

	// verify contents - now generates defaults, not CliConfig values
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	s := string(b)
	// Verify structure contains all required keys
	if !contains(s, "client-id:") {
		t.Fatalf("config content missing client-id; got:\n%s", s)
	}
	if !contains(s, "client-secret:") {
		t.Fatalf("config content missing client-secret; got:\n%s", s)
	}
	if !contains(s, "common:") {
		t.Fatalf("config content missing common section; got:\n%s", s)
	}
	if !contains(s, "hoge:") {
		t.Fatalf("config content missing hoge section; got:\n%s", s)
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

	// verify file now contains new structure with defaults (not old content)
	b, err = os.ReadFile(target)
	if err != nil {
		t.Fatalf("read after force: %v", err)
	}
	s := string(b)
	// File should be overwritten with default structure
	if contains(s, "client-id: old") {
		t.Fatalf("expected file to be overwritten, but still contains old content:\n%s", s)
	}
	if !contains(s, "client-id:") {
		t.Fatalf("expected file to contain client-id key; got:\n%s", s)
	}
	if !contains(s, "common:") {
		t.Fatalf("expected file to contain common section; got:\n%s", s)
	}
	if !contains(s, "hoge:") {
		t.Fatalf("expected file to contain hoge section; got:\n%s", s)
	}
}

// T030: Test --force flag is passed correctly
func TestConfigureCommand_ForceFlag(t *testing.T) {
	// Mock ConfigureFunc
	oldFunc := configure.ConfigureFunc
	defer func() { configure.ConfigureFunc = oldFunc }()

	var capturedTarget string
	var capturedOpts configure.ConfigureOptions

	configure.ConfigureFunc = func(target string, opts configure.ConfigureOptions) error {
		capturedTarget = target
		capturedOpts = opts
		return nil
	}

	// Setup
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	oldCfgForce := cfgForce
	oldProfile := profile
	t.Cleanup(func() {
		cfgForce = oldCfgForce
		profile = oldProfile
	})

	cfgForce = true
	profile = ""

	// Execute
	cmd := &cobra.Command{}
	if err := configureCmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Force flag was passed
	if !capturedOpts.Force {
		t.Error("Force flag not passed correctly")
	}

	// Verify target path is correct
	expectedTarget := filepath.Join(GetConfigPath(), GetConfigFile(DefaultProfile))
	if capturedTarget != expectedTarget {
		t.Errorf("target = %q, want %q", capturedTarget, expectedTarget)
	}
}

// T031: Test --edit flag is passed correctly
func TestConfigureCommand_EditFlag(t *testing.T) {
	// Mock ConfigureFunc
	oldFunc := configure.ConfigureFunc
	defer func() { configure.ConfigureFunc = oldFunc }()

	var capturedOpts configure.ConfigureOptions

	configure.ConfigureFunc = func(target string, opts configure.ConfigureOptions) error {
		capturedOpts = opts
		return nil
	}

	// Setup
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	oldCfgEdit := cfgEdit
	t.Cleanup(func() {
		cfgEdit = oldCfgEdit
	})

	cfgEdit = true

	// Execute
	cmd := &cobra.Command{}
	if err := configureCmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify Edit flag was passed
	if !capturedOpts.Edit {
		t.Error("Edit flag not passed correctly")
	}
}

// T032: Test --no-wait flag is passed correctly
func TestConfigureCommand_NoWaitFlag(t *testing.T) {
	// Mock ConfigureFunc
	oldFunc := configure.ConfigureFunc
	defer func() { configure.ConfigureFunc = oldFunc }()

	var capturedOpts configure.ConfigureOptions

	configure.ConfigureFunc = func(target string, opts configure.ConfigureOptions) error {
		capturedOpts = opts
		return nil
	}

	// Setup
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	oldCfgNoWait := cfgNoWait
	t.Cleanup(func() {
		cfgNoWait = oldCfgNoWait
	})

	cfgNoWait = true

	// Execute
	cmd := &cobra.Command{}
	if err := configureCmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify EditorShouldWait returns false (no wait)
	if capturedOpts.EditorShouldWait == nil {
		t.Fatal("EditorShouldWait function not set")
	}

	shouldWait := capturedOpts.EditorShouldWait("", nil)
	if shouldWait {
		t.Error("Expected EditorShouldWait to return false when NoWait is true")
	}
}

// T033: Test --profile flag changes target path
func TestConfigureCommand_ProfileFlag(t *testing.T) {
	// Mock ConfigureFunc
	oldFunc := configure.ConfigureFunc
	defer func() { configure.ConfigureFunc = oldFunc }()

	var capturedTarget string

	configure.ConfigureFunc = func(target string, opts configure.ConfigureOptions) error {
		capturedTarget = target
		return nil
	}

	// Setup
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	oldProfile := profile
	t.Cleanup(func() {
		profile = oldProfile
	})

	profile = "production"

	// Execute
	cmd := &cobra.Command{}
	if err := configureCmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify target path includes profile name
	expectedTarget := filepath.Join(GetConfigPath(), GetConfigFile("production"))
	if capturedTarget != expectedTarget {
		t.Errorf("target = %q, want %q", capturedTarget, expectedTarget)
	}
}

// TestConfigureCommand_GeneratesNewStructure verifies that the configure
// command generates a YAML file with the new nested structure containing
// all required sections (common, hoge.foo). This validates User Story 4:
// backward compatibility while generating new structure.
func TestConfigureCommand_GeneratesNewStructure(t *testing.T) {
	// Setup temporary HOME
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// Backup and restore globals
	oldCfgForce := cfgForce
	oldCfgEdit := cfgEdit
	oldProfile := profile
	t.Cleanup(func() {
		cfgForce = oldCfgForce
		cfgEdit = oldCfgEdit
		profile = oldProfile
	})

	// Set flags for force creation without editor
	cfgForce = true
	cfgEdit = false
	profile = ""

	// Execute configure command
	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed: %v", err)
	}

	// Verify file exists
	cfgPath := filepath.Join(GetConfigPath(), GetConfigFile(DefaultProfile))
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("expected config file to exist at %s: %v", cfgPath, err)
	}

	// Read and verify YAML structure
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	content := string(b)

	// Verify all required top-level keys
	if !contains(content, "client-id:") {
		t.Errorf("missing client-id in generated YAML:\n%s", content)
	}
	if !contains(content, "client-secret:") {
		t.Errorf("missing client-secret in generated YAML:\n%s", content)
	}

	// Verify nested common section
	if !contains(content, "common:") {
		t.Errorf("missing common section in generated YAML:\n%s", content)
	}
	if !contains(content, "var1:") {
		t.Errorf("missing var1 in common section:\n%s", content)
	}
	if !contains(content, "var2:") {
		t.Errorf("missing var2 in common section:\n%s", content)
	}

	// Verify nested hoge section
	if !contains(content, "hoge:") {
		t.Errorf("missing hoge section in generated YAML:\n%s", content)
	}
	if !contains(content, "fuga:") {
		t.Errorf("missing fuga in hoge section:\n%s", content)
	}

	// Verify deeply nested foo section
	if !contains(content, "foo:") {
		t.Errorf("missing foo section in hoge:\n%s", content)
	}
	if !contains(content, "bar:") {
		t.Errorf("missing bar in foo section:\n%s", content)
	}
}
