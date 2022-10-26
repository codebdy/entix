package imexport

import (
	"archive/zip"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/scalars"
	"rxdrag.com/entify/storage"
	"rxdrag.com/entify/utils"
)

func (m *ImExportModule) MutationFields() []*graphql.Field {
	if !app.Installed {
		return []*graphql.Field{}
	}
	return []*graphql.Field{
		{
			Name: IMPORT_APP,
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				ARG_APP_FILE: &graphql.ArgumentConfig{
					Type: scalars.UploadType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				defer utils.PrintErrorStack()
				return m.importResolve(p)
			},
		},
	}
}

func (m *ImExportModule) importResolve(p graphql.ResolveParams) (interface{}, error) {
	upload := p.Args[ARG_APP_FILE].(storage.File)
	fileInfo := upload.Save(TEMP_DATAS)
	// err := storage.Unzip(fileInfo.Path, fileInfo.Dir+fileInfo.NameBody)
	// if err != nil {
	// 	log.Panic(err.Error())
	// }
	r, err := zip.OpenReader(consts.STATIC_PATH + "/" + fileInfo.Path)
	if err != nil {
		log.Panic(err.Error())
	}

	var appJsonFile *zip.File
	for _, f := range r.File {
		fmt.Println("哈哈", f.Name)
		if f.Name == APP_JON {
			appJsonFile = f
		}
	}

	if appJsonFile == nil {
		log.Panic(fmt.Sprintf("Can not find %s in upload file", APP_JON))
	}

	return true, nil
}

func readAppJsonFile(f *zip.File) {
	rc, err := f.Open()
	if err != nil {
		log.Panic(err.Error())
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

}
