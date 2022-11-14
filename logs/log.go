package logs

import (
	"context"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/service"
)

func WriteModelLog(
	model *graph.Model,
	cls *graph.Class,
	ctx context.Context,
	operate string,
	result string,
	message string,
	gql interface{},
) {
	contextsValues := contexts.Values(ctx)
	logObject := map[string]interface{}{
		"ip":          contextsValues.IP,
		"operateType": operate,
		"classUuid":   cls.Uuid(),
		"className":   cls.Name(),
		"gql":         gql,
		"result":      result,
		"message":     message,
	}
	if contextsValues.Me != nil {
		logObject["user"] = map[string]interface{}{
			"add": map[string]interface{}{
				"id": contextsValues.Me.Id,
			},
		}
	}

	if contextsValues.AppId != 0 {
		logObject["app"] = map[string]interface{}{
			"add": map[string]interface{}{
				"id": contextsValues.AppId,
			},
		}
	}

	instance := data.NewInstance(logObject, model.GetEntityByName("ModelLog"))
	s := service.NewSystem()
	s.SaveOne(instance)
}

func WriteBusinessLog(
	model *model.Model,
	p graphql.ResolveParams,
	operate string,
	result string,
	message string,
) {
	contextsValues := contexts.Values(p.Context)

	useId := ""
	if contextsValues.Me != nil {
		useId = contextsValues.Me.Id
	}

	WriteUserBusinessLog(model, useId, p, operate, result, message)
}

func WriteUserBusinessLog(
	model *model.Model,
	useId string,
	p graphql.ResolveParams,
	operate string,
	result string,
	message string,
) {
	contextsValues := contexts.Values(p.Context)

	logObject := map[string]interface{}{
		"ip":          contextsValues.IP,
		"appUuid":     contextsValues.AppId,
		"operateType": operate,
		"result":      result,
		"message":     message,
	}
	if useId != "" {
		logObject["user"] = map[string]interface{}{
			"add": map[string]interface{}{
				"id": useId,
			},
		}
	}

	if contextsValues.AppId != 0 {
		logObject["app"] = map[string]interface{}{
			"add": map[string]interface{}{
				"id": contextsValues.AppId,
			},
		}
	}

	instance := data.NewInstance(logObject, model.Graph.GetEntityByName("BusinessLog"))
	s := service.NewSystem()
	s.SaveOne(instance)
}
