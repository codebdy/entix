package resolve

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/script"
	"rxdrag.com/entify/utils"
)

func argsString(method *graph.Method) string {
	names := []string{}
	for _, arg := range method.Method.Args {
		names = append(names, arg.Name)
	}
	return strings.Join(names, ", ")
}

func MethodResolveFn(method *graph.Method, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		vm := goja.New()
		script.Enable(vm)
		vm.Set("args", p.Args)
		script.Enable(vm)
		funcStr := fmt.Sprintf(
			`
			class Test2{
				test(){
					return "哈哈"
				}
			}
			function doMethod() {
				const {%s} = args;
			%s
			}`,
			argsString(method),
			method.Method.Script,
		)

		_, err := vm.RunString(funcStr)
		if err != nil {
			panic(err)
		}
		var doMethod func() string
		err = vm.ExportTo(vm.Get("doMethod"), &doMethod)
		if err != nil {
			panic(err)
		}

		result := doMethod()
		return result, nil
	}
}
