package cmdcontext

import (
	"github.com/redhat-developer/app-services-cli/pkg/core/servicecontext"
	"os"
)

const defaultProfileName = "default"

func Init(sc servicecontext.IContext) error {
	ctxFile, err := sc.Load()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		ctxFile = &servicecontext.Context{}
	}

	if _, ok := ctxFile.Contexts[defaultProfileName]; !ok {
		if ctxFile.Contexts == nil {
			ctxFile.Contexts = make(map[string]servicecontext.ServiceConfig)
		}

		ctxFile.Contexts[defaultProfileName] = servicecontext.ServiceConfig{}
		ctxFile.CurrentContext = defaultProfileName
	}

	if err := sc.Save(ctxFile); err != nil {
		return err
	}

	return nil
}

func Create(sc servicecontext.Context, name string) {
	if sc.Contexts == nil {
		sc.Contexts = make(map[string]servicecontext.ServiceConfig)
	}
	if _, ok := sc.Contexts[name]; !ok {
		sc.Contexts[name] = servicecontext.ServiceConfig{}
	}
}

func List(sc servicecontext.Context) []string {
	items := make([]string, len(sc.Contexts))

	for k, _ := range sc.Contexts {
		items = append(items, k)
	}

	return items
}
