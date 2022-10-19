package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/authentication"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/logs"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/utils"
)

func LogoutResolveFn(model *model.Model) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		token := contexts.Values(p.Context).Token
		if token != "" {
			authentication.Logout(token)
		}
		logs.WriteBusinessLog(model, p, logs.LOGOUT, logs.SUCCESS, "")
		return true, nil
	}

}
