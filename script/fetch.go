package script

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
)

type FetchFn func(url string, options map[string]interface{}) *goja.Promise

func GetFetchFn(vm *goja.Runtime) FetchFn {
	return func(url string, options map[string]interface{}) *goja.Promise {
		loop := eventloop.NewEventLoop()
		loop.Start()
		defer loop.Stop()
		fmt.Println("执行Promise,接收的参数", url, options)
		p, resolve, _ := vm.NewPromise()
		loop.RunOnLoop(func(vm *goja.Runtime) {
			go func() {
				time.Sleep(500 * time.Millisecond) // or perform any other blocking operation
				fmt.Println("等待Promise结束")
				resolve("golang返回的的Fetch 结果")
			}()
		})
		return p
	}
}
