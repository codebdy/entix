package authentication

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/common/errorx"
	"rxdrag.com/entify/utils"
)

func resolveMe(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()
	me := contexts.Values(p.Context).Me
	if me == nil {
		return nil, errorx.New(errorx.CODE_LOGIN_EXPIRED, "Login expired!")
	}
	return me, nil
}
