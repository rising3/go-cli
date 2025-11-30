package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	CliName        = "mycli"
	CliVersion     = "dev"
	CliConfigBase  = ".config"
	CliConfigType  = "yaml"
	DefaultProfile = "default"
)

type Config struct {
	ClientID     string       `mapstructure:"client-id"`
	ClientSecret string       `mapstructure:"client-secret"`
	Common       CommonConfig `mapstructure:"common"`
	Hoge         HogeConfig   `mapstructure:"hoge"`
}

type CommonConfig struct {
	Var1 string `mapstructure:"var1"`
	Var2 int    `mapstructure:"var2"`
}

type HogeConfig struct {
	Fuga string    `mapstructure:"fuga"`
	Foo  FooConfig `mapstructure:"foo"`
}

type FooConfig struct {
	Bar string `mapstructure:"bar"`
}

var CliConfig Config

var cfgFile string
var profile string

var rootCmd = &cobra.Command{
	Use:     CliName,
	Version: CliVersion,
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+filepath.Join(GetConfigPath(), GetConfigFile(DefaultProfile))+")")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "config profile (e.g. dev, prod)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	resolveEnvOverrides()

	if readExplicitConfig() {
	} else {
		readDefaultAndMergeProfile()
	}

	viper.SetEnvPrefix(strings.ToUpper(CliName))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&CliConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse configuration:", err)
	}
}

func resolveEnvOverrides() {
	if envCfg := os.Getenv("MYCLI_CONFIG"); envCfg != "" && cfgFile == "" {
		cfgFile = envCfg
	}
	if envProfile := os.Getenv("MYCLI_PROFILE"); envProfile != "" && profile == "" {
		profile = envProfile
	}
}

func readExplicitConfig() bool {
	if cfgFile == "" {
		return false
	}

	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Failed to read config file:", err)
	}
	return true
}

func readDefaultAndMergeProfile() {
	InitViper(viper.GetViper(), DefaultProfile)

	_ = viper.ReadInConfig()

	if profile == "" {
		return
	}

	vp := NewViper(profile)

	if err := vp.ReadInConfig(); err == nil {
		if err := viper.MergeConfigMap(vp.AllSettings()); err == nil {
			fmt.Fprintln(os.Stderr, "Merged profile config:", vp.ConfigFileUsed())
		}
	}
}
