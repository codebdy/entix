package logs

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/service"
)

func WriteModelLog(
	model *model.Model,
	cls *graph.Class,
	p graphql.ResolveParams,
	operate string,
	result string,
	message string,
) {
	contextsValues := contexts.Values(p.Context)
	logObject := map[string]interface{}{
		"ip":          contextsValues.IP,
		"operateType": operate,
		"classUuid":   cls.Uuid(),
		"className":   cls.Name(),
		"gql":         p.Context.Value("gql"),
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

	instance := data.NewInstance(logObject, model.Graph.GetEntityByName("ModelLog"))
	service.SaveOne(instance)
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
	service.SaveOne(instance)
}
