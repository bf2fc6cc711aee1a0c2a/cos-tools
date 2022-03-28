package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/antihax/optional"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/api/public"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/config"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/service"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewNamespaceLsCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "ls",
		Short: "ls",
		Long:  `ls`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := service.NewClient(cmd.Context())
			if err != nil {
				return err
			}

			page := viper.GetInt32("rhoc.namespace.ls.page")
			all := viper.GetBool("rhoc.namespace.ls.all")
			output := viper.GetString("rhoc.namespace.ls.output")

			namespaces := public.ConnectorNamespaceList{
				Kind:  "ConnectorNamespaceList",
				Items: make([]public.ConnectorNamespace, 0),
				Total: 0,
				Size:  0,
			}

			for i := page; i == page || all; i++ {
				list, _, err := client.ConnectorNamespacesApi.ListConnectorNamespaces(cmd.Context(), &public.ListConnectorNamespacesOpts{
					Page:   optional.NewString(strconv.Itoa(int(i))),
					Size:   config.GetOptionalString("rhoc.namespace.ls.size"),
					Search: config.GetOptionalString("rhoc.namespace.ls.search"),
				})

				if err != nil {
					return err
				}
				if len(list.Items) == 0 {
					break
				}

				namespaces.Items = append(namespaces.Items, list.Items...)
				namespaces.Size = int32(len(namespaces.Items))
				namespaces.Total = list.Total
			}

			switch output {
			case config.OutputDefault, config.OutputTable:
				table := tablewriter.NewWriter(cmd.OutOrStdout())
				table.SetBorder(true)
				table.SetHeader([]string{"ID", "Name", "ClusterID", "Owner", "CreatedAt", "Tenant Kind", "Tenant ID", "State"})

				for _, ns := range namespaces.Items {
					table.Append([]string{
						ns.Id,
						ns.Name,
						ns.ClusterId,
						ns.Owner,
						ns.CreatedAt.String(),
						string(ns.Tenant.Kind),
						ns.Tenant.Id,
						string(ns.Status.State),
					})
				}

				table.SetCaption(true, fmt.Sprintf("Items: %d, Total: %d", len(namespaces.Items), namespaces.Total))
				table.Render()
			case config.OutputJson:
				bytes, err := json.Marshal(namespaces)
				if err != nil {
					return err
				}

				fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			case config.OutputYaml:
				bytes, err := yaml.Marshal(namespaces)
				if err != nil {
					return err
				}

				fmt.Fprintln(cmd.OutOrStdout(), string(bytes))
			}

			return nil
		},
	}

	cmd.Flags().Int32("page", 1, "page")
	cmd.Flags().Int32("size", math.MaxInt32, "size")
	cmd.Flags().String("search", "", "search")
	cmd.Flags().BoolP("all", "a", false, "all")
	cmd.Flags().StringP("output", "o", "", "output")

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		viper.BindPFlag("rhoc.namespace.ls."+flag.Name, flag)
	})

	return &cmd
}
