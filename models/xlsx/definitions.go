package xlsx

import (
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/sarulabs/di/v2"
)

var Definitions = []di.Def{
	{
		Build: func(ctn di.Container) (interface{}, error) {
			repositoryResolver := ctn.Get("focus.models.repositories.resolver").(actions.RepositoryResolver)
			return NewExporter(repositoryResolver, 200), nil
		},
		Name: "focus.models.exporter",
	},
}
