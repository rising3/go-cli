package stdio

import (
	"io"
	"os"
	"os/exec"
)

// Streams groups standard IO streams. Use NewDefault() to get OS defaults.
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// NewDefault returns Streams bound to the process' stdio (os.Stdin/Stdout/Stderr).
func NewDefault() Streams {
	return Streams{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}
}

// OpenWriter returns an io.WriteCloser for the given path.
// Special cases:
//   - ""  -> use os.Stdout (no closer)
//   - "-" -> use os.Stdout (no closer)
//
// Otherwise it opens/creates the file for writing (with 0644) and returns it (must be closed).
func OpenWriter(path string) (io.Writer, io.Closer, error) {
	if path == "" || path == "-" {
		return os.Stdout, nil, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

// OpenWriterWithPerm behaves like OpenWriter but allows specifying the
// file permissions used when creating the file. This is useful when callers
// need to create files with non-default permission bits.
func OpenWriterWithPerm(path string, perm os.FileMode) (io.Writer, io.Closer, error) {
	if path == "" || path == "-" {
		return os.Stdout, nil, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

// OpenReader returns an io.ReadCloser for the given path.
// Special cases:
//   - ""  -> use os.Stdin (no closer)
//   - "-" -> use os.Stdin (no closer)
//
// Otherwise it opens the file for reading and returns it (must be closed).
func OpenReader(path string) (io.Reader, io.Closer, error) {
	if path == "" || path == "-" {
		return os.Stdin, nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

// BindCommand attaches the provided Streams to the command's stdio fields.
// This is a small convenience so callers can wire command IO to testing buffers
// or to OS stdio easily.
func BindCommand(cmd *exec.Cmd, s Streams) {
	if s.In != nil {
		cmd.Stdin = s.In
	}
	if s.Out != nil {
		cmd.Stdout = s.Out
	}
	if s.Err != nil {
		cmd.Stderr = s.Err
	}
}

// CloseAll closes any non-nil io.Closers passed and ignores nils/errors.
// Useful to close files returned from OpenWriter.
func CloseAll(closers ...io.Closer) {
	for _, c := range closers {
		if c == nil {
			continue
		}
		_ = c.Close()
	}
}
