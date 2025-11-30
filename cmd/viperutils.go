package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewViper(profile string) *viper.Viper {
	vp := viper.New()
	InitViper(vp, profile)
	return vp
}

func InitViper(vp *viper.Viper, profile string) {
	vp.SetConfigType(CliConfigType)
	vp.AddConfigPath(GetConfigPath())
	vp.SetConfigName(profile)
}

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	return filepath.Join(home, CliConfigBase, CliName)
}

func GetConfigFile(profile string) string {
	if profile == DefaultProfile {
		return DefaultProfile + "." + strings.ToLower(CliConfigType)

	} else {
		return profile + "." + strings.ToLower(CliConfigType)
	}
}

// BuildEffectiveConfig returns the effective configuration as a plain map.
// It returns a nested map with default values for all configuration fields,
// suitable for generating new configuration files via the configure command.
func BuildEffectiveConfig() map[string]interface{} {
	return map[string]interface{}{
		"client-id":     "",
		"client-secret": "",
		"common": map[string]interface{}{
			"var1": "",
			"var2": 123,
		},
		"hoge": map[string]interface{}{
			"fuga": "hello",
			"foo": map[string]interface{}{
				"bar": "hello",
			},
		},
	}
}
