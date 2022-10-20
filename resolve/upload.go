package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/storage"
	"rxdrag.com/entify/utils"
)

func UploadResolveResolve(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()
	file := p.Args[consts.ARG_FILE].(storage.File)
	fileInfo := file.Save(consts.UPLOAD_PATH)
	return GetFileUrl(fileInfo, p)
}

func UploadPluginResolveResolve(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()
	file := p.Args[consts.ARG_FILE].(storage.File)
	fileInfo := file.Save(consts.PLUGINS_PATH)
	err := storage.Unzip(fileInfo.Path, fileInfo.Dir+fileInfo.NameBody)
	if err != nil {
		panic(err)
	}
	return GetFileUrl(fileInfo, p)
}
