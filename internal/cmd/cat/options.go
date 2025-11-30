package cat

import "github.com/spf13/cobra"

// Options holds the formatting options for cat command
type Options struct {
	// NumberAll numbers all output lines (-n flag)
	NumberAll bool

	// NumberNonBlank numbers nonempty output lines (-b flag)
	NumberNonBlank bool

	// ShowEnds displays $ at end of each line (-E flag)
	ShowEnds bool

	// ShowTabs displays TAB characters as ^I (-T flag)
	ShowTabs bool

	// ShowNonPrinting uses ^ and M- notation for control characters (-v flag)
	ShowNonPrinting bool
}

// NewOptions creates Options from Cobra command flags
func NewOptions(cmd *cobra.Command) (Options, error) {
	opts := Options{}

	numberAll, _ := cmd.Flags().GetBool("number")
	numberNonBlank, _ := cmd.Flags().GetBool("number-nonblank")
	showEnds, _ := cmd.Flags().GetBool("show-ends")
	showTabs, _ := cmd.Flags().GetBool("show-tabs")
	showNonPrinting, _ := cmd.Flags().GetBool("show-nonprinting")
	showAll, _ := cmd.Flags().GetBool("show-all")

	// -A flag expands to -vET
	if showAll {
		showNonPrinting = true
		showEnds = true
		showTabs = true
	}

	// T049 [US4] Handle -n and -b conflict: -b takes precedence
	// If both are set, NumberNonBlank wins and NumberAll is disabled
	if numberAll && numberNonBlank {
		opts.NumberNonBlank = true
		opts.NumberAll = false
	} else {
		opts.NumberAll = numberAll
		opts.NumberNonBlank = numberNonBlank
	}

	opts.ShowEnds = showEnds
	opts.ShowTabs = showTabs
	opts.ShowNonPrinting = showNonPrinting

	return opts, nil
}
