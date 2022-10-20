package logs

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/service"
)

func WriteModelLog(
	cls *graph.Class,
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

	instance := data.NewInstance(logObject, a.GetEntityByName("ModelLog"))
	service.SaveOne(instance)
}

func WriteBusinessLog(
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

	a.WriteUserBusinessLog(useId, p, operate, result, message)
}

func WriteUserBusinessLog(
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

	instance := data.NewInstance(logObject, a.GetEntityByName("BusinessLog"))
	service.SaveOne(instance)
}
