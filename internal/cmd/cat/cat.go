package cat

import (
	"fmt"
	"io"
	"os"
)

// CatFunc is the internal cat implementation (can be mocked in tests)
var CatFunc = func(filenames []string, opts Options) error {
	return catImpl(filenames, opts, os.Stdin, os.Stdout, os.Stderr)
}

// catImpl is the actual implementation that can be tested
func catImpl(filenames []string, opts Options, stdin io.Reader, stdout, stderr io.Writer) error {
	formatter := NewDefaultFormatter()
	processor := NewDefaultProcessor(formatter)
	processor.stdinReader = stdin // Inject stdin for testing

	// Handle stdin mode (no args)
	if len(filenames) == 0 {
		return processor.ProcessStdin(opts, stdout)
	}

	// Process each file
	hadError := false
	for _, filename := range filenames {
		if err := processor.ProcessFile(filename, opts, stdout); err != nil {
			_, _ = fmt.Fprintf(stderr, "cat: %s: %v\n", filename, err)
			hadError = true
			// Continue processing remaining files
		}
	}

	// Return error if any file failed (for exit code 1)
	if hadError {
		return fmt.Errorf("one or more files failed")
	}

	return nil
}
