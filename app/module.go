package app

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
)

type AppModule struct {
}

func (m AppModule) QueryFields(ctx context.Context) []graphql.Field {
	return []graphql.Field{}
}
func (m AppModule) MutationFields(ctx context.Context) []graphql.Field {
	return []graphql.Field{}
}
func (m AppModule) SubscriptionFields(ctx context.Context) []graphql.Field {
	return []graphql.Field{}
}
func (m AppModule) Directives(ctx context.Context) []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m AppModule) Types(ctx context.Context) []graphql.Type {
	return []graphql.Type{}
}
func (m AppModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		LoadersMiddleware,
	}
}
