package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app/resolve"
	"rxdrag.com/entify/consts"
)

func (a *AppProcessor) publishField() *graphql.Field {
	return &graphql.Field{
		Type: graphql.Boolean,
		Args: graphql.FieldConfigArgument{
			consts.APPID: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.ID},
			},
		},
		Resolve: resolve.PublishMetaResolveFn(a.Model),
	}
}
