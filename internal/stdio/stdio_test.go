package stdio

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestOpenWriterDashReturnsStdout(t *testing.T) {
	w, c, err := OpenWriter("-")
	if err != nil {
		t.Fatalf("OpenWriter returned error: %v", err)
	}
	if c != nil {
		t.Fatalf("expected nil closer for stdout, got: %v", c)
	}
	if w != os.Stdout {
		t.Fatalf("expected os.Stdout writer, got: %T", w)
	}
}

func TestBindCommandWiresStreams(t *testing.T) {
	var out bytes.Buffer
	s := Streams{In: nil, Out: &out, Err: io.Discard}

	// use a shell to printf so behavior is consistent cross-platform where sh exists
	cmd := exec.Command("sh", "-c", "printf hello")
	BindCommand(cmd, s)
	if err := cmd.Run(); err != nil {
		t.Fatalf("cmd.Run failed: %v", err)
	}
	if out.String() != "hello" {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestOpenReaderDashReturnsStdin(t *testing.T) {
	r, c, err := OpenReader("-")
	if err != nil {
		t.Fatalf("OpenReader returned error: %v", err)
	}
	if c != nil {
		t.Fatalf("expected nil closer for stdin, got: %v", c)
	}
	if r != os.Stdin {
		t.Fatalf("expected os.Stdin reader, got: %T", r)
	}
}

func TestOpenReaderFileReadsAndCloses(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/in.txt"
	content := []byte("input-data")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write tmp file: %v", err)
	}
	r, c, err := OpenReader(path)
	if err != nil {
		t.Fatalf("OpenReader returned error: %v", err)
	}
	if c == nil {
		t.Fatalf("expected closer for file reader")
	}
	buf := make([]byte, len(content))
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("read failed: %v", err)
	}
	if n != len(content) || string(buf) != string(content) {
		t.Fatalf("unexpected read: got %q (n=%d)", string(buf[:n]), n)
	}
	CloseAll(c)
}

func TestOpenWriterWithPermCreatesFile(t *testing.T) {
	tmp := t.TempDir()
	path := tmp + "/out.txt"
	w, c, err := OpenWriterWithPerm(path, 0o640)
	if err != nil {
		t.Fatalf("OpenWriterWithPerm returned error: %v", err)
	}
	if c == nil {
		t.Fatalf("expected closer for file writer")
	}
	n, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}
	if n != 5 {
		t.Fatalf("unexpected write bytes: %d", n)
	}
	CloseAll(c)
	// ensure file exists and has content
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read created file: %v", err)
	}
	if string(b) != "hello" {
		t.Fatalf("unexpected file content: %q", string(b))
	}
}
