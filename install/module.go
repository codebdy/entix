package install

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/entry"
)

type InstalModule struct {
}

func (m InstalModule) QueryFields(ctx context.Context) []*graphql.Field {
	if app.Installed {
		return []*graphql.Field{}
	} else {
		return queryFields()
	}
}
func (m InstalModule) MutationFields(ctx context.Context) []*graphql.Field {
	if app.Installed {
		return []*graphql.Field{}
	} else {
		return mutationFields()
	}
}
func (m InstalModule) SubscriptionFields(ctx context.Context) []*graphql.Field {
	return []*graphql.Field{}
}
func (m InstalModule) Directives(ctx context.Context) []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m InstalModule) Types(ctx context.Context) []graphql.Type {
	return []graphql.Type{}
}
func (m InstalModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	entry.AddModuler(InstalModule{})
}
