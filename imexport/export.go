package imexport

import (
	"fmt"
	"log"
	"strconv"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/utils"
)

var queryGql = `
query($id:ID!){
	oneApp(where:{
		id:{
			_eq:$id
		}
	}){
		id
		uuid
		title
		description
		pages
		menus
		imageUrl
		pageFrames
		publishedMeta
		plugins{
			id
			url
			title
			pluginId
			type
			description
			version
		}
	}
}
`

func (m *ImExportModule) QueryFields() []*graphql.Field {
	if !app.Installed {
		return []*graphql.Field{}
	}
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

	if p.Args[ARG_APP_ID] == nil {
		log.Panic("App id is nil")
	}

	appId, err := strconv.ParseUint(p.Args[ARG_APP_ID].(string), 10, 64)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("哈哈", appId)
	return "", nil
}
