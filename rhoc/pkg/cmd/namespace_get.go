package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/config"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewNamespaceGetCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "get",
		Long:  `get`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := service.NewClient(cmd.Context())
			if err != nil {
				return err
			}

			output := viper.GetString("rhoc.namespace.get.output")
			id := viper.GetString("rhoc.namespace.get.id")

			namespace, _, err := client.ConnectorNamespacesApi.GetConnectorNamespace(cmd.Context(), id)
			if err != nil {
				return err
			}

			switch output {
			case config.OutputDefault, config.OutputTable:
				table := tablewriter.NewWriter(cmd.OutOrStdout())
				table.SetBorder(true)

				table.AppendBulk([][]string{
					{"Id", namespace.Id},
					{"Name", namespace.Name},
					{"ClusterId", namespace.ClusterId},
					{"Owner", namespace.Owner},
					{"CreatedAt", namespace.CreatedAt.String()},
					{"Tenant Kind", string(namespace.Tenant.Kind)},
					{"Tenant Id", namespace.Tenant.Id},
					{"State", string(namespace.Status.State)},
				})

				table.Render()
			case config.OutputJson:
				bytes, err := json.Marshal(namespace)
				if err != nil {
					return err
				}

				fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			case config.OutputYaml:
				bytes, err := yaml.Marshal(namespace)
				if err != nil {
					return err
				}

				fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			}

			return nil
		},
	}

	cmd.Flags().String("id", "", "id")
	cmd.Flags().StringP("output", "o", "", "output")
	cmd.MarkFlagRequired("id")

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag("rhoc.namespace.get."+flag.Name, flag)
	})

	return &cmd
}
