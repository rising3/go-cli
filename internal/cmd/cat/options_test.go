package cat

import (
	"testing"

	"github.com/spf13/cobra"
)

// T034 [P] [US3] TestNewOptions_NumberFlag - verify -n flag sets NumberAll
func TestNewOptions_NumberFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("number", "true"); err != nil {
		t.Fatalf("failed to set number flag: %v", err)
	}

	opts, err := NewOptions(cmd)
	if err != nil {
		t.Fatalf("NewOptions() failed: %v", err)
	}

	if !opts.NumberAll {
		t.Errorf("Expected NumberAll true, got false")
	}
	if opts.NumberNonBlank {
		t.Errorf("Expected NumberNonBlank false, got true")
	}
}

// T045 [P] [US4] TestNewOptions_NumberNonBlankFlag - verify -b flag sets NumberNonBlank
func TestNewOptions_NumberNonBlankFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("number-nonblank", "true"); err != nil {
		t.Fatalf("failed to set number-nonblank flag: %v", err)
	}

	opts, err := NewOptions(cmd)
	if err != nil {
		t.Fatalf("NewOptions() failed: %v", err)
	}

	if !opts.NumberNonBlank {
		t.Errorf("Expected NumberNonBlank true, got false")
	}
	if opts.NumberAll {
		t.Errorf("Expected NumberAll false when -b is set, got true")
	}
}

// T046 [P] [US4] TestNewOptions_NumberConflict - both -n and -b: last one wins
func TestNewOptions_NumberConflict(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	// Set both flags (simulating -n -b)
	if err := cmd.Flags().Set("number", "true"); err != nil {
		t.Fatalf("failed to set number flag: %v", err)
	}
	if err := cmd.Flags().Set("number-nonblank", "true"); err != nil {
		t.Fatalf("failed to set number-nonblank flag: %v", err)
	}

	opts, err := NewOptions(cmd)
	if err != nil {
		t.Fatalf("NewOptions() failed: %v", err)
	}

	// When both are set, -b wins (NumberNonBlank=true, NumberAll=false)
	if !opts.NumberNonBlank {
		t.Errorf("Expected NumberNonBlank true when both -n and -b are set, got false")
	}
	if opts.NumberAll {
		t.Errorf("Expected NumberAll false when both -n and -b are set, got true")
	}
}

// T053 [P] [US5] TestNewOptions_ShowEndsFlag - verify -E flag sets ShowEnds
func TestNewOptions_ShowEndsFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("show-ends", "true"); err != nil {
		t.Fatalf("failed to set show-ends flag: %v", err)
	}

	opts, err := NewOptions(cmd)
	if err != nil {
		t.Fatalf("NewOptions() failed: %v", err)
	}

	if !opts.ShowEnds {
		t.Errorf("Expected ShowEnds true, got false")
	}
}

// T059 [P] [US6] TestNewOptions_ShowTabsFlag - verify -T flag sets ShowTabs
func TestNewOptions_ShowTabsFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("show-tabs", "true"); err != nil {
		t.Fatalf("failed to set show-tabs flag: %v", err)
	}

	opts, err := NewOptions(cmd)
	if err != nil {
		t.Fatalf("NewOptions() failed: %v", err)
	}

	if !opts.ShowTabs {
		t.Errorf("Expected ShowTabs true, got false")
	}
}

// T066 [P] [US7] TestNewOptions_ShowNonPrintingFlag - verify -v flag sets ShowNonPrinting
func TestNewOptions_ShowNonPrintingFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("show-nonprinting", "true"); err != nil {
		t.Fatalf("failed to set show-nonprinting flag: %v", err)
	}

	opts, err := NewOptions(cmd)
	if err != nil {
		t.Fatalf("NewOptions() failed: %v", err)
	}

	if !opts.ShowNonPrinting {
		t.Errorf("Expected ShowNonPrinting true, got false")
	}
}

// T072 [P] [US8] TestNewOptions_ShowAllFlag - verify -A flag is parsed
func TestNewOptions_ShowAllFlag(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().BoolP("number", "n", false, "number all output lines")
	cmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")

	if err := cmd.Flags().Set("show-all", "true"); err != nil {
		t.Fatalf("failed to set show-all flag: %v", err)
	}

	opts, err := NewOptions(cmd)
	if err != nil {
		t.Fatalf("NewOptions() failed: %v", err)
	}

	// -A should expand to -vET
	if !opts.ShowNonPrinting {
		t.Errorf("Expected ShowNonPrinting true when -A is set, got false")
	}
	if !opts.ShowEnds {
		t.Errorf("Expected ShowEnds true when -A is set, got false")
	}
	if !opts.ShowTabs {
		t.Errorf("Expected ShowTabs true when -A is set, got false")
	}
}

// T073 [P] [US8] TestNewOptions_ShowAll_EquivalentToVET - verify -A equals -vET
func TestNewOptions_ShowAll_EquivalentToVET(t *testing.T) {
	// Test -A
	cmdA := &cobra.Command{}
	cmdA.Flags().BoolP("number", "n", false, "number all output lines")
	cmdA.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmdA.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmdA.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmdA.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmdA.Flags().BoolP("show-all", "A", false, "equivalent to -vET")
	_ = cmdA.Flags().Set("show-all", "true")

	optsA, _ := NewOptions(cmdA)

	// Test -vET
	cmdVET := &cobra.Command{}
	cmdVET.Flags().BoolP("number", "n", false, "number all output lines")
	cmdVET.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	cmdVET.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	cmdVET.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	cmdVET.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	cmdVET.Flags().BoolP("show-all", "A", false, "equivalent to -vET")
	_ = cmdVET.Flags().Set("show-nonprinting", "true")
	_ = cmdVET.Flags().Set("show-ends", "true")
	_ = cmdVET.Flags().Set("show-tabs", "true")

	optsVET, _ := NewOptions(cmdVET)

	// Verify equivalence
	if optsA.ShowNonPrinting != optsVET.ShowNonPrinting {
		t.Errorf("ShowNonPrinting mismatch: -A=%v, -vET=%v", optsA.ShowNonPrinting, optsVET.ShowNonPrinting)
	}
	if optsA.ShowEnds != optsVET.ShowEnds {
		t.Errorf("ShowEnds mismatch: -A=%v, -vET=%v", optsA.ShowEnds, optsVET.ShowEnds)
	}
	if optsA.ShowTabs != optsVET.ShowTabs {
		t.Errorf("ShowTabs mismatch: -A=%v, -vET=%v", optsA.ShowTabs, optsVET.ShowTabs)
	}
}
