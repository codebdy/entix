package install

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/entry"
)

type BasicModule struct {
}

func (m BasicModule) Init(ctx context.Context) {
}
func (m BasicModule) QueryFields() []*graphql.Field {
	return installQueryFields()
}
func (m BasicModule) MutationFields() []*graphql.Field {
	if app.Installed {
		return []*graphql.Field{}
	} else {
		return installMutationFields()
	}
}
func (m BasicModule) SubscriptionFields() []*graphql.Field {
	return []*graphql.Field{}
}
func (m BasicModule) Directives() []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m BasicModule) Types() []graphql.Type {
	return []graphql.Type{}
}
func (m BasicModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	entry.AddModuler(BasicModule{})
}
