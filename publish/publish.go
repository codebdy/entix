package publish

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app"
	"rxdrag.com/entify/logs"
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
	app.Publish()
}
