package script

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
)

type GoFetchFn func(url string, options map[string]interface{}) *goja.Promise

func GetFetchFn(vm *goja.Runtime) GoFetchFn {
	return func(url string, options map[string]interface{}) *goja.Promise {
		loop := eventloop.NewEventLoop()
		loop.Start()
		defer loop.Stop()
		p, resolve, _ := vm.NewPromise()
		loop.RunOnLoop(func(vm *goja.Runtime) {
			go func() {
				method := http.MethodGet
				if options != nil && options["method"] != nil {
					method = options["method"].(string)
				}

				reqBody := []byte("")
				if options != nil && options["body"] != nil {
					reqBody = []byte(options["body"].(string))
				}
				client := &http.Client{}
				req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))

				if err != nil {
					fmt.Println(err)
					return
				}
				if options != nil && options["headers"] != nil {
					headers := options["headers"].(map[string]interface{})
					for key, header := range headers {
						if header != nil {
							req.Header.Add(key, header.(string))
						}
					}
				}

				res, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				resolve(string(body))
			}()
		})
		return p
	}
}
