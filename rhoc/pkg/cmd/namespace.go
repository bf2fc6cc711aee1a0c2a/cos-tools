package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewNamespaceCommand() *cobra.Command {
	var cmd = cobra.Command{
		Use:     "namespace",
		Short:   "namespace",
		Long:    `namespace`,
		Aliases: []string{"ns"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag("rhoc.namespace."+flag.Name, flag)
	})
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag("rhoc.namespace."+flag.Name, flag)
	})

	cmd.AddCommand(NewNamespaceLsCommand())
	cmd.AddCommand(NewNamespaceGetCommand())

	return &cmd
}
