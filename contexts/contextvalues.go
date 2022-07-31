package contexts

import (
	"context"

	"rxdrag.com/entify/common"
	"rxdrag.com/entify/consts"
)

type ContextValues struct {
	Token          string
	Me             *common.User
	QueryUserCache map[string][]string
}

func ParseContextValues(ctx context.Context) ContextValues {
	values := ctx.Value(consts.CONTEXT_VALUES)
	if values == nil {
		panic("Not set CONTEXT_VALUES in context")
	}

	return values.(ContextValues)
}

func ParseAppUuid(ctx context.Context) string {
	value := ctx.Value(consts.APPUUID)
	if value == nil {
		panic("Not set appUuid in context")
	}

	return value.(string)
}
