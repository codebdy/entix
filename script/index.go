package script

import "github.com/dop251/goja"

func Enable(vm *goja.Runtime) {
	vm.Set("fetch", GetFetchFn(vm))
}
