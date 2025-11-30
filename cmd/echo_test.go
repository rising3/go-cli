package cmd

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// captureOutput captures stdout and stderr from a command execution
func captureOutput(t *testing.T, cmd *cobra.Command, args []string) (stdout, stderr string, err error) {
	t.Helper()

	// Reset all flags to their default values before each test
	// This is necessary because Cobra persists flag values between SetArgs() calls
	for _, subcmd := range cmd.Commands() {
		subcmd.Flags().VisitAll(func(flag *pflag.Flag) {
			flag.Changed = false
			_ = flag.Value.Set(flag.DefValue)
		})
	}

	// Create buffers to capture output
	stdoutBuf := new(bytes.Buffer)
	stderrBuf := new(bytes.Buffer)

	// Set output buffers
	cmd.SetOut(stdoutBuf)
	cmd.SetErr(stderrBuf)
	cmd.SetArgs(args)

	// Execute command
	err = cmd.Execute()

	return stdoutBuf.String(), stderrBuf.String(), err
}

// T006: Single argument test
func TestEchoCommand_SingleArgument(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "Hello"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello\n"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T007: Multiple arguments test
func TestEchoCommand_MultipleArguments(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "A", "B", "C"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "A B C\n"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T008: No arguments test
func TestEchoCommand_NoArguments(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "\n"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T009: Special characters test
func TestEchoCommand_SpecialCharacters(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "!@#$%"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "!@#$%\n"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T016: -n flag test
func TestEchoCommand_NoNewlineFlag(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "-n", "Hello"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T017: -n flag with multiple arguments test
func TestEchoCommand_NoNewlineFlagMultipleArgs(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "-n", "A", "B"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "A B"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T018: -n flag with no arguments test
func TestEchoCommand_NoNewlineFlagNoArgs(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "-n"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := ""
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T035: -e flag with escape sequences test
func TestEchoCommand_EscapeFlag(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "-e", "Hello\\nWorld"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello\nWorld\n"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T036: Without -e flag, escapes should be literal
func TestEchoCommand_NoEscapeFlag(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "Hello\\nWorld"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Hello\\nWorld\n"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T053: -n -e combination test
func TestEchoCommand_NoNewlineAndEscapeFlags(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "-n", "-e", "Tab\\there"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Tab\there"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T054: -e -n combination test (reversed order)
func TestEchoCommand_EscapeAndNoNewlineFlags(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "-e", "-n", "Line\\nNo"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Line\nNo"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T056: Comprehensive flag combination test
func TestEchoCommand_FlagCombinations(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "no flags",
			args: []string{"echo", "Hello"},
			want: "Hello\n",
		},
		{
			name: "only -n",
			args: []string{"echo", "-n", "Hello"},
			want: "Hello",
		},
		{
			name: "only -e",
			args: []string{"echo", "-e", "Hello\\nWorld"},
			want: "Hello\nWorld\n",
		},
		{
			name: "-n -e combination",
			args: []string{"echo", "-n", "-e", "Hello\\tWorld"},
			want: "Hello\tWorld",
		},
		{
			name: "-e -n combination (reversed)",
			args: []string{"echo", "-e", "-n", "Hello\\tWorld"},
			want: "Hello\tWorld",
		},
		{
			name: "-e with \\c (overrides -n)",
			args: []string{"echo", "-n", "-e", "Hello\\cWorld"},
			want: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := rootCmd
			stdout, _, err := captureOutput(t, cmd, tt.args)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if stdout != tt.want {
				t.Errorf("output = %q, want %q", stdout, tt.want)
			}
		})
	}
}

// T060: Invalid flag error test
func TestEchoCommand_InvalidFlag(t *testing.T) {
	cmd := rootCmd
	_, _, err := captureOutput(t, cmd, []string{"echo", "-x", "test"})

	if err == nil {
		t.Fatal("expected error for invalid flag, got nil")
	}
}

// T061: Exit code test
func TestEchoCommand_ExitCodes(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "success case",
			args:      []string{"echo", "test"},
			wantError: false,
		},
		{
			name:      "invalid flag",
			args:      []string{"echo", "-x"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := rootCmd
			_, _, err := captureOutput(t, cmd, tt.args)

			if tt.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// T067: Empty string argument test
func TestEchoCommand_EmptyStringArgument(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "", "test"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := " test\n"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T068: Double dash argument separator test
func TestEchoCommand_DoubleDashSeparator(t *testing.T) {
	cmd := rootCmd
	stdout, _, err := captureOutput(t, cmd, []string{"echo", "-n", "--", "-e"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "-e"
	if stdout != want {
		t.Errorf("output = %q, want %q", stdout, want)
	}
}

// T064: Verbose flag test
func TestEchoCommand_VerboseFlag(t *testing.T) {
	cmd := rootCmd
	stdout, stderr, err := captureOutput(t, cmd, []string{"echo", "--verbose", "test"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stdout != "test\n" {
		t.Errorf("stdout = %q, want %q", stdout, "test\n")
	}

	// Verify debug output in stderr
	if stderr == "" {
		t.Error("expected debug output in stderr, got empty string")
	}
	if !strings.Contains(stderr, "[DEBUG]") {
		t.Errorf("stderr missing [DEBUG] prefix: %q", stderr)
	}
}

// T065: Performance test - Memory usage
func TestEchoCommand_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	// Test startup memory (< 50MB)
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	cmd := rootCmd
	_, _, err := captureOutput(t, cmd, []string{"echo", "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	startupMemMB := float64(m2.Alloc-m1.Alloc) / 1024 / 1024

	if startupMemMB > 50 {
		t.Errorf("startup memory usage %.2f MB exceeds 50 MB limit", startupMemMB)
	}

	// Test with 10,000 arguments (< 100MB)
	runtime.ReadMemStats(&m1)

	largeArgs := make([]string, 10002)
	largeArgs[0] = "echo"
	for i := 1; i < 10002; i++ {
		largeArgs[i] = "arg"
	}

	_, _, err = captureOutput(t, cmd, largeArgs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	runtime.ReadMemStats(&m2)
	largeMemMB := float64(m2.Alloc-m1.Alloc) / 1024 / 1024

	if largeMemMB > 100 {
		t.Errorf("10,000 args memory usage %.2f MB exceeds 100 MB limit", largeMemMB)
	}
}

// T066: Performance test - Help display speed
func TestEchoCommand_HelpPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	start := time.Now()

	cmd := rootCmd
	_, _, err := captureOutput(t, cmd, []string{"echo", "--help"})

	elapsed := time.Since(start)

	// Help should display in < 50ms
	if elapsed > 50*time.Millisecond {
		t.Errorf("help display took %v, exceeds 50ms limit", elapsed)
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// T076a: UTF-8 test cases
func TestEchoCommand_UTF8(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "Japanese characters",
			args: []string{"echo", "こんにちは世界"},
			want: "こんにちは世界\n",
		},
		{
			name: "Emoji",
			args: []string{"echo", "🚀✨"},
			want: "🚀✨\n",
		},
		{
			name: "Mixed UTF-8 with escape",
			args: []string{"echo", "-e", "Hello\\n世界🌍"},
			want: "Hello\n世界🌍\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := rootCmd
			stdout, _, err := captureOutput(t, cmd, tt.args)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if stdout != tt.want {
				t.Errorf("output = %q, want %q", stdout, tt.want)
			}
		})
	}
}
