package service

import (
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/orm"
)

// func QueryInterface(intf *graph.Interface, args graph.QueryArg) orm.QueryResponse {
// 	session, err := orm.Open()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return session.QueryInterface(intf, args)
// }

// func QueryOneInterface(intf *graph.Interface, args graph.QueryArg) interface{} {
// 	session, err := orm.Open()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return session.QueryOneInterface(intf, args)
// }

func (s *Service) QueryEntity(entity *graph.Entity, args graph.QueryArg) orm.QueryResponse {
	session, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return session.QueryEntity(entity, args)
}

func (s *Service) QueryOneEntity(entity *graph.Entity, args graph.QueryArg) interface{} {
	session, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return session.QueryOneEntity(entity, args)
}

func (s *Service) QueryById(entity *graph.Entity, id uint64) interface{} {
	return s.QueryOneEntity(entity, graph.QueryArg{
		consts.ARG_WHERE: graph.QueryArg{
			consts.ID: graph.QueryArg{
				consts.ARG_EQ: id,
			},
		},
	})
}

func (s *Service) BatchQueryAssociations(
	association *graph.Association,
	ids []uint64,
	args graph.QueryArg,
) []map[string]interface{} {
	session, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return session.BatchRealAssociations(association, ids, args)
}
