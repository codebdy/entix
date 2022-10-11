package script

import "github.com/dop251/goja"

func Enable(vm *goja.Runtime) {
	vm.Set("iFetch", GoFetchFn)
	vm.Set("writeToCache", WriteToCache)
	vm.Set("readFromCache", ReadFromCache)
}

func GetPackageMethods() string {
	return ""
}
