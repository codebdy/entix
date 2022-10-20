package app

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/orm"
)

func (a *App) WriteModelLog(
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
	sesson, err := orm.Open()
	if err != nil {
		panic(err)
	}
	sesson.SaveOne(instance)
}

func (a *App) WriteBusinessLog(
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

	a.WriteUserBusinessLog(useId, model, p, operate, result, message)
}

func (a *App) WriteUserBusinessLog(
	useId string,
	model *model.Model,
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

	instance := data.NewInstance(logObject, model.Graph.GetEntityByName("BusinessLog"))
	sesson, err := orm.Open()
	if err != nil {
		panic(err)
	}
	sesson.SaveOne(instance)
}
