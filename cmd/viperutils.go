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
