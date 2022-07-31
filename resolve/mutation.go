package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/utils"
)

func PostOneResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		object := p.Args[consts.ARG_OBJECT].(map[string]interface{})
		ConvertObjectId(object)
		repos := repository.New(model)
		repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		instance := data.NewInstance(object, entity)
		return repos.SaveOne(instance)
	}
}

func DeleteByIdResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		argId := p.Args[consts.ID]
		repos := repository.New(model)
		repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		instance := data.NewInstance(map[string]interface{}{
			consts.ID: ConvertId(argId),
		}, entity)
		return repos.DeleteInstance(instance)
	}
}
