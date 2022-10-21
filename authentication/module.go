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

func (m AuthenticationModule) Init(ctx context.Context) {
}
func (m AuthenticationModule) QueryFields() []*graphql.Field {
	return []*graphql.Field{}
}
func (m AuthenticationModule) MutationFields() []*graphql.Field {
	if app.Installed {
		return mutationFields()
	} else {
		return []*graphql.Field{}
	}
}
func (m AuthenticationModule) SubscriptionFields() []*graphql.Field {
	return []*graphql.Field{}
}
func (m AuthenticationModule) Directives() []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m AuthenticationModule) Types() []graphql.Type {
	return []graphql.Type{}
}
func (m AuthenticationModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	entry.AddModuler(AuthenticationModule{})
}
