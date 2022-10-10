package resolve

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
		wait := make(chan interface{}, 1)
		timeout := false
		error := false
		timer := loop.SetTimeout(func(*goja.Runtime) {
			timeout = true
			wait <- nil
		}, 2*time.Second)

		vm.Set("Return", func(call goja.FunctionCall) goja.Value {
			fmt.Println("Go收到返回值")
			error = true
			wait <- call.Argument(0).ToString().String()
			loop.ClearTimeout(timer)
			return nil
		})

		vm.Set("Error", func(call goja.FunctionCall) goja.Value {
			fmt.Println("Go收到Error")
			wait <- call.Argument(0).ToString().String()
			loop.ClearTimeout(timer)
			return nil
		})

		vm.RunString(third.Domain.QueryScript)

		result := <-wait
		if timeout {
			return nil, errors.New("Time out")
		} else if error {
			return nil, errors.New(result.(string))
		}
		return result, nil
	}
}

func QueryOneThirdPartyResolveFn(third *graph.ThirdParty, model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		return nil, nil
	}
}
