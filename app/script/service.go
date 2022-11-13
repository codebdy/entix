package script

import (
	"context"
	"log"

	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/orm"
	"rxdrag.com/entify/service"
)

type ScriptService struct {
	ctx     context.Context
	roleIds []uint64
	model   *graph.Model
	session *orm.Session
}

func NewService(ctx context.Context, model *graph.Model) *ScriptService {

	return &ScriptService{
		ctx:     ctx,
		model:   model,
		roleIds: service.QueryRoleIds(ctx, model),
	}
}

func (s *ScriptService) BeginTx() error {
	session, err := orm.Open()
	if err != nil {
		return err
	}
	s.session = session
	return session.BeginTx()
}

func (s *ScriptService) Commit() error {
	if s.session == nil {
		log.Panic("No session to commit")
	}
	return s.session.Commit()
}

func (s *ScriptService) ClearTx() {
	if s.session == nil {
		log.Panic("No session to ClearTx")
	}
	s.session.ClearTx()
	s.session = nil
}

func (s *ScriptService) Rollback() error {
	if s.session == nil {
		log.Panic("No session to Rollback")
	}

	err := s.session.Dbx.Rollback()
	if err != nil {
		return err
	}
	s.session = nil
	return nil
}
