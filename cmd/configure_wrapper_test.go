package cmd

import (
	"path/filepath"
	"testing"

	"github.com/rising3/go-cli/internal/cmd/configure"
	"github.com/spf13/cobra"
)

func TestConfigureWrapperCallsInternal(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// capture arguments passed to ConfigureFunc
	var calledTarget string
	var calledForce bool
	var calledEdit bool
	var calledFormat string
	var calledData map[string]interface{}
	var calledShouldWait bool

	// stub internal implementation
	old := configure.ConfigureFunc
	configure.ConfigureFunc = func(target string, opts configure.ConfigureOptions) error {
		calledTarget = target
		calledForce = opts.Force
		calledEdit = opts.Edit
		calledFormat = opts.Format
		calledData = opts.Data
		// capture EditorShouldWait result
		if opts.EditorShouldWait != nil {
			calledShouldWait = opts.EditorShouldWait("editor", []string{})
		}
		return nil
	}
	defer func() { configure.ConfigureFunc = old }()

	// ensure flags are set as expected
	cfgForce = false
	cfgEdit = true
	profile = "dev"
	// set CliConfig so cmd builds data map
	CliConfig.ClientID = "wrap-id"
	CliConfig.ClientSecret = "wrap-secret"

	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed: %v", err)
	}

	if calledTarget == "" {
		t.Fatalf("internal ConfigureFunc was not called")
	}
	// target should be inside HOME/.config/mycli
	if filepath.Base(calledTarget) != GetConfigFile("dev") {
		t.Fatalf("unexpected target file: %s", calledTarget)
	}
	if !calledEdit {
		t.Fatalf("expected Edit true passed to internal func")
	}
	if calledForce {
		t.Fatalf("expected Force false passed to internal func")
	}
	if calledFormat != CliConfigType {
		t.Fatalf("expected format %s, got %s", CliConfigType, calledFormat)
	}
	if calledData == nil {
		t.Fatalf("expected data to be passed to internal func")
	}
	if v, _ := calledData["client-id"].(string); v != "wrap-id" {
		t.Fatalf("expected client-id wrap-id passed, got %v", calledData["client-id"])
	}

	// by default --no-wait is false, so EditorShouldWait should be true
	if !calledShouldWait {
		t.Fatalf("expected EditorShouldWait true by default, got false")
	}
}

func TestConfigureWrapper_NoWaitFlagPassed(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	var calledShouldWait *bool

	// stub internal implementation
	old := configure.ConfigureFunc
	configure.ConfigureFunc = func(target string, opts configure.ConfigureOptions) error {
		if opts.EditorShouldWait != nil {
			v := opts.EditorShouldWait("editor", []string{})
			calledShouldWait = &v
		}
		return nil
	}
	defer func() { configure.ConfigureFunc = old }()

	// set flags: enable edit and set no-wait
	cfgForce = false
	cfgEdit = true
	cfgNoWait = true
	profile = "dev"

	// ensure CliConfig has values so BuildEffectiveConfig produces a map
	CliConfig.ClientID = "wrap-id"
	CliConfig.ClientSecret = "wrap-secret"

	if err := configureCmd.RunE(&cobra.Command{}, []string{}); err != nil {
		t.Fatalf("configure RunE failed: %v", err)
	}

	if calledShouldWait == nil {
		t.Fatalf("EditorShouldWait was not provided to internal func")
	}
	if *calledShouldWait {
		t.Fatalf("expected EditorShouldWait false when --no-wait is set, got true")
	}
}
