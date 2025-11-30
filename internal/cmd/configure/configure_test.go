package configure_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rising3/go-cli/internal/cmd/configure"
	"github.com/rising3/go-cli/internal/proc"
)

// T006: Basic file creation test
func TestConfigure_BasicFileCreation(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.yaml")

	data := map[string]interface{}{
		"key": "value",
	}

	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            false,
		Edit:             false,
		NoWait:           false,
		Data:             data,
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}

	// Execute
	err := configure.Configure(target, opts)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Errorf("file was not created: %s", target)
	}

	// Verify file content
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "key: value\n"
	if string(content) != expected {
		t.Errorf("content = %q, want %q", string(content), expected)
	}

	// Verify success message
	if !strings.Contains(errBuf.String(), "Wrote config:") {
		t.Errorf("expected success message, got: %s", errBuf.String())
	}
}

// T007: File exists without Force flag test
func TestConfigure_FileExists_NoForce(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "existing.yaml")

	// Create existing file
	if err := os.WriteFile(target, []byte("old: data\n"), 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            false,
		Data:             map[string]interface{}{"new": "data"},
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}

	err := configure.Configure(target, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file content unchanged
	content, _ := os.ReadFile(target)
	if string(content) != "old: data\n" {
		t.Errorf("file was modified when it shouldn't be")
	}

	// Verify message
	if !strings.Contains(errBuf.String(), "already exists") {
		t.Errorf("expected 'already exists' message, got: %s", errBuf.String())
	}
}

// T008: File exists with Force flag test
func TestConfigure_FileExists_Force(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "existing.yaml")

	// Create existing file
	if err := os.WriteFile(target, []byte("old: data\n"), 0o644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            true,
		Data:             map[string]interface{}{"new": "data"},
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}

	err := configure.Configure(target, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file content updated
	content, _ := os.ReadFile(target)
	expected := "new: data\n"
	if string(content) != expected {
		t.Errorf("content = %q, want %q", string(content), expected)
	}
}

// T009: Directory creation test
func TestConfigure_DirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "subdir", "nested", "config.yaml")

	data := map[string]interface{}{
		"test": "value",
	}

	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            false,
		Edit:             false,
		NoWait:           false,
		Data:             data,
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}

	err := configure.Configure(target, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify directory was created
	dir := filepath.Dir(target)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("directory was not created: %s", dir)
	}

	// Verify file exists
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Errorf("file was not created: %s", target)
	}
}

// T010: YAML format test
func TestConfigure_YAMLFormat(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.yaml")

	data := map[string]interface{}{
		"name":    "test",
		"enabled": true,
		"count":   42,
	}

	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            false,
		Data:             data,
		Format:           "yaml",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}

	err := configure.Configure(target, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file content is valid YAML
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// Check for YAML-specific formatting
	contentStr := string(content)
	if !strings.Contains(contentStr, "name: test") {
		t.Errorf("expected YAML format with 'name: test', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "enabled: true") {
		t.Errorf("expected YAML format with 'enabled: true', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "count: 42") {
		t.Errorf("expected YAML format with 'count: 42', got: %s", contentStr)
	}
}

// T011: JSON format test
func TestConfigure_JSONFormat(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.json")

	data := map[string]interface{}{
		"name":    "test",
		"enabled": true,
		"count":   42,
	}

	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:            false,
		Data:             data,
		Format:           "json",
		Output:           &bytes.Buffer{},
		ErrOutput:        &errBuf,
		EditorLookup:     func() (string, []string, error) { return "", nil, nil },
		EditorShouldWait: func(string, []string) bool { return true },
	}

	err := configure.Configure(target, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file content is valid JSON
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// Check for JSON-specific formatting
	contentStr := string(content)
	if !strings.Contains(contentStr, `"name": "test"`) {
		t.Errorf("expected JSON format with '\"name\": \"test\"', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, `"enabled": true`) {
		t.Errorf("expected JSON format with '\"enabled\": true', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, `"count": 42`) {
		t.Errorf("expected JSON format with '\"count\": 42', got: %s", contentStr)
	}
}

// T020: Editor found test
func TestConfigure_Edit_EditorFound(t *testing.T) {
	// Mock proc.ExecCommand to avoid launching actual editor
	oldExecCommand := proc.ExecCommand
	defer func() { proc.ExecCommand = oldExecCommand }()

	editorCommandCaptured := false
	proc.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		editorCommandCaptured = true
		// Return a command that exits immediately
		return exec.Command("true")
	}

	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.yaml")

	data := map[string]interface{}{
		"key": "value",
	}

	var errBuf bytes.Buffer
	editorCalled := false

	opts := configure.ConfigureOptions{
		Force:     false,
		Edit:      true,
		NoWait:    false,
		Data:      data,
		Format:    "yaml",
		Output:    &bytes.Buffer{},
		ErrOutput: &errBuf,
		EditorLookup: func() (string, []string, error) {
			editorCalled = true
			return "/usr/bin/vi", []string{}, nil
		},
		EditorShouldWait: func(ed string, args []string) bool {
			return true
		},
	}

	err := configure.Configure(target, opts)

	// Should not error
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify EditorLookup was called
	if !editorCalled {
		t.Error("EditorLookup was not called")
	}

	// Verify editor command was constructed
	if !editorCommandCaptured {
		t.Error("Editor command was not executed")
	}

	// Verify file was created before editor launch
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Errorf("file was not created: %s", target)
	}
}

// T021: Editor not found test
func TestConfigure_Edit_EditorNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.yaml")

	data := map[string]interface{}{
		"key": "value",
	}

	var errBuf bytes.Buffer
	opts := configure.ConfigureOptions{
		Force:     false,
		Edit:      true,
		NoWait:    false,
		Data:      data,
		Format:    "yaml",
		Output:    &bytes.Buffer{},
		ErrOutput: &errBuf,
		EditorLookup: func() (string, []string, error) {
			return "", nil, os.ErrNotExist
		},
		EditorShouldWait: func(string, []string) bool { return true },
	}

	err := configure.Configure(target, opts)

	// Error should be absorbed (return nil)
	if err != nil {
		t.Fatalf("expected nil error when editor not found, got: %v", err)
	}

	// Verify error message was written
	if !strings.Contains(errBuf.String(), "No editor found:") {
		t.Errorf("expected 'No editor found' message, got: %s", errBuf.String())
	}

	// Verify file was still created
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Errorf("file was not created: %s", target)
	}
}

// T022: NoWait test
func TestConfigure_Edit_NoWait(t *testing.T) {
	// Mock proc.ExecCommand to avoid launching actual editor
	oldExecCommand := proc.ExecCommand
	defer func() { proc.ExecCommand = oldExecCommand }()

	proc.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}

	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "config.yaml")

	data := map[string]interface{}{
		"key": "value",
	}

	var errBuf bytes.Buffer
	shouldWaitCalled := false
	var shouldWaitResult bool

	opts := configure.ConfigureOptions{
		Force:     false,
		Edit:      true,
		NoWait:    true,
		Data:      data,
		Format:    "yaml",
		Output:    &bytes.Buffer{},
		ErrOutput: &errBuf,
		EditorLookup: func() (string, []string, error) {
			return "/usr/bin/vi", []string{}, nil
		},
		EditorShouldWait: func(ed string, args []string) bool {
			shouldWaitCalled = true
			shouldWaitResult = false // NoWait = true means shouldWait = false
			return shouldWaitResult
		},
	}

	_ = configure.Configure(target, opts)

	// Verify EditorShouldWait was called
	if !shouldWaitCalled {
		t.Error("EditorShouldWait was not called")
	}

	// Verify it returned false (no wait)
	if shouldWaitResult != false {
		t.Errorf("expected EditorShouldWait to return false, got: %v", shouldWaitResult)
	}
}
