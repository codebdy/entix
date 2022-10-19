package resolve

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/authentication"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/log"
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
		fmt.Println("哈哈", model, p)
		log.WriteBusinessLog(model, p, log.LOGOUT, log.SUCCESS, "")
		return true, nil
	}

}
