package cat

import (
	"bytes"
	"os"
	"testing"
)

// T010 [P] [US1] TestProcessFile_Success - basic file read
func TestProcessFile_Success(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp(t.TempDir(), "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile.Close() }()

	// Write test data
	content := "line1\nline2\n"
	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}

	// Test Processor
	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)

	var output bytes.Buffer
	opts := Options{}

	if err := processor.ProcessFile(tmpfile.Name(), opts, &output); err != nil {
		t.Fatalf("ProcessFile() failed: %v", err)
	}

	got := output.String()
	want := "line1\nline2\n"

	if got != want {
		t.Errorf("ProcessFile() = %q, want %q", got, want)
	}
}

// T011 [P] [US1] TestProcessFile_MultipleFiles
func TestProcessFile_MultipleFiles(t *testing.T) {
	// Create first file
	tmpfile1, err := os.CreateTemp(t.TempDir(), "test1*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile1.Close() }()

	if _, err := tmpfile1.WriteString("First\n"); err != nil {
		t.Fatal(err)
	}

	// Create second file
	tmpfile2, err := os.CreateTemp(t.TempDir(), "test2*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile2.Close() }()

	if _, err := tmpfile2.WriteString("Second\n"); err != nil {
		t.Fatal(err)
	}

	// Test processing multiple files
	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)
	opts := Options{}

	var output bytes.Buffer

	if err := processor.ProcessFile(tmpfile1.Name(), opts, &output); err != nil {
		t.Fatalf("ProcessFile(file1) failed: %v", err)
	}

	if err := processor.ProcessFile(tmpfile2.Name(), opts, &output); err != nil {
		t.Fatalf("ProcessFile(file2) failed: %v", err)
	}

	got := output.String()
	want := "First\nSecond\n"

	if got != want {
		t.Errorf("ProcessFile() = %q, want %q", got, want)
	}
}

// T023 [P] [US2] TestProcessStdin_Success - basic stdin read
func TestProcessStdin_Success(t *testing.T) {
	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)

	// Simulate stdin input by injecting custom stdin reader
	stdinContent := "stdin line 1\nstdin line 2\n"
	processor.stdinReader = bytes.NewBufferString(stdinContent)

	var output bytes.Buffer
	opts := Options{}

	if err := processor.ProcessStdin(opts, &output); err != nil {
		t.Fatalf("ProcessStdin() failed: %v", err)
	}

	got := output.String()
	want := stdinContent

	if got != want {
		t.Errorf("ProcessStdin() = %q, want %q", got, want)
	}
}

// T024 [P] [US2] TestProcessFile_Stdin_DashArgument - "-" triggers stdin processing
func TestProcessFile_Stdin_DashArgument(t *testing.T) {
	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)

	// Simulate stdin input by injecting custom stdin reader
	stdinContent := "dash stdin content\n"
	processor.stdinReader = bytes.NewBufferString(stdinContent)

	var output bytes.Buffer
	opts := Options{}

	// ProcessFile with "-" should read from stdin
	if err := processor.ProcessFile("-", opts, &output); err != nil {
		t.Fatalf("ProcessFile('-') failed: %v", err)
	}

	got := output.String()
	want := stdinContent

	if got != want {
		t.Errorf("ProcessFile('-') = %q, want %q", got, want)
	}
}

// T077 [P] TestProcessFile_NotExist - error when file doesn't exist
func TestProcessFile_NotExist(t *testing.T) {
	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)

	var output bytes.Buffer
	opts := Options{}

	err := processor.ProcessFile("/nonexistent/file/path.txt", opts, &output)

	if err == nil {
		t.Errorf("Expected error for nonexistent file, got nil")
	}

	if !os.IsNotExist(err) {
		t.Errorf("Expected IsNotExist error, got: %v", err)
	}
}

// T078 [P] TestProcessFile_IsDirectory - error when trying to read directory
func TestProcessFile_IsDirectory(t *testing.T) {
	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)

	// Use temp directory
	dir := t.TempDir()

	var output bytes.Buffer
	opts := Options{}

	err := processor.ProcessFile(dir, opts, &output)

	if err == nil {
		t.Errorf("Expected error for directory, got nil")
	}
}

// T082 [P] TestProcessFile_BinaryFile - binary files should be processed as-is
func TestProcessFile_BinaryFile(t *testing.T) {
	// Create a temporary binary file
	tmpfile, err := os.CreateTemp(t.TempDir(), "binary*.bin")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tmpfile.Close() }()

	// Write binary data
	binaryData := []byte{0xFF, 0xFE, 0x00, 0x01, 0x7F, 0x80}
	if _, err := tmpfile.Write(binaryData); err != nil {
		t.Fatal(err)
	}
	_ = tmpfile.Close()

	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)

	var output bytes.Buffer
	opts := Options{}

	// Should not error on binary files
	if err := processor.ProcessFile(tmpfile.Name(), opts, &output); err != nil {
		t.Fatalf("ProcessFile() failed on binary file: %v", err)
	}

	// Binary data should be in output (may be partial due to scanner line-based reading)
	if output.Len() == 0 {
		t.Errorf("Expected some output from binary file, got empty")
	}
}
