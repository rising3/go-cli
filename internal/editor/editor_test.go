package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func createFakeExecutable(t *testing.T, dir, name string) string {
	t.Helper()
	exeName := name
	if runtime.GOOS == "windows" {
		exeName = name + ".exe"
	}
	path := filepath.Join(dir, exeName)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create fake exe: %v", err)
	}
	if runtime.GOOS != "windows" {
		if _, err := f.WriteString("#!/bin/sh\nexit 0\n"); err != nil {
			if cerr := f.Close(); cerr != nil {
				t.Fatalf("failed to write to fake exe: %v; close error: %v", err, cerr)
			}
			t.Fatalf("failed to write to fake exe: %v", err)
		}
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close fake exe: %v", err)
	}
	if err := os.Chmod(path, 0755); err != nil {
		t.Fatalf("failed to chmod fake exe: %v", err)
	}
	return path
}

func TestGetEditor_EditorEnvFound(t *testing.T) {
	dir := t.TempDir()
	createFakeExecutable(t, dir, "myedit")

	oldPath := os.Getenv("PATH")
	defer func() {
		if err := os.Setenv("PATH", oldPath); err != nil {
			t.Fatalf("failed to restore PATH: %v", err)
		}
	}()
	if err := os.Setenv("PATH", fmt.Sprintf("%s%c%s", dir, os.PathListSeparator, oldPath)); err != nil {
		t.Fatalf("failed to set PATH: %v", err)
	}

	oldEditor := os.Getenv("EDITOR")
	defer func() {
		if err := os.Setenv("EDITOR", oldEditor); err != nil {
			t.Fatalf("failed to restore EDITOR: %v", err)
		}
	}()
	if err := os.Setenv("EDITOR", "myedit --flag"); err != nil {
		t.Fatalf("failed to set EDITOR: %v", err)
	}

	name, args, err := GetEditor()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if name != "myedit" {
		t.Fatalf("expected name myedit, got %s", name)
	}
	if len(args) != 1 || args[0] != "--flag" {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestGetEditor_EditorEnvNotFound(t *testing.T) {
	oldEditor := os.Getenv("EDITOR")
	defer func() {
		if err := os.Setenv("EDITOR", oldEditor); err != nil {
			t.Fatalf("failed to restore EDITOR: %v", err)
		}
	}()
	if err := os.Setenv("EDITOR", "no-such-editor"); err != nil {
		t.Fatalf("failed to set EDITOR: %v", err)
	}

	oldPath := os.Getenv("PATH")
	defer func() {
		if err := os.Setenv("PATH", oldPath); err != nil {
			t.Fatalf("failed to restore PATH: %v", err)
		}
	}()
	if err := os.Setenv("PATH", ""); err != nil {
		t.Fatalf("failed to set PATH: %v", err)
	}

	_, _, err := GetEditor()
	if err == nil {
		t.Fatalf("expected error when EDITOR not found")
	}
}

func TestGetEditor_CandidateFound(t *testing.T) {
	oldEditor := os.Getenv("EDITOR")
	defer func() {
		if err := os.Setenv("EDITOR", oldEditor); err != nil {
			t.Fatalf("failed to restore EDITOR: %v", err)
		}
	}()
	if err := os.Unsetenv("EDITOR"); err != nil {
		t.Fatalf("failed to unset EDITOR: %v", err)
	}

	dir := t.TempDir()
	var candidate string
	switch runtime.GOOS {
	case "darwin":
		candidate = "open"
	case "windows":
		candidate = "notepad"
	default:
		candidate = "vim"
	}
	createFakeExecutable(t, dir, candidate)

	oldPath := os.Getenv("PATH")
	defer func() {
		if err := os.Setenv("PATH", oldPath); err != nil {
			t.Fatalf("failed to restore PATH: %v", err)
		}
	}()
	if err := os.Setenv("PATH", fmt.Sprintf("%s%c%s", dir, os.PathListSeparator, oldPath)); err != nil {
		t.Fatalf("failed to set PATH: %v", err)
	}

	name, args, err := GetEditor()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if name != candidate {
		t.Fatalf("expected candidate %s, got %s", candidate, name)
	}
	if args != nil {
		t.Fatalf("expected nil args for candidate, got %#v", args)
	}
}

func TestGetEditor_NoEditorFound(t *testing.T) {
	oldEditor := os.Getenv("EDITOR")
	defer func() {
		if err := os.Setenv("EDITOR", oldEditor); err != nil {
			t.Fatalf("failed to restore EDITOR: %v", err)
		}
	}()
	if err := os.Unsetenv("EDITOR"); err != nil {
		t.Fatalf("failed to unset EDITOR: %v", err)
	}

	oldPath := os.Getenv("PATH")
	defer func() {
		if err := os.Setenv("PATH", oldPath); err != nil {
			t.Fatalf("failed to restore PATH: %v", err)
		}
	}()
	empty := t.TempDir()
	if err := os.Setenv("PATH", empty); err != nil {
		t.Fatalf("failed to set PATH: %v", err)
	}

	_, _, err := GetEditor()
	if err == nil {
		t.Fatalf("expected error when no editor found")
	}
}
