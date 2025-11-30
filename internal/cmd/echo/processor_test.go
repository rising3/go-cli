package echo

import "testing"

// T025: \n (newline) escape test
func TestProcessEscapes_Newline(t *testing.T) {
	input := `Hello\nWorld`
	want := "Hello\nWorld"
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T026: \t (tab) escape test
func TestProcessEscapes_Tab(t *testing.T) {
	input := `Hello\tWorld`
	want := "Hello\tWorld"
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T027: \\ (backslash) escape test
func TestProcessEscapes_Backslash(t *testing.T) {
	input := `Hello\\World`
	want := `Hello\World`
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T028: \" (double quote) escape test
func TestProcessEscapes_DoubleQuote(t *testing.T) {
	input := `Hello\"World`
	want := `Hello"World`
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T029: \a (bell/alert) escape test
func TestProcessEscapes_Alert(t *testing.T) {
	input := `Hello\aWorld`
	want := "Hello\aWorld"
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T030: \b (backspace) escape test
func TestProcessEscapes_Backspace(t *testing.T) {
	input := `Hello\bWorld`
	want := "Hello\bWorld"
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T031: \c (suppress further output) escape test
func TestProcessEscapes_SuppressOutput(t *testing.T) {
	input := `Hello\cWorld`
	want := "Hello"
	got, suppressNewline := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
	if !suppressNewline {
		t.Errorf("ProcessEscapes(%q) suppressNewline = false, want true", input)
	}
}

// T032: \r (carriage return) escape test
func TestProcessEscapes_CarriageReturn(t *testing.T) {
	input := `Hello\rWorld`
	want := "Hello\rWorld"
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T033: \v (vertical tab) escape test
func TestProcessEscapes_VerticalTab(t *testing.T) {
	input := `Hello\vWorld`
	want := "Hello\vWorld"
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}

// T034: Invalid escape sequence test
func TestProcessEscapes_InvalidEscape(t *testing.T) {
	input := `Hello\zWorld`
	want := `Hello\zWorld` // Should remain literal
	got, _ := ProcessEscapes(input)
	if got != want {
		t.Errorf("ProcessEscapes(%q) = %q, want %q", input, got, want)
	}
}
