package resolve

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app/script"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/utils"
)

func QueryThirdPartyResolveFn(third *graph.ThirdParty, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		if strings.Trim(third.Domain.QueryScript, " ") == "" {
			return nil, nil
		}
		defer utils.PrintErrorStack()
		vm := goja.New()
		script.Enable(vm)
		vm.Set("args", p.Args)
		script.Enable(vm)
		funcStr := fmt.Sprintf(
			`
			%s

			function doMethod() {
				
			%s
			}`,
			script.GetCommonCodes()+script.GetPackageCodes(model, third.Class.Domain.PackageUuid),
			//argsString(method),
			third.Domain.QueryScript,
		)

		_, err := vm.RunString(funcStr)
		if err != nil {
			panic(err)
		}
		var doMethod func() interface{}
		err = vm.ExportTo(vm.Get("doMethod"), &doMethod)
		if err != nil {
			panic(err)
		}

		result := doMethod()
		return result, nil
	}
}

func QueryOneThirdPartyResolveFn(third *graph.ThirdParty, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		return nil, nil
	}
}
