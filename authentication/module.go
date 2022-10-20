package authentication

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/entry"
)

type AuthenticationModule struct {
}

func (m AuthenticationModule) QueryFields(ctx context.Context) []*graphql.Field {
	return []*graphql.Field{}
}
func (m AuthenticationModule) MutationFields(ctx context.Context) []*graphql.Field {
	if !app.Installed {
		return []*graphql.Field{}
	} else {
		return mutationFields()
	}
}
func (m AuthenticationModule) SubscriptionFields(ctx context.Context) []*graphql.Field {
	return []*graphql.Field{}
}
func (m AuthenticationModule) Directives(ctx context.Context) []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m AuthenticationModule) Types(ctx context.Context) []graphql.Type {
	return []graphql.Type{}
}
func (m AuthenticationModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	entry.AddModuler(AuthenticationModule{})
}
