package script

import "github.com/dop251/goja"

func Enable(vm *goja.Runtime) {
	vm.Set("goFetch", GetFetchFn(vm))
}
