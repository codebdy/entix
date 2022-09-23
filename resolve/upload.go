package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/storage"
)

func UploadResolveResolveFn(appId uint64) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		file := p.Args[consts.ARG_FILE].(storage.File)
		fileInfo := file.Save()
		return GetFileUrl(fileInfo, p)
	}
}
