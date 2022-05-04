package show

import (
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/dumper"
	"github.com/redhat-developer/app-services-cli/pkg/shared/factory"
	"github.com/spf13/cobra"
)

type keyVal struct {
	key string
	val string
}

func NewConfigShowCommand(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "show",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Config.Load()
			if err != nil {
				return err
			}

			t := dumper.Table[keyVal]{}
			t.Field("Key", func(in *keyVal) string {
				return in.key
			})
			t.Field("Val", func(in *keyVal) string {
				return in.val
			})

			t.Dump(f.IOStreams.Out, []keyVal{
				{key: "API URL", val: c.APIUrl},
				{key: "Auth URL", val: c.AuthURL},
				{key: "Mas Auth URL", val: c.MasAuthURL},
			})

			return nil
		},
	}

	return cmd
}
