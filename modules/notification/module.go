package notification

import (
	"context"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/modules/app"
	"rxdrag.com/entify/modules/register"
)

type SubscriptionModule struct {
	app *app.App
	ctx context.Context
}

func (m *SubscriptionModule) Init(ctx context.Context) {
	if contexts.Values(ctx).AppId == 0 {
		return
	}

	//没有安装
	if !app.Installed {
		return
	}

	app, err := app.Get(contexts.Values(ctx).AppId)
	if err != nil {
		log.Panic(err.Error())
	}
	m.app = app
	m.ctx = ctx
}

func (m *SubscriptionModule) QueryFields() []*graphql.Field {
	return []*graphql.Field{}
}

func (m *SubscriptionModule) MutationFields() []*graphql.Field {
	return []*graphql.Field{}
}

func (m *SubscriptionModule) Directives() []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m *SubscriptionModule) Types() []graphql.Type {
	return []graphql.Type{}
}
func (m *SubscriptionModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	register.RegisterModule(&SubscriptionModule{})
}
