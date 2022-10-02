package graph

import (
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/domain"
)

type Service struct {
	Entity
}

func NewService(c *domain.Class) *Service {
	return &Service{
		Entity: Entity{
			Class: *NewClass(c),
		},
	}
}

func (s *Service) QueryMethods() []*Method {
	methods := []*Method{}
	for _, method := range s.AllMethods() {
		if method.Method.MethodMeta.OperateType == consts.QUERY {
			methods = append(methods, method)
		}
	}
	return methods
}

func (s *Service) MetationMethods() []*Method {
	methods := []*Method{}
	for _, method := range s.AllMethods() {
		if method.Method.MethodMeta.OperateType == consts.MUTATION {
			methods = append(methods, method)
		}
	}
	return methods
}
