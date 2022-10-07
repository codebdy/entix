package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/utils"
)

func DeployProcessResolveFn(model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		argId := p.Args[consts.ID]
		repos := repository.New(model)
		//@@@后面需要修改权限
		repos.MakeEntityAbilityVerifier(p, model.Graph.GetEntityByName("Process").Uuid())
		process := repos.QueryOneEntity(model.Graph.GetEntityByName("Process"), graph.QueryArg{
			consts.ARG_WHERE: graph.QueryArg{
				consts.ID: graph.QueryArg{
					consts.ARG_EQ: argId,
				},
			},
		})

		if process == nil {
			panic("can not find process by id")
		}
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		// instance := data.NewInstance(map[string]interface{}{
		// 	consts.ID: ConvertId(argId),
		// }, entity)
		// return repos.DeleteInstance(instance)
		return argId, nil
	}
}
