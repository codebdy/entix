package resolve

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/script"
	"rxdrag.com/entify/utils"
)

func QueryThirdPartyResolveFn(third *graph.ThirdParty, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		if strings.Trim(third.Domain.QueryScript, " ") == "" {
			return nil, nil
		}
		vm := goja.New()
		script.Enable(vm)
		loop := eventloop.NewEventLoop()
		loop.Start()
		defer loop.Stop()

		wait := make(chan string, 1)
		vm.Set("callback", func(call goja.FunctionCall) goja.Value {
			fmt.Println("Go收到返回值")
			wait <- call.Argument(0).ToString().String()
			return nil
		})
		vm.RunString(third.Domain.QueryScript)

		return <-wait, nil
	}
}

func QueryOneThirdPartyResolveFn(third *graph.ThirdParty, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		return nil, nil
	}
}
