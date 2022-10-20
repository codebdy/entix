package resolve

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/logs"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/utils"
)

func PublishMetaResolveFn(model *model.Model) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		//repos := repository.New(model)
		//repos.MakeSupperVerifier()
		appId := p.Args[consts.APPID]
		if appId == nil {
			appId = "1"
		}
		err := doPublish(model, appId.(string))
		if err != nil {
			logs.WriteBusinessLog(model, p, logs.PUBLISH_META, logs.FAILURE, err.Error())
		} else {
			logs.WriteBusinessLog(model, p, logs.PUBLISH_META, logs.SUCCESS, "")
		}
		return "success", err
	}
}

func doPublish(model *model.Model, appUuid string) error {
	// publishedMeta := r.QueryPublishedMeta(appUuid)
	// nextMeta := r.QueryNextMeta(appUuid)
	// appId := r.QueryAppId(appUuid)
	// fmt.Println("Start to publish")
	// // fmt.Println("Published Meta ID:", publishedMeta.(utils.Object)["id"])
	// // fmt.Println("Next Meta ID:", nextMeta.(utils.Object)["id"])

	// if nextMeta == nil {
	// 	panic("Can not find unpublished meta")
	// }
	// publishedModel := model.New(r.Model.AppUuid, r.MergeModel(appUuid, repository.DecodeContent(publishedMeta, appId)))
	// nextModel := model.New(r.Model.AppUuid, r.MergeModel(appUuid, repository.DecodeContent(nextMeta, appId)))
	// nextModel.Graph.Validate()
	// diff := model.CreateDiff(publishedModel, nextModel)
	// r.ExcuteDiff(diff)
	// fmt.Println("ExcuteDiff success")
	// metaObj := nextMeta.(utils.Object)
	// metaObj[consts.META_STATUS] = meta.META_STATUS_PUBLISHED
	// metaObj[consts.META_PUBLISHEDAT] = time.Now()
	// _, err := r.SaveOne(data.NewInstance(metaObj, r.Model.Graph.GetMetaEntity()))
	// if err != nil {
	// 	return err
	// }

	return nil
}
