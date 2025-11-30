package internalcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	proc "github.com/rising3/go-cli/internal/proc"
	stdio "github.com/rising3/go-cli/internal/stdio"

	"gopkg.in/yaml.v3"
)

type ConfigureOptions struct {
	Force        bool
	Edit         bool
	Data         map[string]interface{}
	Format       string // "yaml" or "json"
	Streams      stdio.Streams
	EditorLookup func() (string, []string, error)
	// EditorShouldWait, if provided, determines whether runEditor should wait for
	// the editor process to exit. If nil, the default is to wait.
	EditorShouldWait func(editor string, args []string) bool
}

// ConfigureFile ensures the config at target exists according to options.
// It delegates small tasks to helpers for readability and testability.
func ConfigureFile(target string, opts ConfigureOptions) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	if err := writeConfigIfNeeded(target, opts); err != nil {
		return err
	}

	if opts.Edit {
		return runEditor(target, opts)
	}
	return nil
}

// writeConfigIfNeeded writes the config to target when it does not exist or when Force is true.
// If the file exists and Force is false, it prints a message and returns nil.
func writeConfigIfNeeded(target string, opts ConfigureOptions) error {
	if _, err := os.Stat(target); err == nil && !opts.Force {
		if opts.Streams.Err != nil {
			_, _ = fmt.Fprintln(opts.Streams.Err, "Config already exists, skipping initialization:", target)
		}
		return nil
	}

	if opts.Force {
		_ = os.Remove(target)
	}

	out, err := marshalData(opts.Data, opts.Format)
	if err != nil {
		return err
	}

	// use stdio.OpenWriter so callers can use "-" to write to stdout or a file path
	w, closer, err := stdio.OpenWriterWithPerm(target, 0o644)
	if err != nil {
		return err
	}
	defer stdio.CloseAll(closer)

	if _, err := w.Write(out); err != nil {
		return err
	}
	if opts.Streams.Err != nil {
		_, _ = fmt.Fprintln(opts.Streams.Err, "Wrote config:", target)
	}
	return nil
}

// marshalData turns the provided map into YAML or JSON according to format.
func marshalData(data map[string]interface{}, format string) ([]byte, error) {
	switch format {
	case "yaml", "yml":
		return yaml.Marshal(data)
	default:
		return json.MarshalIndent(data, "", "  ")
	}
}

// runEditor launches the configured editor for the given target. It mirrors the previous
// behavior of swallowing certain editor discovery/launch errors after logging to Stderr.
func runEditor(target string, opts ConfigureOptions) error {
	if opts.EditorLookup == nil {
		return fmt.Errorf("no editor lookup provided")
	}
	ed, edArgs, err := opts.EditorLookup()
	if err != nil {
		if opts.Streams.Err != nil {
			_, _ = fmt.Fprintln(opts.Streams.Err, "No editor found:", err)
		}
		return nil
	}
	args := append(edArgs, target)
	cmd := proc.ExecCommand(ed, args...)
	// bind standard IO streams via helper
	stdio.BindCommand(cmd, opts.Streams)
	shouldWait := true
	if opts.EditorShouldWait != nil {
		shouldWait = opts.EditorShouldWait(ed, args)
	}

	return proc.Run(cmd, shouldWait, opts.Streams.Err)
}

// execCommand was previously a variable so tests could override it; the
// replacement is `proc.ExecCommand` in the `internal/proc` package.

// ConfigureFunc is a variable indirection so callers (cmd package tests)
// can replace the implementation with a stub. By default it points to ConfigureFile.
var ConfigureFunc = ConfigureFile
