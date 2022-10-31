package app

import (
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
	"rxdrag.com/entify/service"
)

func PublishMeta(published, next *meta.MetaContent, appId uint64) {
	publishedModel := model.New(published, appId)
	nextModel := model.New(next, appId)
	diff := model.CreateDiff(publishedModel, nextModel)
	orm.Migrage(diff)
}

func (a *App) Publish() {
	entity := a.GetEntityByName(meta.APP_ENTITY_NAME)

	appData := service.QueryById(
		entity,
		a.AppId,
	)

	appMap := appData.(map[string]interface{})

	nextMeta := meta.MetaContent{}
	err := mapstructure.Decode(appMap["meta"], &nextMeta)
	if err != nil {
		log.Println(err.Error())
	}
	oldMeta := meta.MetaContent{}
	err = mapstructure.Decode(appMap["publishedMeta"], &oldMeta)
	if err != nil {
		log.Println(err.Error())
	}

	PublishMeta(MergeSystemModel(&oldMeta), MergeSystemModel(&nextMeta), a.AppId)

	appMap["publishedMeta"] = appMap["meta"]
	appMap["publishMetaAt"] = time.Now()
	instance := data.NewInstance(
		appMap,
		entity,
	)

	_, err = service.SaveOne(instance)

	if err != nil {
		log.Panic(err.Error())
	}

	a.ReLoad()
}

// func (a *App) MergeModel(content *meta.MetaContent) *meta.MetaContent {
// 	//合并系统Schema
// 	if a.AppId != meta.SYSTEM_APP_ID {
// 		systemAppData := service.QueryById(
// 			a.GetEntityByName(meta.APP_ENTITY_NAME),
// 			meta.SYSTEM_APP_ID,
// 		)

// 		systemContent := systemAppData.(map[string]interface{})["publishedMeta"].(meta.MetaContent)
// 		//systemMetaContent := r.LoadAndDecodeMeta(consts.SYSTEM_APP_UUID)
// 		for i := range systemContent.Classes {
// 			content.Classes = append(content.Classes, systemContent.Classes[i])
// 		}

// 		for i := range systemContent.Relations {
// 			content.Relations = append(content.Relations, systemContent.Relations[i])
// 		}
// 	}

// 	return content
// }
