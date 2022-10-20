package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/logs"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/orm"
	"rxdrag.com/entify/service"
	"rxdrag.com/entify/utils"
)

func PostResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		objects := p.Args[consts.ARG_OBJECTS].([]interface{})
		instances := []*data.Instance{}
		for i := range objects {
			object := objects[i]
			ConvertObjectId(object.(map[string]interface{}))
			instance := data.NewInstance(object.(map[string]interface{}), entity)
			instances = append(instances, instance)
		}
		returing, err := service.Save(instances)

		if err != nil {
			return nil, err
		}
		logs.WriteModelLog(model, &entity.Class, p, logs.UPSERT, logs.SUCCESS, "")
		return returing, nil
	}
}

//未实现
func SetResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())

		set := p.Args[consts.ARG_SET].(map[string]interface{})
		objs := service.QueryEntity(entity, p.Args)[consts.NODES]
		convertedObjs := objs.([]orm.InsanceData)
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
		returing, err := service.Save(instances)

		if err != nil {
			return nil, err
		}

		logs.WriteModelLog(model, &entity.Class, p, logs.SET, logs.SUCCESS, "")

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
		//repos := repository.New(model)
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		instance := data.NewInstance(object, entity)
		result, err := service.SaveOne(instance)
		logs.WriteModelLog(model, &entity.Class, p, logs.UPSERT, logs.SUCCESS, "")
		return result, err
	}
}

func DeleteByIdResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		argId := p.Args[consts.ID]
		//repos := repository.New(model)
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		instance := data.NewInstance(map[string]interface{}{
			consts.ID: ConvertId(argId),
		}, entity)
		result, err := service.DeleteInstance(instance)
		logs.WriteModelLog(model, &entity.Class, p, logs.DELETE, logs.SUCCESS, "")
		return result, err
	}
}

func DeleteResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())

		objs := service.QueryEntity(entity, p.Args)[consts.NODES]

		if objs == nil || len(objs.([]orm.InsanceData)) == 0 {
			return map[string]interface{}{
				consts.RESPONSE_RETURNING:    []interface{}{},
				consts.RESPONSE_AFFECTEDROWS: 0,
			}, nil
		}

		convertedObjs := objs.([]orm.InsanceData)

		instances := []*data.Instance{}
		for i := range convertedObjs {
			instance := data.NewInstance(map[string]interface{}{
				consts.ID: ConvertId(convertedObjs[i][consts.ID]),
			}, entity)

			instances = append(instances, instance)
		}

		service.DeleteInstances(instances)
		logs.WriteModelLog(model, &entity.Class, p, logs.DELETE, logs.SUCCESS, "")
		return map[string]interface{}{
			consts.RESPONSE_RETURNING:    objs,
			consts.RESPONSE_AFFECTEDROWS: len(instances),
		}, nil
	}
}
