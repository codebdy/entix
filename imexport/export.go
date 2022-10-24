package imexport

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/utils"
)

func exportQueryFields() []*graphql.Field {
	return []*graphql.Field{
		{
			Name: EXPORT_APP,
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				ARG_APP_ID: &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.ID,
					},
				},
			},
			Resolve: exportResolve,
		},
	}
}

func exportResolve(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()

	return "", nil
}
