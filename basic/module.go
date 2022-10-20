package basic

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/entry"
)

type BasicModule struct {
}

func (m BasicModule) QueryFields(ctx context.Context) []*graphql.Field {
	if app.Installed {
		return []*graphql.Field{}
	} else {
		return installQueryFields()
	}
}
func (m BasicModule) MutationFields(ctx context.Context) []*graphql.Field {
	if app.Installed {
		return []*graphql.Field{}
	} else {
		return installMutationFields()
	}
}
func (m BasicModule) SubscriptionFields(ctx context.Context) []*graphql.Field {
	return []*graphql.Field{}
}
func (m BasicModule) Directives(ctx context.Context) []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m BasicModule) Types(ctx context.Context) []graphql.Type {
	return []graphql.Type{}
}
func (m BasicModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	entry.AddModuler(BasicModule{})
}
