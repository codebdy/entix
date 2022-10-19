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
		result, err := auth.Login(p.Args[consts.LOGIN_NAME].(string), p.Args[consts.PASSWORD].(string))
		if err != nil {
			log.WriteBusinessLog(model, p, log.LOGIN, log.FAILURE, err.Error())
		} else {
			log.WriteBusinessLog(model, p, log.LOGIN, log.SUCCESS, "")
		}
		return result, err
	}
}
