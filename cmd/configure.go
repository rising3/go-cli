package cmd

import (
	"path/filepath"

	"github.com/rising3/go-cli/internal/cmd/configure"
	"github.com/rising3/go-cli/internal/editor"
	"github.com/spf13/cobra"
)

var cfgForce bool
var cfgEdit bool
var cfgNoWait bool

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().BoolVar(&cfgForce, "force", false, "overwrite existing config")
	configureCmd.Flags().BoolVar(&cfgEdit, "edit", false, "edit the created file in $EDITOR")
	configureCmd.Flags().BoolVar(&cfgNoWait, "no-wait", false, "do not wait for editor to exit")
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create a scaffold config file based on Config struct",
	RunE: func(cmd *cobra.Command, args []string) error {
		// T036: Determine target path
		dir := GetConfigPath()
		cfgName := DefaultProfile
		if profile != "" {
			cfgName = profile
		}
		target := filepath.Join(dir, GetConfigFile(cfgName))

		// T037-T039: Build ConfigureOptions
		opts := configure.ConfigureOptions{
			Force:            cfgForce,
			Edit:             cfgEdit,
			NoWait:           cfgNoWait,
			Data:             BuildEffectiveConfig(),
			Format:           CliConfigType,
			Output:           cmd.OutOrStdout(),
			ErrOutput:        cmd.ErrOrStderr(),
			EditorLookup:     func() (string, []string, error) { return editor.GetEditor() },
			EditorShouldWait: func(string, []string) bool { return !cfgNoWait },
		}

		// T040: Call internal function
		return configure.ConfigureFunc(target, opts)
	},
}
