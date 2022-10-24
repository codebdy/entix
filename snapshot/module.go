package snapshot

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/register"
)

type SnapshotModule struct {
}

func (m *SnapshotModule) Init(ctx context.Context) {
}
func (m *SnapshotModule) QueryFields() []*graphql.Field {
	return []*graphql.Field{}
}
func (m *SnapshotModule) MutationFields() []*graphql.Field {
	if app.Installed {
		return []*graphql.Field{}
	} else {
		return []*graphql.Field{}
	}
}
func (m *SnapshotModule) SubscriptionFields() []*graphql.Field {
	return []*graphql.Field{}
}
func (m *SnapshotModule) Directives() []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m *SnapshotModule) Types() []graphql.Type {
	return []graphql.Type{}
}
func (m *SnapshotModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	register.RegisterModule(&SnapshotModule{})
}
