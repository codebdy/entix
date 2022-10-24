package imexport

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/scalars"
	"rxdrag.com/entify/utils"
)

func importMutationFields() []*graphql.Field {
	return []*graphql.Field{
		{
			Name: EXPORT_APP,
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				ARG_APP_FILE: &graphql.ArgumentConfig{
					Type: scalars.UploadType,
				},
			},
			Resolve: exportResolve,
		},
	}
}

func importResolve(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()

	return false, nil
}
