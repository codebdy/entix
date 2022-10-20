package app

import (
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
)

func (a *App) PublishMeta(published, next *meta.MetaContent) error {
	diff := model.CreateDiff(published, next)
	orm.Migrage(diff)
}

func (a *App) Publish() error {
	return nil
}
