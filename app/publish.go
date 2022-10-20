package app

import (
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
)

func PublishMeta(published, next *meta.MetaContent, appId uint64) {
	publishedModel := model.New(published, appId)
	nextModel := model.New(next, appId)
	diff := model.CreateDiff(publishedModel, nextModel)
	orm.Migrage(diff)
}

func (a *App) Publish() error {
	//需要MergeModel
	return nil
}
