package cmd

import (
	"bytes"
	"testing"

	"github.com/rising3/go-cli/internal/cmd/echo"
	"github.com/spf13/cobra"
)

func TestEchoWrapperCallsInternal(t *testing.T) {
	// capture arguments passed to EchoFunc
	var calledSuppressNewline bool
	var calledInterpretEscapes bool
	var calledVerbose bool
	var calledArgs []string

	// stub internal implementation
	old := echo.EchoFunc
	echo.EchoFunc = func(opts echo.EchoOptions) error {
		calledSuppressNewline = opts.SuppressNewline
		calledInterpretEscapes = opts.InterpretEscapes
		calledVerbose = opts.Verbose
		calledArgs = opts.Args
		return nil
	}
	defer func() { echo.EchoFunc = old }()

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	// Add flags to the test command (copy from echoCmd)
	cmd.Flags().BoolP("no-newline", "n", false, "do not output the trailing newline")
	cmd.Flags().BoolP("escape", "e", false, "interpret backslash escapes")
	cmd.Flags().Bool("verbose", false, "enable debug logging to stderr")

	// Set flags
	if err := cmd.Flags().Set("no-newline", "true"); err != nil {
		t.Fatalf("failed to set no-newline flag: %v", err)
	}
	if err := cmd.Flags().Set("escape", "true"); err != nil {
		t.Fatalf("failed to set escape flag: %v", err)
	}
	if err := cmd.Flags().Set("verbose", "false"); err != nil {
		t.Fatalf("failed to set verbose flag: %v", err)
	}

	// Execute
	args := []string{"Hello", "World"}
	if err := echoCmd.RunE(cmd, args); err != nil {
		t.Fatalf("echo RunE failed: %v", err)
	}

	// Verify internal EchoFunc was called with correct options
	if !calledSuppressNewline {
		t.Errorf("expected SuppressNewline true, got false")
	}
	if !calledInterpretEscapes {
		t.Errorf("expected InterpretEscapes true, got false")
	}
	if calledVerbose {
		t.Errorf("expected Verbose false, got true")
	}
	if len(calledArgs) != 2 || calledArgs[0] != "Hello" || calledArgs[1] != "World" {
		t.Errorf("expected args [Hello World], got %v", calledArgs)
	}
}

func TestEchoWrapper_DefaultFlags(t *testing.T) {
	var calledSuppressNewline bool
	var calledInterpretEscapes bool
	var calledVerbose bool

	// stub internal implementation
	old := echo.EchoFunc
	echo.EchoFunc = func(opts echo.EchoOptions) error {
		calledSuppressNewline = opts.SuppressNewline
		calledInterpretEscapes = opts.InterpretEscapes
		calledVerbose = opts.Verbose
		return nil
	}
	defer func() { echo.EchoFunc = old }()

	// Create command with flags at default values
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("no-newline", "n", false, "do not output the trailing newline")
	cmd.Flags().BoolP("escape", "e", false, "interpret backslash escapes")
	cmd.Flags().Bool("verbose", false, "enable debug logging to stderr")

	// Execute with default flags
	args := []string{"test"}
	if err := echoCmd.RunE(cmd, args); err != nil {
		t.Fatalf("echo RunE failed: %v", err)
	}

	// Verify all flags are false by default
	if calledSuppressNewline {
		t.Errorf("expected SuppressNewline false by default, got true")
	}
	if calledInterpretEscapes {
		t.Errorf("expected InterpretEscapes false by default, got true")
	}
	if calledVerbose {
		t.Errorf("expected Verbose false by default, got true")
	}
}

func TestEchoWrapper_VerboseFlag(t *testing.T) {
	var calledVerbose bool

	// stub internal implementation
	old := echo.EchoFunc
	echo.EchoFunc = func(opts echo.EchoOptions) error {
		calledVerbose = opts.Verbose
		return nil
	}
	defer func() { echo.EchoFunc = old }()

	// Create command with verbose flag enabled
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("no-newline", "n", false, "do not output the trailing newline")
	cmd.Flags().BoolP("escape", "e", false, "interpret backslash escapes")
	cmd.Flags().Bool("verbose", false, "enable debug logging to stderr")

	if err := cmd.Flags().Set("verbose", "true"); err != nil {
		t.Fatalf("failed to set verbose flag: %v", err)
	}

	// Execute
	args := []string{"verbose", "test"}
	if err := echoCmd.RunE(cmd, args); err != nil {
		t.Fatalf("echo RunE failed: %v", err)
	}

	// Verify verbose flag was passed correctly
	if !calledVerbose {
		t.Errorf("expected Verbose true, got false")
	}
}

func TestEchoWrapper_OutputStreams(t *testing.T) {
	var capturedOutput *bytes.Buffer
	var capturedErrOutput *bytes.Buffer

	// stub internal implementation
	old := echo.EchoFunc
	echo.EchoFunc = func(opts echo.EchoOptions) error {
		// Capture the output streams
		if buf, ok := opts.Output.(*bytes.Buffer); ok {
			capturedOutput = buf
		}
		if buf, ok := opts.ErrOutput.(*bytes.Buffer); ok {
			capturedErrOutput = buf
		}
		return nil
	}
	defer func() { echo.EchoFunc = old }()

	// Create command with custom output streams
	cmd := &cobra.Command{}
	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)

	cmd.Flags().BoolP("no-newline", "n", false, "do not output the trailing newline")
	cmd.Flags().BoolP("escape", "e", false, "interpret backslash escapes")
	cmd.Flags().Bool("verbose", false, "enable debug logging to stderr")

	// Execute
	if err := echoCmd.RunE(cmd, []string{"test"}); err != nil {
		t.Fatalf("echo RunE failed: %v", err)
	}

	// Verify output streams were passed correctly
	if capturedOutput != outBuf {
		t.Errorf("Output stream not passed correctly")
	}
	if capturedErrOutput != errBuf {
		t.Errorf("ErrOutput stream not passed correctly")
	}
}

func TestEchoWrapper_EmptyArgs(t *testing.T) {
	var calledArgs []string

	// stub internal implementation
	old := echo.EchoFunc
	echo.EchoFunc = func(opts echo.EchoOptions) error {
		calledArgs = opts.Args
		return nil
	}
	defer func() { echo.EchoFunc = old }()

	// Create command
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("no-newline", "n", false, "do not output the trailing newline")
	cmd.Flags().BoolP("escape", "e", false, "interpret backslash escapes")
	cmd.Flags().Bool("verbose", false, "enable debug logging to stderr")

	// Execute with empty args
	if err := echoCmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("echo RunE failed: %v", err)
	}

	// Verify empty args were passed
	if len(calledArgs) != 0 {
		t.Errorf("expected empty args, got %v", calledArgs)
	}
}
