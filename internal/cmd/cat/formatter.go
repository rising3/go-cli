package cat

// Formatter formats lines for output
type Formatter interface {
	// FormatLine formats a single line with the given options
	// Parameters:
	//   - line: the line content (without newline)
	//   - lineNum: the current line number (1-indexed)
	//   - isEmpty: whether the line is empty
	//   - opts: formatting options
	// Returns:
	//   - formatted line (without newline)
	FormatLine(line string, lineNum int, isEmpty bool, opts Options) string
}

// DefaultFormatter is the standard implementation of Formatter
type DefaultFormatter struct {
	controlCharMap map[byte]string
}

// NewDefaultFormatter creates a new DefaultFormatter
func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{
		controlCharMap: buildControlCharMap(),
	}
}

// FormatLine implements Formatter interface
func (f *DefaultFormatter) FormatLine(line string, lineNum int, isEmpty bool, opts Options) string {
	result := line

	// T060 [US6] ShowTabs: replace tabs with ^I
	if opts.ShowTabs {
		result = replaceAllBytes(result, '\t', "^I")
	}

	// T067 [US7] ShowNonPrinting: convert control chars using map
	if opts.ShowNonPrinting {
		result = f.convertControlChars(result)
	}

	// T047 [US4] NumberNonBlank: skip numbering for empty lines
	// T036 [US3] Line numbering with %6d format
	// T037 [US3] Overflow handling: lineNum % 1000000
	if opts.NumberAll || (opts.NumberNonBlank && !isEmpty) {
		displayNum := lineNum % 1000000
		result = formatLineNumber(displayNum) + result
	}

	// T054 [US5] ShowEnds: append $ at end of line
	if opts.ShowEnds {
		result = result + "$"
	}

	return result
}

// convertControlChars converts control characters using the map
func (f *DefaultFormatter) convertControlChars(s string) string {
	if len(s) == 0 {
		return s
	}

	var result []byte
	for i := 0; i < len(s); i++ {
		b := s[i]
		if converted, ok := f.controlCharMap[b]; ok {
			result = append(result, converted...)
		} else {
			result = append(result, b)
		}
	}
	return string(result)
}

// replaceAllBytes replaces all occurrences of a byte with a string
func replaceAllBytes(s string, old byte, new string) string {
	if len(s) == 0 {
		return s
	}

	var result []byte
	for i := 0; i < len(s); i++ {
		if s[i] == old {
			result = append(result, new...)
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// formatLineNumber formats a line number with 6-digit right-aligned format
func formatLineNumber(num int) string {
	// Format: "%6d  " (6 digits right-aligned + 2 spaces)
	const maxDigits = 6
	numStr := ""

	// Convert number to string
	if num == 0 {
		numStr = "0"
	} else {
		n := num
		for n > 0 {
			digit := n % 10
			numStr = string('0'+byte(digit)) + numStr
			n /= 10
		}
	}

	// Pad with spaces to 6 digits
	for len(numStr) < maxDigits {
		numStr = " " + numStr
	}

	return numStr + "  "
}

// buildControlCharMap creates the control character conversion map
// Maps ASCII 0-31 (excluding tab 9, newline 10) and ASCII 127 (DEL)
func buildControlCharMap() map[byte]string {
	m := make(map[byte]string)

	// ASCII 0-8
	m[0] = "^@"
	m[1] = "^A"
	m[2] = "^B"
	m[3] = "^C"
	m[4] = "^D"
	m[5] = "^E"
	m[6] = "^F"
	m[7] = "^G"
	m[8] = "^H"
	// 9 is tab - skip
	// 10 is newline - skip

	// ASCII 11-31
	m[11] = "^K"
	m[12] = "^L"
	m[13] = "^M"
	m[14] = "^N"
	m[15] = "^O"
	m[16] = "^P"
	m[17] = "^Q"
	m[18] = "^R"
	m[19] = "^S"
	m[20] = "^T"
	m[21] = "^U"
	m[22] = "^V"
	m[23] = "^W"
	m[24] = "^X"
	m[25] = "^Y"
	m[26] = "^Z"
	m[27] = "^["
	m[28] = "^\\"
	m[29] = "^]"
	m[30] = "^^"
	m[31] = "^_"

	// ASCII 127 (DEL)
	m[127] = "^?"

	return m
}
