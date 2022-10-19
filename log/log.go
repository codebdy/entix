package log

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/repository"
)

func WriteModelLog(
	model *model.Model,
	cls *graph.Class,
	p graphql.ResolveParams,
	operate string,
	result string,
	message string,
) {
	repos := repository.New(model)
	repos.MakeSupperVerifier()
	contextsValues := contexts.Values(p.Context)
	logObject := map[string]interface{}{
		"ip":          contextsValues.IP,
		"appUuid":     contextsValues.AppUuid,
		"operateType": operate,
		"classUuid":   cls.Uuid(),
		"className":   cls.Name(),
		"gql":         p.Context.Value("gql"),
		"result":      result,
		"message":     message,
		"user": map[string]interface{}{
			"add": map[string]interface{}{
				"id": contextsValues.Me.Id,
			},
		},
	}
	instance := data.NewInstance(logObject, model.Graph.GetEntityByName("ModelLog"))
	repos.SaveOne(instance)
}
