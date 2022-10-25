package imexport

import (
	"log"
	"strconv"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/service"
	"rxdrag.com/entify/utils"
)

func (m *ImExportModule) QueryFields() []*graphql.Field {
	if !app.Installed {
		return []*graphql.Field{}
	}
	return []*graphql.Field{
		{
			Name: EXPORT_APP,
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				ARG_SNAPSHOT_ID: &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.ID,
					},
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				defer utils.PrintErrorStack()
				return m.exportResolve(p)
			},
		},
	}
}

func (m *ImExportModule) exportResolve(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()

	if p.Args[ARG_SNAPSHOT_ID] == nil {
		log.Panic("Snapshot id is nil")
	}

	snapshotId, err := strconv.ParseUint(p.Args[ARG_SNAPSHOT_ID].(string), 10, 64)

	if err != nil {
		log.Panic(err)
	}

	appSnapshot := service.QueryById(m.app.GetEntityByName("Snapshot"), snapshotId)

	if appSnapshot == nil {
		log.Panicf("App snapshot is nil on id:%d", snapshotId)
	}
	appJson := appSnapshot.(map[string]interface{})["content"]

	if appJson == nil {
		log.Panic("App json in snapshot is nil")
	}

	return "", nil
}
