package contexts

import (
	"context"

	"rxdrag.com/entify/common"
	"rxdrag.com/entify/consts"
)

type ContextValues struct {
	Token   string
	Me      *common.User
	AppUuid string
	Host    string
}

func Values(ctx context.Context) ContextValues {
	values := ctx.Value(consts.CONTEXT_VALUES)
	if values == nil {
		panic("Not set CONTEXT_VALUES in context")
	}
	return values.(ContextValues)
}
