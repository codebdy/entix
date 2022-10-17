package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/utils"
)

func Me(p graphql.ResolveParams) (interface{}, error) {
	defer utils.PrintErrorStack()
	me := contexts.Values(p.Context).Me
	if me == nil {
		panic("Login expired!")
	}
	return me, nil
}
