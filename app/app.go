package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/opentracing/opentracing-go/log"
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

func GetAppByIdArg(idArg interface{}) (*App, error) {
	if idArg == nil {
		err := errors.New("Nil app id")
		log.Error(err)
		return nil, err
	}
	appIdStr := idArg.(string)
	appId, err := strconv.ParseUint(appIdStr, 10, 64)

	if err != nil {
		err := errors.New(fmt.Sprintf("App id error:%s", appIdStr))
		log.Error(err)
		return nil, err
	}
	return Get(appId)
}

func Get(appId uint64) (*App, error) {
	if appCache[appId] == nil {
		//appCache[appId] = NewAppSchema(appId)
	}

	return appCache[appId], nil
}

func GetSystemApp() *App {
	for key := range appCache {
		if appCache[key].AppUuid == meta.SYSTEM_APP_UUID {
			return appCache[key]
		}
	}
	return nil
}
