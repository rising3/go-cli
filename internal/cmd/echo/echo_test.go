package echo

import (
	"bytes"
	"testing"
)

func TestEchoOptions_Initialization(t *testing.T) {
	tests := []struct {
		name string
		opts EchoOptions
		want EchoOptions
	}{
		{
			name: "default values",
			opts: EchoOptions{},
			want: EchoOptions{
				SuppressNewline:  false,
				InterpretEscapes: false,
				Verbose:          false,
				Args:             nil,
			},
		},
		{
			name: "with flags set",
			opts: EchoOptions{
				SuppressNewline:  true,
				InterpretEscapes: true,
				Verbose:          true,
				Args:             []string{"test"},
			},
			want: EchoOptions{
				SuppressNewline:  true,
				InterpretEscapes: true,
				Verbose:          true,
				Args:             []string{"test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opts.SuppressNewline != tt.want.SuppressNewline {
				t.Errorf("SuppressNewline = %v, want %v", tt.opts.SuppressNewline, tt.want.SuppressNewline)
			}
			if tt.opts.InterpretEscapes != tt.want.InterpretEscapes {
				t.Errorf("InterpretEscapes = %v, want %v", tt.opts.InterpretEscapes, tt.want.InterpretEscapes)
			}
			if tt.opts.Verbose != tt.want.Verbose {
				t.Errorf("Verbose = %v, want %v", tt.opts.Verbose, tt.want.Verbose)
			}
			if len(tt.opts.Args) != len(tt.want.Args) {
				t.Errorf("Args length = %d, want %d", len(tt.opts.Args), len(tt.want.Args))
			}
		})
	}
}

// T021: GenerateOutput unit test
func TestGenerateOutput(t *testing.T) {
	tests := []struct {
		name                string
		opts                EchoOptions
		wantOutput          string
		wantSuppressNewline bool
	}{
		{
			name:                "single argument",
			opts:                EchoOptions{Args: []string{"Hello"}},
			wantOutput:          "Hello",
			wantSuppressNewline: false,
		},
		{
			name:                "multiple arguments",
			opts:                EchoOptions{Args: []string{"A", "B", "C"}},
			wantOutput:          "A B C",
			wantSuppressNewline: false,
		},
		{
			name:                "no arguments",
			opts:                EchoOptions{Args: []string{}},
			wantOutput:          "",
			wantSuppressNewline: false,
		},
		{
			name:                "with escape sequences",
			opts:                EchoOptions{InterpretEscapes: true, Args: []string{"Hello\\nWorld"}},
			wantOutput:          "Hello\nWorld",
			wantSuppressNewline: false,
		},
		{
			name:                "with suppress flag",
			opts:                EchoOptions{SuppressNewline: true, Args: []string{"Hello"}},
			wantOutput:          "Hello",
			wantSuppressNewline: true,
		},
		{
			name:                "with \\c escape",
			opts:                EchoOptions{InterpretEscapes: true, Args: []string{"Hello\\cWorld"}},
			wantOutput:          "Hello",
			wantSuppressNewline: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotSuppress := GenerateOutput(tt.opts)
			if gotOutput != tt.wantOutput {
				t.Errorf("GenerateOutput() output = %q, want %q", gotOutput, tt.wantOutput)
			}
			if gotSuppress != tt.wantSuppressNewline {
				t.Errorf("GenerateOutput() suppress = %v, want %v", gotSuppress, tt.wantSuppressNewline)
			}
		})
	}
}

// T021: WriteOutput unit test
func TestWriteOutput(t *testing.T) {
	tests := []struct {
		name            string
		output          string
		suppressNewline bool
		want            string
	}{
		{
			name:            "with newline",
			output:          "Hello",
			suppressNewline: false,
			want:            "Hello\n",
		},
		{
			name:            "without newline",
			output:          "Hello",
			suppressNewline: true,
			want:            "Hello",
		},
		{
			name:            "empty with newline",
			output:          "",
			suppressNewline: false,
			want:            "\n",
		},
		{
			name:            "empty without newline",
			output:          "",
			suppressNewline: true,
			want:            "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := WriteOutput(buf, tt.output, tt.suppressNewline)
			if err != nil {
				t.Fatalf("WriteOutput() error = %v", err)
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("WriteOutput() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Refactored Echo function test
func TestEcho(t *testing.T) {
	tests := []struct {
		name string
		opts EchoOptions
		want string
	}{
		{
			name: "basic output",
			opts: EchoOptions{
				Args:   []string{"Hello", "World"},
				Output: nil, // will be set in test
			},
			want: "Hello World\n",
		},
		{
			name: "with suppress newline",
			opts: EchoOptions{
				Args:            []string{"Hello"},
				SuppressNewline: true,
				Output:          nil,
			},
			want: "Hello",
		},
		{
			name: "with escape sequences",
			opts: EchoOptions{
				Args:             []string{"Line1\\nLine2"},
				InterpretEscapes: true,
				Output:           nil,
			},
			want: "Line1\nLine2\n",
		},
		{
			name: "combined flags",
			opts: EchoOptions{
				Args:             []string{"Hello\\tWorld"},
				SuppressNewline:  true,
				InterpretEscapes: true,
				Output:           nil,
			},
			want: "Hello\tWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tt.opts.Output = buf
			tt.opts.ErrOutput = new(bytes.Buffer)

			err := Echo(tt.opts)
			if err != nil {
				t.Fatalf("Echo() error = %v", err)
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("Echo() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Test EchoFunc indirection for testability
func TestEchoFunc_Indirection(t *testing.T) {
	// Save original
	originalEchoFunc := EchoFunc
	defer func() { EchoFunc = originalEchoFunc }()

	// Mock implementation
	called := false
	EchoFunc = func(opts EchoOptions) error {
		called = true
		return nil
	}

	// Test that the mock is called
	_ = EchoFunc(EchoOptions{})
	if !called {
		t.Error("EchoFunc was not called through indirection")
	}
}

// Test ProcessEscapesFunc indirection for testability
func TestProcessEscapesFunc_Indirection(t *testing.T) {
	// Save original
	originalProcessEscapesFunc := ProcessEscapesFunc
	defer func() { ProcessEscapesFunc = originalProcessEscapesFunc }()

	// Mock implementation
	ProcessEscapesFunc = func(input string) (string, bool) {
		return "mocked", true
	}

	// Test with InterpretEscapes=true to trigger ProcessEscapesFunc
	buf := new(bytes.Buffer)
	opts := EchoOptions{
		Args:             []string{"test"},
		InterpretEscapes: true,
		Output:           buf,
		ErrOutput:        new(bytes.Buffer),
	}

	_ = Echo(opts)
	got := buf.String()
	if got != "mocked" {
		t.Errorf("ProcessEscapesFunc mock not used, got %q", got)
	}
}
