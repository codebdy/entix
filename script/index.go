package script

import "github.com/dop251/goja"

func Enable(vm *goja.Runtime) {
	vm.Set("goFetch", GoFetchFn)
	vm.Set("writeToCache", WriteToCache)
	vm.Set("readFromCache", ReadFromCache)
}
