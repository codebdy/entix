package schema

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/opentracing/opentracing-go/log"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/resolve"
)

func publishResolve(p graphql.ResolveParams) (interface{}, error) {
	if p.Args[consts.APPID] == nil {
		err := errors.New("No App Id!")
		log.Error(err)
		return nil, err
	}

	app, err := app.GetAppByIdArg(p.Args[consts.APPID])
	if err != nil {
		return nil, err
	}
	app.Publish()
	return true, nil
}

func installResolve(p graphql.ResolveParams) (interface{}, error) {
	if !repository.IsEntityExists(consts.META_ENTITY_NAME) {
		repository.InstallMeta()
	}
	appUuid := contexts.Values(p.Context).AppId
	appSchema := Get(appUuid)
	appSchema.Make()

	result, err := resolve.InstallResolve(p, appSchema.Model())
	if err != nil {
		return result, err
	}
	appSchema.Make()
	return result, err
}
