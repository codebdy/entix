package app

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"rxdrag.com/entify/app/schema"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/orm"
	"rxdrag.com/entify/service"
)

//节省开支，运行时使用，初始化时请使用orm.IsEntityExists
var Installed = false

type App struct {
	AppId  uint64
	Model  *model.Model
	Schema schema.AppGraphqlSchema
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
		log.Panic(err.Error())
	}
	appIdStr := idArg.(string)
	appId, err := strconv.ParseUint(appIdStr, 10, 64)

	if err != nil {
		err := errors.New(fmt.Sprintf("App id error:%s", appIdStr))
		log.Panic(err.Error())
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

func (a *App) ReLoad() {
	newApp := NewApp(a.AppId)
	a.Model = newApp.Model
	a.Schema = newApp.Schema
}

func NewApp(appId uint64) *App {
	systemApp := GetSystemApp()

	appMeta := service.QueryById(
		systemApp.GetEntityByName(meta.APP_ENTITY_NAME),
		appId,
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

		model := model.New(content, appId)
		return &App{
			AppId:  appId,
			Model:  model,
			Schema: schema.New(model),
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
