package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/storage"
	"rxdrag.com/entify/utils"
)

func UploadResolveResolveFn(appId uint64) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		file := p.Args[consts.ARG_FILE].(storage.File)
		fileInfo := file.Save(appId, consts.UPLOAD_PATH)
		return GetFileUrl(fileInfo, p)
	}
}

func UploadPluginResolveResolveFn(appId uint64) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		file := p.Args[consts.ARG_FILE].(storage.File)
		fileInfo := file.Save(appId, consts.PLUGINS_PATH)
		err := storage.UnZip(fileInfo.Path, fileInfo.Dir)
		if err != nil {
			panic(err)
		}
		return GetFileUrl(fileInfo, p)
	}
}
