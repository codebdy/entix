package snapshot

import (
	"fmt"
	"log"
	"strings"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/model/graph"
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

	entityInnerId := utils.DecodeEntityInnerId(instanceId)

	entity := m.app.GetEntityByInnerId(entityInnerId)

	if entity == nil {
		log.Panic(fmt.Sprintf("Can not find entity by inner id:%d", entityInnerId))
	}

	queryGql := fmt.Sprintf(`
	query ($id:ID!){
		one%s(where:{
			id:{
				_eq:$id
			}
		})
		%s
	}
	`,
		entity.Name(),
		m.makeFieldsGql(entity),
	)
	//gqlSchema := register.GetSchema(p.Context)

	fmt.Println(queryGql)
	return false, nil
}

func (m *SnapshotModule) makeFieldsGql(entity *graph.Entity) string {
	fieldStrings := strings.Join(entity.AllAttributeNames(), "\n")
	for _, assoc := range entity.Associations() {
		if assoc.IsCombination() {
			subFields := m.makeFieldsGql(assoc.TypeEntity())
			fieldStrings = fieldStrings + subFields
		}
	}
	return fmt.Sprintf("{\n%s\n}", fieldStrings)
}
