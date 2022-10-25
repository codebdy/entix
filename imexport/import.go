package imexport

import (
	"log"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/scalars"
	"rxdrag.com/entify/storage"
	"rxdrag.com/entify/utils"
)

func (m *ImExportModule) MutationFields() []*graphql.Field {
	if !app.Installed {
		return []*graphql.Field{}
	}
	return []*graphql.Field{
		{
			Name: IMPORT_APP,
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				ARG_APP_FILE: &graphql.ArgumentConfig{
					Type: scalars.UploadType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				defer utils.PrintErrorStack()
				return m.importResolve(p)
			},
		},
	}
}

func (m *ImExportModule) importResolve(p graphql.ResolveParams) (interface{}, error) {
	file := p.Args[ARG_APP_FILE].(storage.File)
	fileInfo := file.Save(TEMP_DATAS)
	err := storage.Unzip(fileInfo.Path, fileInfo.Dir+fileInfo.NameBody)
	if err != nil {
		log.Panic(err.Error())
	}
	return true, nil
}
