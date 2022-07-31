package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/contexts"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/resolve"
)

func publishResolve(p graphql.ResolveParams) (interface{}, error) {
	appUuid := contexts.ParseAppUuid(p.Context)
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
	appUuid := contexts.ParseAppUuid(p.Context)
	appSchema := Get(appUuid)
	appSchema.Make()

	result, err := resolve.InstallResolve(p, appSchema.Model())
	if err != nil {
		return result, err
	}
	appSchema.Make()
	return result, err
}
