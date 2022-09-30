package main

import (
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/util/cmdutil"
	"os"

	"github.com/bf2fc6cc711aee1a0c2a/cos-tools/rhoc/pkg/cmd/root"
	"github.com/redhat-developer/app-services-cli/pkg/core/config"
)

func main() {
	f, err := cmdutil.NewFactory()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	cfgFile, err := f.Config.Load()
	if cfgFile == nil {
		if !os.IsNotExist(err) {
			_, _ = fmt.Fprintln(f.IOStreams.ErrOut, err)
			os.Exit(1)
		}

		cfgFile = &config.Config{}
		if err := f.Config.Save(cfgFile); err != nil {
			_, _ = fmt.Fprintln(f.IOStreams.ErrOut, err)
			os.Exit(1)
		}
	}

	err = root.NewRootCommand(f).Execute()
	if err == nil {
		return
	}
}
