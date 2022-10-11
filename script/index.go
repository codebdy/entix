package script

import (
	"github.com/dop251/goja"
	"rxdrag.com/entify/model"
)

func Enable(vm *goja.Runtime) {
	vm.Set("iFetch", FetchFn)
	vm.Set("writeToCache", WriteToCache)
	vm.Set("readFromCache", ReadFromCache)
}

func GetPackageMethods(model *model.Model) string {
	return ""
}
