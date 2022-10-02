package resolve

import (
	"github.com/dop251/goja"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/utils"
)

func MethodResolveFn(method *graph.Method, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		vm := goja.New()
		v, err := vm.RunString(method.Method.Script)
		if err != nil {
			panic(err)
		}
		if num := v.Export().(int64); num != 4 {
			panic(num)
		}
		return v, nil
	}
}
