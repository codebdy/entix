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

func PostResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		repos := repository.New(model)
		repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		objects := p.Args[consts.ARG_OBJECTS].([]map[string]interface{})
		instances := []*data.Instance{}
		for i := range objects {
			object := objects[i]
			ConvertObjectId(object)
			instance := data.NewInstance(object, entity)
			instances = append(instances, instance)
		}
		returing, err := repos.Save(instances)

		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			consts.RESPONSE_RETURNING:    returing,
			consts.RESPONSE_AFFECTEDROWS: len(instances),
		}, nil
	}
}

//未实现
func SetResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		repos := repository.New(model)
		repos.MakeEntityAbilityVerifier(p, entity.Uuid())

		set := p.Args[consts.ARG_SET].(map[string]interface{})
		objs := repos.QueryEntity(entity, p.Args)[consts.NODES]
		convertedObjs := objs.([]repository.InsanceData)
		instances := []*data.Instance{}

		for i := range convertedObjs {
			obj := convertedObjs[i]
			object := map[string]interface{}{}

			object[consts.ID] = obj[consts.ID]

			for key := range set {
				object[key] = set[key]
				instance := data.NewInstance(object, entity)
				instances = append(instances, instance)
			}
		}
		returing, err := repos.Save(instances)

		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			consts.RESPONSE_RETURNING:    returing,
			consts.RESPONSE_AFFECTEDROWS: len(instances),
		}, nil
	}
}

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

func DeleteResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		repos := repository.New(model)
		repos.MakeEntityAbilityVerifier(p, entity.Uuid())

		objs := repos.QueryEntity(entity, p.Args)[consts.NODES]

		if objs == nil || len(objs.([]repository.InsanceData)) == 0 {
			return map[string]interface{}{
				consts.RESPONSE_RETURNING:    []interface{}{},
				consts.RESPONSE_AFFECTEDROWS: 0,
			}, nil
		}

		convertedObjs := objs.([]repository.InsanceData)

		instances := []*data.Instance{}
		for i := range convertedObjs {
			instance := data.NewInstance(map[string]interface{}{
				consts.ID: ConvertId(convertedObjs[i][consts.ID]),
			}, entity)

			instances = append(instances, instance)
		}

		repos.DeleteInstances(instances)

		return map[string]interface{}{
			consts.RESPONSE_RETURNING:    objs,
			consts.RESPONSE_AFFECTEDROWS: len(instances),
		}, nil
	}
}
