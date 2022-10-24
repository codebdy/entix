package imexport

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/register"
)

type ImExportModule struct {
}

func (m *ImExportModule) Init(ctx context.Context) {
}

func (m *ImExportModule) SubscriptionFields() []*graphql.Field {
	return []*graphql.Field{}
}
func (m *ImExportModule) Directives() []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m *ImExportModule) Types() []graphql.Type {
	return []graphql.Type{}
}
func (m *ImExportModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	register.RegisterModule(&ImExportModule{})
}
