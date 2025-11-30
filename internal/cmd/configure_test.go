package internalcmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	proc "github.com/rising3/go-cli/internal/proc"
	stdio "github.com/rising3/go-cli/internal/stdio"
)

func TestConfigureFile_WritesWhenMissing(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config.yaml")

	opts := ConfigureOptions{
		Force: false,
		Edit:  false,
		Data: map[string]interface{}{
			"client-id":     "x-id",
			"client-secret": "x-secret",
		},
		Format:       "yaml",
		Streams:      stdio.Streams{In: nil, Out: nil, Err: os.Stderr},
		EditorLookup: nil,
	}

	if err := ConfigureFile(target, opts); err != nil {
		t.Fatalf("ConfigureFile failed: %v", err)
	}

	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	s := string(b)
	if !contains(s, "client-id: x-id") {
		t.Fatalf("expected client-id in file; got:\n%s", s)
	}
}

func TestConfigureFile_SkipsWhenExists(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(target, []byte("client-id: old\n"), 0o644); err != nil {
		t.Fatalf("write existing: %v", err)
	}

	opts := ConfigureOptions{
		Force: false,
		Edit:  false,
		Data: map[string]interface{}{
			"client-id":     "new-id",
			"client-secret": "new-secret",
		},
		Format:       "yaml",
		Streams:      stdio.Streams{In: nil, Out: nil, Err: os.Stderr},
		EditorLookup: nil,
	}

	if err := ConfigureFile(target, opts); err != nil {
		t.Fatalf("ConfigureFile failed: %v", err)
	}

	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read existing: %v", err)
	}
	if !contains(string(b), "client-id: old") {
		t.Fatalf("expected existing content to remain; got:\n%s", string(b))
	}
}

func TestConfigureFile_ForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(target, []byte("client-id: old\n"), 0o644); err != nil {
		t.Fatalf("write existing: %v", err)
	}

	opts := ConfigureOptions{
		Force: true,
		Edit:  false,
		Data: map[string]interface{}{
			"client-id":     "f-id",
			"client-secret": "f-secret",
		},
		Format:       "yaml",
		Streams:      stdio.Streams{In: nil, Out: nil, Err: os.Stderr},
		EditorLookup: nil,
	}

	if err := ConfigureFile(target, opts); err != nil {
		t.Fatalf("ConfigureFile failed: %v", err)
	}

	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read after force: %v", err)
	}
	if !contains(string(b), "client-id: f-id") {
		t.Fatalf("expected overwritten content; got:\n%s", string(b))
	}
}

func TestConfigureFile_EditInvokesEditor(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config.yaml")

	// ensure file exists so editor will be called
	opts := ConfigureOptions{
		Force: false,
		Edit:  true,
		Data: map[string]interface{}{
			"client-id":     "e-id",
			"client-secret": "e-secret",
		},
		Format:       "yaml",
		Streams:      stdio.Streams{In: nil, Out: nil, Err: os.Stderr},
		EditorLookup: func() (string, []string, error) { return "dummy", []string{"arg"}, nil },
	}

	// override proc.ExecCommand to run a test helper process that simply exits 0
	oldExec := proc.ExecCommand
	proc.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		args := append([]string{"-test.run=TestHelperProcess", "--", name}, arg...)
		cmd := exec.Command(os.Args[0], args...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		return cmd
	}
	defer func() { proc.ExecCommand = oldExec }()

	if err := ConfigureFile(target, opts); err != nil {
		t.Fatalf("ConfigureFile with edit failed: %v", err)
	}
}

// TestHelperProcess is used to mock exec.Command in tests.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

// contains copied from cmd package helper
func contains(s, sub string) bool {
	if len(s) < len(sub) {
		return false
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestMarshalData_JSON(t *testing.T) {
	data := map[string]interface{}{"client-id": "j-id", "client-secret": "s"}
	out, err := marshalData(data, "json")
	if err != nil {
		t.Fatalf("marshalData json failed: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "\"client-id\": \"j-id\"") {
		t.Fatalf("unexpected json output: %s", s)
	}
}

func TestMarshalData_YAML(t *testing.T) {
	data := map[string]interface{}{"client-id": "y-id"}
	out, err := marshalData(data, "yaml")
	if err != nil {
		t.Fatalf("marshalData yaml failed: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "client-id: y-id") {
		t.Fatalf("unexpected yaml output: %s", s)
	}
}

func TestRunEditor_NoLookup(t *testing.T) {
	opts := ConfigureOptions{Streams: stdio.Streams{Err: &bytes.Buffer{}}}
	if err := runEditor("/tmp/foo", opts); err == nil {
		t.Fatalf("expected error when EditorLookup is nil")
	}
}

func TestRunEditor_LookupErrorLogged(t *testing.T) {
	var buf bytes.Buffer
	opts := ConfigureOptions{
		Streams:      stdio.Streams{Err: &buf},
		EditorLookup: func() (string, []string, error) { return "", nil, errors.New("no editor") },
	}
	if err := runEditor("/tmp/foo", opts); err != nil {
		t.Fatalf("runEditor returned unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No editor found") {
		t.Fatalf("expected 'No editor found' logged; got: %s", buf.String())
	}
}

func TestRunEditor_StartFailsLogs(t *testing.T) {
	var buf bytes.Buffer
	// EditorLookup returns a non-existent path so Start will fail
	opts := ConfigureOptions{
		Streams:      stdio.Streams{Err: &buf},
		EditorLookup: func() (string, []string, error) { return "/no/such/editor", nil, nil },
	}

	old := proc.ExecCommand
	proc.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		// Return a command pointing to a non-existent path to simulate Start error
		return exec.Command("/no/such/editor")
	}
	defer func() { proc.ExecCommand = old }()

	if err := runEditor("/tmp/foo", opts); err != nil {
		t.Fatalf("runEditor returned unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Failed to start process") {
		t.Fatalf("expected 'Failed to start process' logged; got: %s", buf.String())
	}
}

func TestRunEditor_NoWaitReturnsImmediately(t *testing.T) {
	// Use a command that sleeps for 1s; when EditorShouldWait returns false,
	// runEditor should return immediately (well before the sleep finishes).
	opts := ConfigureOptions{
		Streams:          stdio.Streams{Err: &bytes.Buffer{}},
		EditorLookup:     func() (string, []string, error) { return "sleep", []string{"1"}, nil },
		EditorShouldWait: func(editor string, args []string) bool { return false },
	}

	start := time.Now()
	if err := runEditor("/tmp/foo", opts); err != nil {
		t.Fatalf("runEditor returned unexpected error: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed > 200*time.Millisecond {
		t.Fatalf("runEditor did not return immediately; elapsed=%v", elapsed)
	}
}
