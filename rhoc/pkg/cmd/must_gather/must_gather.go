package must_gather

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/collections"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/kubernetes"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/kubernetes/resources"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

type options struct {
	id        string
	logs      bool
	file      string
	resources []string
	gvrs      []schema.GroupVersionResource
	o         *cmdutil.OutputWriter
	f         *factory.Factory
}

func NewMustGatherCommand(f *factory.Factory) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)
	_ = printers.NewTypeSetter(scheme.Scheme).ToPrinter(&printers.YAMLPrinter{})

	opts := options{
		f:    f,
		logs: false,
	}

	cmd := &cobra.Command{
		Use:     "must-gather",
		Aliases: []string{"mg"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := configFlags.ToRESTConfig()
			if err != nil {
				return err
			}
			client, err := kubernetes.NewClient(f.Context, config)
			if err != nil {
				return err
			}

			if opts.resources == nil {
				opts.gvrs, err = client.ServerResources()
				if err != nil {
					return err
				}
			} else {
				opts.gvrs, err = resources.Parse(opts.resources)
				if err != nil {
					return err
				}
			}

			opts.gvrs = collections.Filter(opts.gvrs, func(gvr schema.GroupVersionResource) bool {
				if gvr.Group == "" && gvr.Version == "v1" && gvr.Resource == "secrets" {
					return false
				}

				return true
			})

			if opts.file == "" {
				opts.o, err = cmdutil.NewOutputWriter(f.IOStreams.Out)
				if err != nil {
					return err
				}
			} else {
				opts.o, err = cmdutil.NewOutputFileWriter(opts.file)
				if err != nil {
					return err
				}
			}

			defer func() {
				_ = opts.o.Close()
			}()

			items, err := client.List(opts.gvrs, metav1.ListOptions{
				LabelSelector: "cos.bf2.org/connector.id=" + opts.id,
			})

			if err != nil {
				return err
			}

			if len(items) != 0 {
				raw, err := yaml.Marshal(items)
				if err != nil {
					return err
				}

				_, err = opts.o.Write(raw)
				if err != nil {
					return err
				}
			}

			if opts.logs {
				for _, item := range items {
					if item.GetAPIVersion() == "v1" && item.GetKind() == "Pod" {
						err = client.Logs(item.GetNamespace(), item.GetName(), opts.o)
						if err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
	}

	cmdutil.AddID(cmd, &opts.id).Required()
	cmdutil.AddFile(cmd, &opts.file)

	cmd.Flags().BoolVar(&opts.logs, "logs", opts.logs, "Include logs")
	cmd.Flags().StringSliceVar(&opts.resources, "resource", nil, "resources to include apiVersion:plural")

	configFlags.AddFlags(cmd.PersistentFlags())

	return cmd
}
