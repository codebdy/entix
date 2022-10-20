package resolve

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/service"
	"rxdrag.com/entify/utils"
)

func QueryOneInterfaceResolveFn(intf *graph.Interface, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeInterfaceAbilityVerifier(p, intf)
		instance := service.QueryOneInterface(intf, p.Args)
		return instance, nil
	}
}

func QueryInterfaceResolveFn(intf *graph.Interface, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeInterfaceAbilityVerifier(p, intf)
		result := service.QueryInterface(intf, p.Args)
		return result, nil
	}
}

func QueryOneEntityResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		instance := service.QueryOneEntity(entity, p.Args)
		return instance, nil
	}
}

func QueryEntityResolveFn(entity *graph.Entity, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeEntityAbilityVerifier(p, entity.Uuid())
		result := service.QueryEntity(entity, p.Args)
		return result, nil
	}
}

func QueryAssociationFn(asso *graph.Association, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		var (
			source      = p.Source.(map[string]interface{})
			v           = p.Context.Value
			loaders     = v(consts.LOADERS).(*Loaders)
			handleError = func(err error) error {
				return fmt.Errorf(err.Error())
			}
		)
		defer utils.PrintErrorStack()

		if loaders == nil {
			panic("Data loaders is nil")
		}
		loader := loaders.GetLoader(p, asso, p.Args, model)
		thunk := loader.Load(p.Context, NewKey(source[consts.ID].(uint64)))
		return func() (interface{}, error) {
			data, err := thunk()
			if err != nil {
				return nil, handleError(err)
			}

			var retValue interface{}
			if data == nil {
				if asso.IsArray() {
					retValue = []map[string]interface{}{}
				} else {
					retValue = nil
				}
			} else {
				retValue = data
			}
			return retValue, nil
		}, nil
	}
}
