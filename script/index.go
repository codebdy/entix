package script

import (
	"github.com/dop251/goja"
	"rxdrag.com/entify/model"
)

func Enable(vm *goja.Runtime) {
	vm.Set("log", Log)
	vm.Set("iFetch", FetchFn)
	vm.Set("writeToCache", WriteToCache)
	vm.Set("readFromCache", ReadFromCache)
}

func GetPackageCodes(model *model.Model, packageUuid string) string {
	codeStr := ""
	for i := range model.Meta.Codes {
		code := model.Meta.Codes[i]
		if code.PackageUuid == packageUuid {
			codeStr = "\n" + code.Code
		}
	}
	return codeStr
}

func GetCommonCodes() string {
	return `
	const debug = {}
	debug.log = log
	`
}
