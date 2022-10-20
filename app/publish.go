package app

import (
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
)

func (a *App) PublishMeta(published, next *meta.MetaContent) {
	publishedModel := model.New(published)
	nextModel := model.New(next)
	diff := model.CreateDiff(publishedModel, nextModel)
	orm.Migrage(diff)
}

func (a *App) Publish() error {
	//需要MergeModel
	return nil
}
