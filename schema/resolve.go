package schema

import (
	"errors"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/resolve"
)

func publishResolve(p graphql.ResolveParams) (interface{}, error) {
	if p.Args[consts.APPUUID] == nil {
		return nil, errors.New("No appuuid!")
	}
	appUuid := p.Args[consts.APPUUID].(string)
	appSchema := Get(appUuid)
	result, err := resolve.PublishMetaResolve(p, appSchema.Model())
	if err != nil {
		return result, err
	}
	appSchema.Make()
	return result, nil
}

func installResolve(p graphql.ResolveParams) (interface{}, error) {
	if !repository.IsEntityExists(consts.META_ENTITY_NAME) {
		repository.InstallMeta()
	}
	appUuid := contexts.Values(p.Context).AppUuid
	appSchema := Get(appUuid)
	appSchema.Make()

	result, err := resolve.InstallResolve(p, appSchema.Model())
	if err != nil {
		return result, err
	}
	appSchema.Make()
	return result, err
}
