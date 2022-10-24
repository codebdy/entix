package snapshot

import (
	"log"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/utils"
)

const (
	APP_ID      = "appId"
	INSTANCE_ID = "instaneId"
	VERSION     = "version"
	DESCRIPTION = "description"
)

func (m *SnapshotModule) MutationFields() []*graphql.Field {
	if !app.Installed {
		return []*graphql.Field{}
	}
	return []*graphql.Field{
		{
			Name: "makeVersion",
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				APP_ID: &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.ID,
					},
				},
				INSTANCE_ID: &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.ID,
					},
				},
				VERSION: &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
				DESCRIPTION: &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				defer utils.PrintErrorStack()
				return m.makeVersion(p)
			},
		},
	}
}

func (m *SnapshotModule) makeVersion(p graphql.ResolveParams) (interface{}, error) {
	appId := utils.Uint64Value(p.Args[APP_ID])
	if appId == 0 {
		log.Panic("App id is nil")
	}
	instanceId := utils.Uint64Value(p.Args[INSTANCE_ID])

	if instanceId == 0 {
		log.Panic("Instance id is nil")
	}
	//gqlSchema := register.GetSchema(p.Context)
	return false, nil
}
