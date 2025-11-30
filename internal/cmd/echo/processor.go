package echo

import "strings"

// ProcessEscapes processes escape sequences in the input string
// Returns the processed string and whether output should be suppressed (for \c)
func ProcessEscapes(input string) (output string, suppressNewline bool) {
	var builder strings.Builder
	builder.Grow(len(input)) // Pre-allocate for efficiency

	i := 0
	for i < len(input) {
		if input[i] == '\\' && i+1 < len(input) {
			// Process escape sequence
			switch input[i+1] {
			case 'n': // T038: newline
				builder.WriteRune('\n')
				i += 2
			case 't': // T039: tab
				builder.WriteRune('\t')
				i += 2
			case '\\': // T040: backslash
				builder.WriteRune('\\')
				i += 2
			case '"': // T041: double quote
				builder.WriteRune('"')
				i += 2
			case 'a': // T042: alert (bell)
				builder.WriteRune('\a')
				i += 2
			case 'b': // T043: backspace
				builder.WriteRune('\b')
				i += 2
			case 'c': // T044: suppress further output
				return builder.String(), true
			case 'r': // T045: carriage return
				builder.WriteRune('\r')
				i += 2
			case 'v': // T046: vertical tab
				builder.WriteRune('\v')
				i += 2
			default: // T047: invalid escape - keep literal
				builder.WriteByte(input[i])
				i++
			}
		} else {
			builder.WriteByte(input[i])
			i++
		}
	}

	return builder.String(), false
}
