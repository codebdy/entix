package resolve

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app/script"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/utils"
)

func MethodResolveFn(code string, methodArgs []meta.ArgMeta, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		scriptService := script.NewService(p.Context, model.Graph)
		vm := goja.New()
		script.Enable(vm)
		vm.Set("$args", p.Args)
		vm.Set("$BeginTx", scriptService.BeginTx)
		script.Enable(vm)
		funcStr := fmt.Sprintf(
			`
			%s
			%s
			`,
			script.GetCodes(model),
			code,
		)

		result, err := vm.RunString(funcStr)
		if err != nil {
			panic(err)
		}
		return result, nil
	}
}
