package cat

import "testing"

// T009 [P] [US1] TestFormatLine_NoOptions - plain text, no formatting
func TestFormatLine_NoOptions(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{}

	got := formatter.FormatLine("Hello, World!", 1, false, opts)
	want := "Hello, World!"

	if got != want {
		t.Errorf("FormatLine() = %q, want %q", got, want)
	}
}

func TestFormatLine_EmptyLine(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{}

	got := formatter.FormatLine("", 1, true, opts)
	want := ""

	if got != want {
		t.Errorf("FormatLine() = %q, want %q", got, want)
	}
}

// T031 [P] [US3] TestFormatLine_NumberAll - line numbering for all lines
func TestFormatLine_NumberAll(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{NumberAll: true}

	got := formatter.FormatLine("Hello", 42, false, opts)
	want := "    42  Hello"

	if got != want {
		t.Errorf("FormatLine() with NumberAll = %q, want %q", got, want)
	}
}

// T032 [P] [US3] TestFormatLine_NumberAll_EmptyLine - numbering includes empty lines
func TestFormatLine_NumberAll_EmptyLine(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{NumberAll: true}

	got := formatter.FormatLine("", 5, true, opts)
	want := "     5  "

	if got != want {
		t.Errorf("FormatLine() with NumberAll on empty line = %q, want %q", got, want)
	}
}

// T033 [P] [US3] TestFormatLine_NumberAll_Overflow - line numbers wrap at 1,000,000
func TestFormatLine_NumberAll_Overflow(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{NumberAll: true}

	// Test line 1,000,000 (wraps to 0)
	got := formatter.FormatLine("overflow", 1000000, false, opts)
	want := "     0  overflow"

	if got != want {
		t.Errorf("FormatLine() at line 1,000,000 = %q, want %q", got, want)
	}

	// Test line 1,000,001 (wraps to 1)
	got2 := formatter.FormatLine("wrapped", 1000001, false, opts)
	want2 := "     1  wrapped"

	if got2 != want2 {
		t.Errorf("FormatLine() at line 1,000,001 = %q, want %q", got2, want2)
	}
}

// T043 [P] [US4] TestFormatLine_NumberNonBlank_EmptyLine - empty lines not numbered
func TestFormatLine_NumberNonBlank_EmptyLine(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{NumberNonBlank: true}

	got := formatter.FormatLine("", 5, true, opts)
	want := ""

	if got != want {
		t.Errorf("FormatLine() with NumberNonBlank on empty line = %q, want %q", got, want)
	}
}

// T044 [P] [US4] TestFormatLine_NumberNonBlank_NonEmptyLine - non-empty lines are numbered
func TestFormatLine_NumberNonBlank_NonEmptyLine(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{NumberNonBlank: true}

	got := formatter.FormatLine("Content", 10, false, opts)
	want := "    10  Content"

	if got != want {
		t.Errorf("FormatLine() with NumberNonBlank = %q, want %q", got, want)
	}
}

// T051 [P] [US5] TestFormatLine_ShowEnds - line ends with $
func TestFormatLine_ShowEnds(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{ShowEnds: true}

	got := formatter.FormatLine("Hello", 1, false, opts)
	want := "Hello$"

	if got != want {
		t.Errorf("FormatLine() with ShowEnds = %q, want %q", got, want)
	}
}

// T052 [P] [US5] TestFormatLine_ShowEnds_EmptyLine - empty line shows $
func TestFormatLine_ShowEnds_EmptyLine(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{ShowEnds: true}

	got := formatter.FormatLine("", 1, true, opts)
	want := "$"

	if got != want {
		t.Errorf("FormatLine() with ShowEnds on empty line = %q, want %q", got, want)
	}
}

// T057 [P] [US6] TestFormatLine_ShowTabs - tabs shown as ^I
func TestFormatLine_ShowTabs(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{ShowTabs: true}

	got := formatter.FormatLine("Hello\tWorld", 1, false, opts)
	want := "Hello^IWorld"

	if got != want {
		t.Errorf("FormatLine() with ShowTabs = %q, want %q", got, want)
	}
}

// T058 [P] [US6] TestFormatLine_ShowTabs_MultipleTabs - multiple tabs converted
func TestFormatLine_ShowTabs_MultipleTabs(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{ShowTabs: true}

	got := formatter.FormatLine("\tA\t\tB\t", 1, false, opts)
	want := "^IA^I^IB^I"

	if got != want {
		t.Errorf("FormatLine() with ShowTabs = %q, want %q", got, want)
	}
}

// T063 [P] [US7] TestFormatLine_ShowNonPrinting_ControlChars - control chars visible
func TestFormatLine_ShowNonPrinting_ControlChars(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{ShowNonPrinting: true}

	// ASCII 7 (BEL) -> ^G, ASCII 27 (ESC) -> ^[
	got := formatter.FormatLine("Hello\x07World\x1B", 1, false, opts)
	want := "Hello^GWorld^["

	if got != want {
		t.Errorf("FormatLine() with ShowNonPrinting = %q, want %q", got, want)
	}
}

// T064 [P] [US7] TestFormatLine_ShowNonPrinting_DEL - DEL (127) shown as ^?
func TestFormatLine_ShowNonPrinting_DEL(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{ShowNonPrinting: true}

	got := formatter.FormatLine("Test\x7FEnd", 1, false, opts)
	want := "Test^?End"

	if got != want {
		t.Errorf("FormatLine() with ShowNonPrinting = %q, want %q", got, want)
	}
}

// T065 [P] [US7] TestFormatLine_ShowNonPrinting_NoControlChars - normal text unchanged
func TestFormatLine_ShowNonPrinting_NoControlChars(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{ShowNonPrinting: true}

	got := formatter.FormatLine("Normal Text", 1, false, opts)
	want := "Normal Text"

	if got != want {
		t.Errorf("FormatLine() with ShowNonPrinting = %q, want %q", got, want)
	}
}

// T071 [P] [US8] TestFormatLine_AllOptions - combination of all formatting
func TestFormatLine_AllOptions(t *testing.T) {
	formatter := NewDefaultFormatter()
	opts := Options{
		NumberAll:       true,
		ShowEnds:        true,
		ShowTabs:        true,
		ShowNonPrinting: true,
	}

	got := formatter.FormatLine("A\tB\x07C", 3, false, opts)
	want := "     3  A^IB^GC$"

	if got != want {
		t.Errorf("FormatLine() with all options = %q, want %q", got, want)
	}
}
