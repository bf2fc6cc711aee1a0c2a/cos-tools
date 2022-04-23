package root

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/clusters"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/commands"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/config"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/connectors"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/deployments"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/namespaces"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/completion"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/login"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/logout"
	"github.com/redhat-developer/app-services-cli/pkg/cmd/whoami"
	"github.com/redhat-developer/app-services-cli/pkg/core/cmdutil/flagutil"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewRootCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: false,
		Use:           "rhoc",
		Short:         "rhoc",
		Long:          "",
		Example:       "",
	}

	fs := cmd.PersistentFlags()
	flagutil.VerboseFlag(fs)

	// this flag comes out of the box, but has its own basic usage text, so this overrides that
	var help bool
	fs.BoolVarP(&help, "help", "h", false, "Prints help information")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	commands.Bind(
		cmd,
		login.NewLoginCmd(f),
		logout.NewLogoutCommand(f),
		whoami.NewWhoAmICmd(f),
		completion.NewCompletionCommand(f),
		config.NewConfigCommand(f),
		namespaces.NewNamespacesCommand(f),
		connectors.NewConnectorsCommand(f),
		deployments.NewDeploymentsCommand(f),
		clusters.NeClustersCommand(f))

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		f.Logger.SetDebug(flagutil.DebugEnabled())

		configName := os.Getenv("RHOC_CONFIG_NAME")
		if configName == "" {
			configName = "rhoc"
		}

		v := viper.New()
		v.SetConfigName(configName)

		configPath := os.Getenv("RHOC_CONFIG_PATH")
		if configPath != "" {
			// if a specific config path is set, don't add
			// default locations
			v.AddConfigPath(configPath)
		} else {
			v.AddConfigPath(".")
			v.AddConfigPath(".rhoc")
			v.AddConfigPath("$HOME/.rhoc")
		}

		v.SetEnvPrefix("RHOC")
		v.SetEnvKeyReplacer(strings.NewReplacer(
			".", "_",
			"-", "_",
		))

		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				panic(err)
			}
		}

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if !f.Changed && v.IsSet(f.Name) {
				val := v.Get(f.Name)
				if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
					panic(fmt.Errorf("error setting flag %s, %s", f.Name, err))
				}
			}
		})

		return nil
	}

	cmd.InitDefaultHelpCmd()

	return cmd
}
