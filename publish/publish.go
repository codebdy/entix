package publish

import (
	"log"
	"time"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/logs"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/service"
	"rxdrag.com/entify/utils"
)

func PublishMetaResolveFn(app *app.App) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		defer utils.PrintErrorStack()
		doPublish(app)
		logs.WriteBusinessLog(app.Model, p, logs.PUBLISH_META, logs.SUCCESS, "")
		return true, nil
	}
}

func doPublish(app *app.App) {
	entity := app.GetEntityByName(meta.APP_ENTITY_NAME)
	appData := service.QueryById(
		entity,
		app.AppId,
	)
	if app == nil {
		log.Panic("App is nil")
	}

	appMap := appData.(map[string]interface{})
	appMap["publishedMeta"] = appMap["meta"]
	appMap["publishMetaAt"] = time.Now()
	instance := data.NewInstance(
		appMap,
		entity,
	)

	_, err := service.SaveOne(instance)

	if err != nil {
		log.Panic(err.Error())
	}

	app.ReLoad()
}
