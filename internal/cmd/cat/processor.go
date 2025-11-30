package cat

import (
	"bufio"
	"io"
	"os"
)

// Processor handles file reading and streaming
type Processor interface {
	// ProcessFile reads and processes a file
	// Parameters:
	//   - filename: path to the file (or "-" for stdin)
	//   - opts: formatting options
	//   - output: where to write formatted output
	// Returns:
	//   - error if file cannot be opened or read
	ProcessFile(filename string, opts Options, output io.Writer) error

	// ProcessStdin reads and processes standard input
	// Parameters:
	//   - opts: formatting options
	//   - output: where to write formatted output
	// Returns:
	//   - error if stdin cannot be read
	ProcessStdin(opts Options, output io.Writer) error
}

// DefaultProcessor is the standard implementation of Processor
type DefaultProcessor struct {
	formatter   Formatter
	stdinReader io.Reader // injectable for testing
}

// NewDefaultProcessor creates a new DefaultProcessor
func NewDefaultProcessor(formatter Formatter) *DefaultProcessor {
	return &DefaultProcessor{
		formatter:   formatter,
		stdinReader: os.Stdin,
	}
}

// ProcessFile implements Processor interface
func (p *DefaultProcessor) ProcessFile(filename string, opts Options, output io.Writer) error {
	// Handle "-" as stdin
	if filename == "-" {
		return p.ProcessStdin(opts, output)
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return p.processReader(file, opts, output)
}

// processReader is a helper that reads from any io.Reader
func (p *DefaultProcessor) processReader(reader io.Reader, opts Options, output io.Writer) error {
	scanner := bufio.NewScanner(reader)

	// Set 32KB buffer size per research.md
	buf := make([]byte, 32*1024)
	scanner.Buffer(buf, 32*1024)

	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		isEmpty := len(line) == 0

		// T048 [US4] Increment lineNum only for non-blank lines when NumberNonBlank is set
		if opts.NumberNonBlank {
			if !isEmpty {
				lineNum++
			}
		} else {
			lineNum++
		}

		formatted := p.formatter.FormatLine(line, lineNum, isEmpty, opts)
		if _, err := output.Write([]byte(formatted + "\n")); err != nil {
			return err
		}
	}

	return scanner.Err()
}

// ProcessStdin implements Processor interface
func (p *DefaultProcessor) ProcessStdin(opts Options, output io.Writer) error {
	return p.processReader(p.stdinReader, opts, output)
}
