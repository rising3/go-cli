package cmd

import (
	"path/filepath"

	internalcmd "github.com/rising3/go-cli/internal/cmd"
	editor "github.com/rising3/go-cli/internal/editor"
	stdio "github.com/rising3/go-cli/internal/stdio"
	"github.com/spf13/cobra"
)

var cfgForce bool
var cfgEdit bool
var cfgNoWait bool
var cfgInput string
var cfgOutput string

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().BoolVar(&cfgForce, "force", false, "overwrite existing config")
	configureCmd.Flags().BoolVar(&cfgEdit, "edit", false, "edit the created file in $EDITOR")
	configureCmd.Flags().BoolVar(&cfgNoWait, "no-wait", false, "do not wait for editor to exit")
	configureCmd.Flags().StringVar(&cfgInput, "input", "", "read input from path ('-' for stdin)")
	configureCmd.Flags().StringVar(&cfgOutput, "output", "", "write output to path ('-' for stdout)")
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create a scaffold config file based on Config struct",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := GetConfigPath()
		cfgName := DefaultProfile
		if profile != "" {
			cfgName = profile
		}
		target := filepath.Join(dir, GetConfigFile(cfgName))

		opts := internalcmd.ConfigureOptions{
			Force:            cfgForce,
			Edit:             cfgEdit,
			Data:             BuildEffectiveConfig(),
			Format:           CliConfigType,
			Streams:          stdio.NewDefault(),
			EditorLookup:     func() (string, []string, error) { return editor.GetEditor() },
			EditorShouldWait: func(string, []string) bool { return !cfgNoWait },
		}

		return internalcmd.ConfigureFunc(target, opts)
	},
}
