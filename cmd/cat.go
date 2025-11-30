package cmd

import (
	"github.com/rising3/go-cli/internal/cmd/cat"
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat [flags] [file...]",
	Short: "Concatenate files and print on the standard output",
	Long: `Concatenate FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

Examples:
  mycli cat file.txt           # Display file content
  mycli cat file1.txt file2.txt  # Concatenate multiple files
  echo "test" | mycli cat       # Read from stdin
  mycli cat -n file.txt         # Number all lines
  mycli cat -b file.txt         # Number non-empty lines
  mycli cat -E file.txt         # Show line ends with $
  mycli cat -T file.txt         # Show tabs as ^I
  mycli cat -v file.txt         # Show control characters
  mycli cat -A file.txt         # Show all (equivalent to -vET)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts, err := cat.NewOptions(cmd)
		if err != nil {
			return err
		}

		return cat.CatFunc(args, opts)
	},
}

func init() {
	rootCmd.AddCommand(catCmd)

	catCmd.Flags().BoolP("number", "n", false, "number all output lines")
	catCmd.Flags().BoolP("number-nonblank", "b", false, "number nonempty output lines")
	catCmd.Flags().BoolP("show-ends", "E", false, "display $ at end of each line")
	catCmd.Flags().BoolP("show-tabs", "T", false, "display TAB characters as ^I")
	catCmd.Flags().BoolP("show-nonprinting", "v", false, "use ^ and M- notation")
	catCmd.Flags().BoolP("show-all", "A", false, "equivalent to -vET")
}
