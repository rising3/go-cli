package proc

import (
	"bytes"
	"os/exec"
	"testing"
	"time"
)

func TestRun_StartFailsLogs(t *testing.T) {
	var buf bytes.Buffer
	// command that does not exist
	cmd := exec.Command("/no/such/editor")
	if err := Run(cmd, true, &buf); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Failed to start process")) {
		t.Fatalf("expected failure message in stderr, got: %s", buf.String())
	}
}

func TestRun_WaitLogsOnNonZeroExit(t *testing.T) {
	var buf bytes.Buffer
	// `false` exits with non-zero
	cmd := exec.Command("false")
	if err := Run(cmd, true, &buf); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("exited with error")) {
		t.Fatalf("expected exit error logged, got: %s", buf.String())
	}
}

func TestRun_NoWaitReturnsImmediately(t *testing.T) {
	// Use a short sleep and ensure Run returns immediately when shouldWait == false
	start := time.Now()
	cmd := exec.Command("sleep", "1")
	var buf bytes.Buffer
	if err := Run(cmd, false, &buf); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Fatalf("Run did not return immediately; elapsed=%v", time.Since(start))
	}
}

func TestExecCommandOverrideUsed(t *testing.T) {
	// ensure ExecCommand can be overridden by tests/callers
	old := ExecCommand
	defer func() { ExecCommand = old }()

	// override so created command is `true` which exits 0 immediately
	ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("true")
	}

	// use the overridden ExecCommand to build a command and run
	cmd := ExecCommand("dummy")
	var buf bytes.Buffer
	if err := Run(cmd, true, &buf); err != nil {
		t.Fatalf("Run returned unexpected error: %v", err)
	}
}
