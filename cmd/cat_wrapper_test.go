package cmd

import (
	"bytes"
	"testing"

	"github.com/rising3/go-cli/internal/cmd/cat"
	"github.com/spf13/cobra"
)

func TestCatWrapperCallsInternal(t *testing.T) {
	// capture arguments passed to CatFunc
	var calledOptions cat.Options
	var calledFilenames []string

	// stub internal implementation
	old := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledFilenames = filenames
		calledOptions = opts
		return nil
	}
	defer func() { cat.CatFunc = old }()

	// Create command with flags
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	// Add flags to the test command (copy from catCmd)
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	// Set flags
	if err := cmd.Flags().Set("number", "true"); err != nil {
		t.Fatalf("failed to set number flag: %v", err)
	}
	if err := cmd.Flags().Set("show-ends", "true"); err != nil {
		t.Fatalf("failed to set show-ends flag: %v", err)
	}

	// Execute
	args := []string{"file1.txt", "file2.txt"}
	if err := catCmd.RunE(cmd, args); err != nil {
		t.Fatalf("cat RunE failed: %v", err)
	}

	// Verify internal CatFunc was called with correct options
	if !calledOptions.NumberAll {
		t.Errorf("expected NumberAll true, got false")
	}
	if calledOptions.NumberNonBlank {
		t.Errorf("expected NumberNonBlank false, got true")
	}
	if !calledOptions.ShowEnds {
		t.Errorf("expected ShowEnds true, got false")
	}
	if calledOptions.ShowTabs {
		t.Errorf("expected ShowTabs false, got true")
	}
	if calledOptions.ShowNonPrinting {
		t.Errorf("expected ShowNonPrinting false, got true")
	}

	// Verify filenames were passed correctly
	if len(calledFilenames) != 2 || calledFilenames[0] != "file1.txt" || calledFilenames[1] != "file2.txt" {
		t.Errorf("expected filenames [file1.txt file2.txt], got %v", calledFilenames)
	}
}

func TestCatWrapper_DefaultFlags(t *testing.T) {
	var calledOptions cat.Options

	// stub internal implementation
	old := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledOptions = opts
		return nil
	}
	defer func() { cat.CatFunc = old }()

	// Create command with flags at default values
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	// Execute with default flags
	args := []string{"test.txt"}
	if err := catCmd.RunE(cmd, args); err != nil {
		t.Fatalf("cat RunE failed: %v", err)
	}

	// Verify all flags are false by default
	if calledOptions.NumberAll {
		t.Errorf("expected NumberAll false by default, got true")
	}
	if calledOptions.NumberNonBlank {
		t.Errorf("expected NumberNonBlank false by default, got true")
	}
	if calledOptions.ShowEnds {
		t.Errorf("expected ShowEnds false by default, got true")
	}
	if calledOptions.ShowTabs {
		t.Errorf("expected ShowTabs false by default, got true")
	}
	if calledOptions.ShowNonPrinting {
		t.Errorf("expected ShowNonPrinting false by default, got true")
	}
}

func TestCatWrapper_ShowAllFlag(t *testing.T) {
	var calledOptions cat.Options

	// stub internal implementation
	old := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledOptions = opts
		return nil
	}
	defer func() { cat.CatFunc = old }()

	// Create command with show-all flag enabled
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("show-all", "true"); err != nil {
		t.Fatalf("failed to set show-all flag: %v", err)
	}

	// Execute
	args := []string{"test.txt"}
	if err := catCmd.RunE(cmd, args); err != nil {
		t.Fatalf("cat RunE failed: %v", err)
	}

	// Verify -A flag expands to -vET
	if !calledOptions.ShowNonPrinting {
		t.Errorf("expected ShowNonPrinting true when -A is set, got false")
	}
	if !calledOptions.ShowEnds {
		t.Errorf("expected ShowEnds true when -A is set, got false")
	}
	if !calledOptions.ShowTabs {
		t.Errorf("expected ShowTabs true when -A is set, got false")
	}
}

func TestCatWrapper_NumberNonBlankFlag(t *testing.T) {
	var calledOptions cat.Options

	// stub internal implementation
	old := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledOptions = opts
		return nil
	}
	defer func() { cat.CatFunc = old }()

	// Create command with number-nonblank flag enabled
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("number-nonblank", "true"); err != nil {
		t.Fatalf("failed to set number-nonblank flag: %v", err)
	}

	// Execute
	args := []string{"test.txt"}
	if err := catCmd.RunE(cmd, args); err != nil {
		t.Fatalf("cat RunE failed: %v", err)
	}

	// Verify number-nonblank flag was passed correctly
	if !calledOptions.NumberNonBlank {
		t.Errorf("expected NumberNonBlank true, got false")
	}
	if calledOptions.NumberAll {
		t.Errorf("expected NumberAll false when -b is set, got true")
	}
}

func TestCatWrapper_EmptyArgs(t *testing.T) {
	var calledFilenames []string

	// stub internal implementation
	old := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledFilenames = filenames
		return nil
	}
	defer func() { cat.CatFunc = old }()

	// Create command
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	// Execute with empty args (should read from stdin)
	if err := catCmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("cat RunE failed: %v", err)
	}

	// Verify empty args were passed (stdin mode)
	if len(calledFilenames) != 0 {
		t.Errorf("expected empty filenames for stdin, got %v", calledFilenames)
	}
}

func TestCatWrapper_MultipleFlags(t *testing.T) {
	var calledOptions cat.Options

	// stub internal implementation
	old := cat.CatFunc
	cat.CatFunc = func(filenames []string, opts cat.Options) error {
		calledOptions = opts
		return nil
	}
	defer func() { cat.CatFunc = old }()

	// Create command with multiple flags enabled
	cmd := &cobra.Command{}
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))

	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	// Set multiple flags
	if err := cmd.Flags().Set("number", "true"); err != nil {
		t.Fatalf("failed to set number flag: %v", err)
	}
	if err := cmd.Flags().Set("show-tabs", "true"); err != nil {
		t.Fatalf("failed to set show-tabs flag: %v", err)
	}
	if err := cmd.Flags().Set("show-nonprinting", "true"); err != nil {
		t.Fatalf("failed to set show-nonprinting flag: %v", err)
	}

	// Execute
	args := []string{"test.txt"}
	if err := catCmd.RunE(cmd, args); err != nil {
		t.Fatalf("cat RunE failed: %v", err)
	}

	// Verify all flags were passed correctly
	if !calledOptions.NumberAll {
		t.Errorf("expected NumberAll true, got false")
	}
	if !calledOptions.ShowTabs {
		t.Errorf("expected ShowTabs true, got false")
	}
	if !calledOptions.ShowNonPrinting {
		t.Errorf("expected ShowNonPrinting true, got false")
	}
	if calledOptions.ShowEnds {
		t.Errorf("expected ShowEnds false, got true")
	}
}
