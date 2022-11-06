package service

import (
	"context"
	"log"

	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
)

type Service struct {
	isSystem bool
	ctx      context.Context
	roleIds  []uint64
}

func New(ctx context.Context) *Service {

	return &Service{
		isSystem: false,
		ctx:      ctx,
	}
}

func NewSystem() *Service {
	return &Service{
		isSystem: true,
	}
}

func QueryRoleIds(ctx context.Context, model *graph.Model) []uint64 {
	ids := []uint64{
		consts.GUEST_ROLE_ID,
	}

	me := contexts.Values(ctx).Me

	session, err := orm.Open()
	if err != nil {
		log.Panic(err.Error())
	}

	result := session.QueryEntity(model.GetEntityByName(meta.ROLE_ENTITY_NAME), map[string]interface{}{
		"users": map[string]interface{}{
			"id": map[string]interface{}{
				consts.ARG_EQ: me.Id,
			},
		},
	})

	for _, role := range result.Nodes {
		ids = append(ids, role[consts.ID].(uint64))
	}

	return ids
}
