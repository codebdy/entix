package authentication

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/logs"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/utils"
)

func LoginResolveFn(model *model.Model) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		auth := New()
		loginName := p.Args[consts.LOGIN_NAME].(string)
		result, err := auth.Login(loginName, p.Args[consts.PASSWORD].(string))
		if err != nil {
			logs.WriteBusinessLog(model, p, logs.LOGIN, logs.FAILURE, ("Login name:"+loginName+", ")+err.Error())
		} else {
			logs.WriteBusinessLog(model, p, logs.LOGIN, logs.SUCCESS, ("Login name:" + loginName))
		}
		return result, err
	}
}

func LogoutResolveFn(model *model.Model) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		token := contexts.Values(p.Context).Token
		if token != "" {
			Logout(token)
		}
		logs.WriteBusinessLog(model, p, logs.LOGOUT, logs.SUCCESS, "")
		return true, nil
	}

}
