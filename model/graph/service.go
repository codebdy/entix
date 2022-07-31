package graph

import (
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
