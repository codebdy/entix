package authentication

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/app/schema/parser"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
	"rxdrag.com/entify/register"
)

type AuthenticationModule struct {
}

func OutputFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		fields[attr.Name] = &graphql.Field{
			Type:        parser.PropertyType(attr.GetType()),
			Description: attr.Description,
		}
	}

	return fields
}
func (m *AuthenticationModule) Init(ctx context.Context) {
}
func (m *AuthenticationModule) QueryFields() []*graphql.Field {
	if orm.IsEntityExists(meta.USER_ENTITY_NAME) {
		systemApp := app.GetSystemApp()
		userType := systemApp.GetEntityByName(meta.USER_ENTITY_NAME)
		return []*graphql.Field{
			{
				Name: consts.ME,
				Type: graphql.NewObject(
					graphql.ObjectConfig{
						Name:   "Me",
						Fields: OutputFields(userType.AllAttributes()),
					},
				),
				Resolve: resolveMe,
			},
		}
	} else {
		return []*graphql.Field{}
	}
}
func (m *AuthenticationModule) MutationFields() []*graphql.Field {
	if app.Installed {
		return mutationFields()
	} else {
		return []*graphql.Field{}
	}
}
func (m *AuthenticationModule) SubscriptionFields() []*graphql.Field {
	return []*graphql.Field{}
}
func (m *AuthenticationModule) Directives() []*graphql.Directive {
	return []*graphql.Directive{}
}
func (m *AuthenticationModule) Types() []graphql.Type {
	return []graphql.Type{}
}
func (m *AuthenticationModule) Middlewares() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{}
}

func init() {
	register.RegisterModule(&AuthenticationModule{})
}
