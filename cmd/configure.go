package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	ieditor "github.com/rising3/go-cli/internal/editor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgForce bool
var cfgEdit bool

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.Flags().BoolVar(&cfgForce, "force", false, "overwrite existing config")
	configureCmd.Flags().BoolVar(&cfgEdit, "edit", false, "edit the created file in $EDITOR")
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create a scaffold config file based on Config struct",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := GetConfigPath()
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		cfgName := DefaultProfile
		if profile != "" {
			cfgName = profile
		}
		target := filepath.Join(dir, GetConfigFile(cfgName))

		if _, err := os.Stat(target); err == nil && !cfgForce {
			// Config already exists; do not initialize/overwrite when --force is not provided.
			// Continue processing so options like --edit still work against the existing file.
			fmt.Fprintln(os.Stderr, "Config already exists, skipping initialization:", target)
		}

		// Use a fresh viper to build scaffold content
		vp := viper.New()
		vp.SetConfigType(CliConfigType)

		// If CliConfig has existing values, use them as defaults in scaffold
		vp.Set("client-id", CliConfig.ClientID)
		vp.Set("client-secret", CliConfig.ClientSecret)

		// Only write the scaffold when forcing or when the file does not exist.
		if _, err := os.Stat(target); os.IsNotExist(err) || cfgForce {
			if cfgForce {
				// best-effort remove before write
				_ = os.Remove(target)
			}

			if err := vp.WriteConfigAs(target); err != nil {
				return err
			}

			fmt.Fprintln(os.Stderr, "Wrote config:", target)
		}

		// Edit in editor only when requested via --edit
		if cfgEdit {
			ed, edArgs, err := ieditor.GetEditor()
			if err != nil {
				fmt.Fprintln(os.Stderr, "No editor found:", err)
			} else {
				args := append(edArgs, target)
				cmd := exec.Command(ed, args...)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Start(); err != nil {
					fmt.Fprintln(os.Stderr, "Failed to start editor:", err)
				}
				if err := cmd.Wait(); err != nil {
					fmt.Fprintln(os.Stderr, "Editor exited with error:", err)
				}
			}
		}
		return nil
	},
}
