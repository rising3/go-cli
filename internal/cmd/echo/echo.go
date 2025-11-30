package echo

import (
	"fmt"
	"io"
	"strings"
)

// EchoOptions represents the configuration for the echo command.
type EchoOptions struct {
	// SuppressNewline indicates whether to suppress the trailing newline (-n flag)
	SuppressNewline bool

	// InterpretEscapes indicates whether to interpret backslash escape sequences (-e flag)
	InterpretEscapes bool

	// Verbose enables debug logging to stderr (--verbose flag)
	Verbose bool

	// Args contains the text arguments to echo
	Args []string

	// Output is the writer for standard output
	Output io.Writer

	// ErrOutput is the writer for error/debug output
	ErrOutput io.Writer
}

// Echo performs the echo operation with the given options.
// It generates the output string, applies escape processing if needed,
// and writes to the configured output writer.
func Echo(opts EchoOptions) error {
	output, suppress := generateOutput(opts)
	return writeOutput(opts.Output, output, suppress)
}

// generateOutput generates the output string based on EchoOptions.
// Returns the processed output and whether newline should be suppressed.
func generateOutput(opts EchoOptions) (string, bool) {
	output := strings.Join(opts.Args, " ")

	// Process escape sequences if -e flag is set
	if opts.InterpretEscapes {
		processed, suppressFromEscape := ProcessEscapesFunc(output)
		// If \c is encountered, it takes precedence over -n flag
		return processed, suppressFromEscape || opts.SuppressNewline
	}

	return output, opts.SuppressNewline
}

// writeOutput writes the output to the given writer with optional newline.
func writeOutput(w io.Writer, output string, suppressNewline bool) error {
	if suppressNewline {
		_, err := fmt.Fprint(w, output)
		return err
	}
	_, err := fmt.Fprintln(w, output)
	return err
}

// Legacy functions for backward compatibility
// These will be deprecated in favor of Echo()

// GenerateOutput generates the output string based on EchoOptions.
// Deprecated: Use Echo() instead for better testability.
func GenerateOutput(opts EchoOptions) (string, bool) {
	return generateOutput(opts)
}

// WriteOutput writes the output to the given writer with optional newline.
// Deprecated: Use Echo() instead for better testability.
func WriteOutput(w io.Writer, output string, suppressNewline bool) error {
	return writeOutput(w, output, suppressNewline)
}

// ProcessEscapesFunc is a variable indirection so callers can replace
// the implementation with a stub for testing. By default it points to ProcessEscapes.
var ProcessEscapesFunc = ProcessEscapes

// EchoFunc is a variable indirection so callers (cmd package tests)
// can replace the implementation with a stub. By default it points to Echo.
var EchoFunc = Echo
