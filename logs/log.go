package logs

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app/model"
	"rxdrag.com/entify/app/model/data"
	"rxdrag.com/entify/app/model/graph"
	"rxdrag.com/entify/common/contexts"
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

	instance := data.NewInstance(logObject, model.Graph.GetEntityByName("ModelLog"))
	repos.SaveOne(instance)
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

	WriteUserBusinessLog(useId, model, p, operate, result, message)
}

func WriteUserBusinessLog(
	useId string,
	model *model.Model,
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

	instance := data.NewInstance(logObject, model.Graph.GetEntityByName("BusinessLog"))
	repos.SaveOne(instance)
}
