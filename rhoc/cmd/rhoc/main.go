package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	DefaultConfigName     = "config"
	DefaultConfigLocation = DefaultConfigName + ".yaml"
)

func main() {
	cobra.OnInitialize(func() {
		configName := os.Getenv("RHOC_CONFIG_NAME")
		if configName == "" {
			configName = DefaultConfigName
		}

		viper.SetConfigName(configName)

		configPath := os.Getenv("RHOC_CONFIG_PATH")
		if configPath != "" {
			viper.AddConfigPath(configPath)
		} else {
			viper.AddConfigPath(".")
			viper.AddConfigPath(".rhoc")
			viper.AddConfigPath("$HOME/.rhoc")
		}

		viper.AutomaticEnv()
		viper.SetEnvPrefix("RHOC")
		viper.SetEnvKeyReplacer(strings.NewReplacer(
			".", "_",
			"-", "_",
		))

		_ = viper.ReadInConfig()
	})

	var rhoc = cobra.Command{
		Use:   "rhoc",
		Short: "rhoc",
		Long:  `rhoc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	rhoc.PersistentFlags().String("api-token", "", "api-token")
	rhoc.PersistentFlags().String("api-url", "", "api-url")

	rhoc.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag("rhoc."+flag.Name, flag)
	})

	rhoc.AddCommand(cmd.NewNamespaceCommand())

	if err := rhoc.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
