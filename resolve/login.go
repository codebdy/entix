package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/authentication"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/log"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/utils"
)

func LoginResolveFn(model *model.Model) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		auth := authentication.New()
		loginName := p.Args[consts.LOGIN_NAME].(string)
		result, err := auth.Login(loginName, p.Args[consts.PASSWORD].(string))
		if err != nil {
			log.WriteBusinessLog(model, p, log.LOGIN, log.FAILURE, ("Login name:"+loginName+", ")+err.Error())
		} else {
			log.WriteBusinessLog(model, p, log.LOGIN, log.SUCCESS, ("Login name:" + loginName))
		}
		return result, err
	}
}
