package service

import "context"

type Service struct {
	isSystem bool
	ctx      context.Context
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
