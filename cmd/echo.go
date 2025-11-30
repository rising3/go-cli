// Package cmd implements CLI commands for the mycli application.
// T069: Package documentation for echo command
package cmd

import (
	"github.com/rising3/go-cli/internal/cmd/echo"
	"github.com/spf13/cobra"
)

var echoCmd = &cobra.Command{
	Use:   "echo [string...]",
	Short: "Output text to standard output",
	Long: `Echo writes the specified string(s) to standard output, separated by spaces,
followed by a newline.

This is a UNIX-compatible echo command implementation with support for
escape sequences and newline suppression options.`,
	SilenceUsage: false, // T058: Show usage on errors
	Example: `  # Basic output
  mycli echo "Hello, World!"
  
  # Multiple arguments
  mycli echo Hello World
  
  # Suppress newline
  mycli echo -n "Prompt: "
  
  # Interpret escape sequences
  mycli echo -e "Line1\nLine2\tTab"
  
  # Combine flags
  mycli echo -n -e "No newline\twith tab"
  
  # Special escape: \c suppresses output
  mycli echo -e "Stop here\cIgnored text"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flag values
		suppressNewline, _ := cmd.Flags().GetBool("no-newline")
		interpretEscapes, _ := cmd.Flags().GetBool("escape")
		verbose, _ := cmd.Flags().GetBool("verbose")

		// T063: Verbose logging
		if verbose {
			cmd.PrintErrf("[DEBUG] Args: %v\n", args)
			cmd.PrintErrf("[DEBUG] SuppressNewline: %v\n", suppressNewline)
			cmd.PrintErrf("[DEBUG] InterpretEscapes: %v\n", interpretEscapes)
		}

		// Create options with configured streams
		opts := echo.EchoOptions{
			SuppressNewline:  suppressNewline,
			InterpretEscapes: interpretEscapes,
			Verbose:          verbose,
			Args:             args,
			Output:           cmd.OutOrStdout(),
			ErrOutput:        cmd.ErrOrStderr(),
		}

		// Use the refactored Echo function via EchoFunc for testability
		return echo.EchoFunc(opts)
	},
}

func init() {
	// T019: Add -n/--no-newline flag
	echoCmd.Flags().BoolP("no-newline", "n", false, "do not output the trailing newline")

	// T048: Add -e/--escape flag
	echoCmd.Flags().BoolP("escape", "e", false, "interpret backslash escapes")

	// T062: Add --verbose flag
	echoCmd.Flags().Bool("verbose", false, "enable debug logging to stderr")

	rootCmd.AddCommand(echoCmd)
}
