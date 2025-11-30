package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/rising3/go-cli/internal/cmd/cat"
)

// T012 [P] [US1] TestCatCommand_BasicFile - Cobra integration test
func TestCatCommand_BasicFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp(t.TempDir(), "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile.Close() }()

	content := "test line 1\ntest line 2\n"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	_ = tmpfile.Close()

	// Mock CatFunc
	var calledFilenames []string

	oldCatFunc := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledFilenames = filenames
		return nil
	}
	defer func() { cat.CatFunc = oldCatFunc }()

	// Execute command
	rootCmd.SetArgs([]string{"cat", tmpfile.Name()})
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	// Verify CatFunc was called with correct args
	if len(calledFilenames) != 1 || calledFilenames[0] != tmpfile.Name() {
		t.Errorf("Expected filenames [%s], got %v", tmpfile.Name(), calledFilenames)
	}
}

// T025 [P] [US2] TestCatCommand_StdinOnly - no args should read from stdin
func TestCatCommand_StdinOnly(t *testing.T) {
	// Mock CatFunc
	var calledFilenames []string

	oldCatFunc := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledFilenames = filenames
		return nil
	}
	defer func() { cat.CatFunc = oldCatFunc }()

	// Execute with no file args (stdin mode)
	rootCmd.SetArgs([]string{"cat"})
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	// Verify CatFunc was called with empty filenames (stdin mode)
	if len(calledFilenames) != 0 {
		t.Errorf("Expected empty filenames for stdin, got %v", calledFilenames)
	}
}

// T026 [P] [US2] TestCatCommand_MixedStdinAndFiles - mix of "-" and regular files
func TestCatCommand_MixedStdinAndFiles(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp(t.TempDir(), "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile.Close() }()

	if _, err := tmpfile.WriteString("file content\n"); err != nil {
		t.Fatal(err)
	}
	_ = tmpfile.Close()

	// Mock CatFunc
	var calledFilenames []string

	oldCatFunc := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledFilenames = filenames
		return nil
	}
	defer func() { cat.CatFunc = oldCatFunc }()

	// Execute with mix of file and "-" (stdin)
	rootCmd.SetArgs([]string{"cat", tmpfile.Name(), "-", tmpfile.Name()})
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	// Verify filenames include "-" for stdin
	want := []string{tmpfile.Name(), "-", tmpfile.Name()}
	if len(calledFilenames) != len(want) {
		t.Errorf("Expected %d filenames, got %d", len(want), len(calledFilenames))
	}

	for i := range want {
		if i >= len(calledFilenames) || calledFilenames[i] != want[i] {
			t.Errorf("Expected filenames %v, got %v", want, calledFilenames)
			break
		}
	}
}

// T035 [P] [US3] TestCatCommand_NumberFlag - verify -n flag is passed correctly
func TestCatCommand_NumberFlag(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp(t.TempDir(), "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile.Close() }()

	if _, err := tmpfile.WriteString("line 1\nline 2\n"); err != nil {
		t.Fatal(err)
	}
	_ = tmpfile.Close()

	// Mock CatFunc
	var calledOpts cat.Options

	oldCatFunc := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledOpts = opts
		return nil
	}
	defer func() { cat.CatFunc = oldCatFunc }()

	// Execute with -n flag
	rootCmd.SetArgs([]string{"cat", "-n", tmpfile.Name()})
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	// Verify NumberAll is set
	if !calledOpts.NumberAll {
		t.Errorf("Expected NumberAll true, got false")
	}
}

// T080 [P] TestCatCommand_PartialError - some files fail but processing continues
func TestCatCommand_PartialError(t *testing.T) {
	// Create one valid file
	tmpfile, err := os.CreateTemp(t.TempDir(), "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile.Close() }()

	if _, err := tmpfile.WriteString("valid content\n"); err != nil {
		t.Fatal(err)
	}
	_ = tmpfile.Close()

	// Mock CatFunc to simulate partial error
	var calledFilenames []string
	oldCatFunc := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledFilenames = filenames
		// Simulate the actual behavior: process valid files, error on invalid
		return fmt.Errorf("one or more files failed")
	}
	defer func() { cat.CatFunc = oldCatFunc }()

	// Execute with mix of valid and invalid files
	rootCmd.SetArgs([]string{"cat", tmpfile.Name(), "/nonexistent.txt", tmpfile.Name()})
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	err = rootCmd.Execute()

	// Should report error
	if err == nil {
		t.Errorf("Expected error for nonexistent file, got nil")
	}

	// Verify all filenames were passed to CatFunc
	if len(calledFilenames) != 3 {
		t.Errorf("Expected 3 filenames, got %d", len(calledFilenames))
	}
}

// T081 [P] TestCatCommand_EmptyFile - empty files should work without error
func TestCatCommand_EmptyFile(t *testing.T) {
	// Create an empty file
	tmpfile, err := os.CreateTemp(t.TempDir(), "empty*.txt")
	if err != nil {
		t.Fatal(err)
	}
	_ = tmpfile.Close()

	rootCmd.SetArgs([]string{"cat", tmpfile.Name()})
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() failed on empty file: %v", err)
	}

	// Empty file should produce no output
	if output.Len() != 0 {
		t.Errorf("Expected no output from empty file, got: %q", output.String())
	}
}
