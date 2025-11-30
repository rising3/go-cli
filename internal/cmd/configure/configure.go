package configure

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rising3/go-cli/internal/proc"
	"gopkg.in/yaml.v3"
)

// ConfigureOptions represents the configuration for the configure command.
// It encapsulates all parameters needed to create and optionally edit a configuration file.
type ConfigureOptions struct {
	Force            bool                             // Force overwrites existing configuration files without prompting
	Edit             bool                             // Edit launches an editor after creating the configuration file
	NoWait           bool                             // NoWait runs the editor in background without blocking
	Data             map[string]interface{}           // Data contains the configuration data to be serialized
	Format           string                           // Format specifies the output format ("yaml", "yml", or "json")
	Output           io.Writer                        // Output is the standard output stream (currently unused, reserved for future use)
	ErrOutput        io.Writer                        // ErrOutput is the error output stream for messages
	EditorLookup     func() (string, []string, error) // EditorLookup is a function that returns the editor command and arguments
	EditorShouldWait func(string, []string) bool      // EditorShouldWait determines whether to wait for the editor to exit
}

// Configure creates or overwrites a configuration file at the specified target path.
// It marshals the opts.Data according to opts.Format and writes it to the file.
// If opts.Edit is true, it launches the configured editor after file creation.
//
// Parameters:
//   - target: Absolute path to the configuration file
//   - opts: Configuration options including data, format, and I/O streams
//
// Returns:
//   - error: Returns error if file creation fails or editor launch fails
//     (unless editor detection fails, in which case error is logged and nil is returned)
func Configure(target string, opts ConfigureOptions) error {
	// T013: Create parent directory
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	// T015: Check if file exists
	if _, err := os.Stat(target); err == nil && !opts.Force {
		if opts.ErrOutput != nil {
			_, _ = fmt.Fprintln(opts.ErrOutput, "Config already exists, skipping initialization:", target)
		}
		return nil
	}

	// Remove existing file if Force
	if opts.Force {
		_ = os.Remove(target)
	}

	// T014: Marshal data
	var out []byte
	var err error
	switch opts.Format {
	case "yaml", "yml":
		out, err = yaml.Marshal(opts.Data)
	default:
		out, err = json.MarshalIndent(opts.Data, "", "  ")
	}
	if err != nil {
		return err
	}

	// T013: Write file
	if err := os.WriteFile(target, out, 0o644); err != nil {
		return err
	}

	// T016: Write success message
	if opts.ErrOutput != nil {
		_, _ = fmt.Fprintln(opts.ErrOutput, "Wrote config:", target)
	}

	// T023-T025: Editor launch (Phase 3)
	if opts.Edit {
		// T023: Call EditorLookup
		ed, edArgs, err := opts.EditorLookup()
		if err != nil {
			// T024: Absorb EditorLookup errors
			if opts.ErrOutput != nil {
				_, _ = fmt.Fprintln(opts.ErrOutput, "No editor found:", err)
			}
			return nil
		}

		// T026: Use proc package for editor launch
		args := append(edArgs, target)
		cmd := proc.ExecCommand(ed, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// T025: Determine wait based on EditorShouldWait
		shouldWait := true
		if opts.EditorShouldWait != nil {
			shouldWait = opts.EditorShouldWait(ed, args)
		}

		// T026: Run editor via proc.Run
		return proc.Run(cmd, shouldWait, opts.ErrOutput)
	}

	return nil
}

// ConfigureFunc is a variable indirection for testing.
// By default it points to Configure.
var ConfigureFunc = Configure
