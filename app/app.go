package app

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
	"rxdrag.com/entify/service"
)

var Installed = false

type App struct {
	AppId  uint64
	Model  *model.Model
	Schema *graphql.Schema
}

var appCache = map[uint64]*App{}

func init() {
	//先加载系统APP
	if orm.IsEntityExists(meta.APP_ENTITY_NAME) {
		Get(meta.SYSTEM_APP_ID)
	}
}

func GetAppByIdArg(idArg interface{}) (*App, error) {
	if idArg == nil {
		err := errors.New("Nil app id")
		log.Panic(err)
	}
	appIdStr := idArg.(string)
	appId, err := strconv.ParseUint(appIdStr, 10, 64)

	if err != nil {
		err := errors.New(fmt.Sprintf("App id error:%s", appIdStr))
		log.Panic(err)
	}
	return Get(appId)
}

func Get(appId uint64) (*App, error) {
	if appCache[appId] == nil {
		app := NewApp(appId)
		if app != nil {
			appCache[appId] = app
		} else {
			log.Panic(errors.New("Cant load app"))
		}
	}

	return appCache[appId], nil
}

func GetSystemApp() *App {
	if appCache[meta.SYSTEM_APP_ID] != nil {
		return appCache[meta.SYSTEM_APP_ID]
	}

	metaConent := meta.SystemAppData["meta"].(meta.MetaContent)

	return &App{
		AppId: meta.SystemAppData["id"].(uint64),
		Model: model.New(&metaConent, meta.SYSTEM_APP_ID),
	}
}

func (a *App) GetEntityByName(name string) *graph.Entity {
	return a.Model.Graph.GetEntityByName(name)
}

func NewApp(appId uint64) *App {
	systemApp := GetSystemApp()

	appMeta := service.QueryOneEntity(
		systemApp.GetEntityByName(meta.APP_ENTITY_NAME),
		graph.QueryArg{
			consts.ARG_WHERE: graph.QueryArg{
				consts.ID: graph.QueryArg{
					consts.ARG_EQ: appId,
				},
			},
		},
	)

	if appMeta != nil {
		publishedMeta := appMeta.(map[string]interface{})["publishedMeta"]
		var content *meta.MetaContent
		if publishedMeta != nil {
			content = DecodeContent(publishedMeta)
		}
		if appId != meta.SYSTEM_APP_ID {
			content = MergeSystemModel(content)
		}

		return &App{
			AppId: appId,
			Model: model.New(content, appId),
		}
	}

	return nil
}

func DecodeContent(obj interface{}) *meta.MetaContent {
	content := meta.MetaContent{}
	if obj != nil {
		err := mapstructure.Decode(obj, &content)
		if err != nil {
			panic("Decode content failure:" + err.Error())
		}
	}
	return &content
}

func MergeSystemModel(content *meta.MetaContent) *meta.MetaContent {
	if content == nil {
		content = &meta.MetaContent{}
	}
	//合并系统Schema
	systemModel := GetSystemApp().Model
	for i := range systemModel.Meta.Classes {
		content.Classes = append(content.Classes, *systemModel.Meta.Classes[i])
	}

	for i := range systemModel.Meta.Relations {
		content.Relations = append(content.Relations, *systemModel.Meta.Relations[i])
	}
	return content
}
