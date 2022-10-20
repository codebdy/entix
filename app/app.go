package app

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/meta"
)

var Installed = false

type App struct {
	AppId   uint64
	AppUuid string
	model   *model.Model
	schema  *graphql.Schema
}

var appCache = map[uint64]*App{}

func Get(appId uint64) *App {
	if appCache[appId] == nil {
		appCache[appId] = NewAppSchema(appId)
	}

	return appCache[appId]
}

func GetSystemApp() *App {
	for key := range appCache {
		if appCache[key].AppUuid == meta.SYSTEM_APP_UUID {
			return appCache[key]
		}
	}
	return nil
}
