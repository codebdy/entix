package imexport

import (
	"log"
	"strconv"

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

	if p.Args[ARG_APP_ID] == nil {
		log.Panic("App id is nil")
	}

	appId, err := strconv.ParseUint(p.Args[ARG_APP_ID].(string), 10, 64)

	if err != nil {
		log.Panic(err)
	}

	return false, nil
}
