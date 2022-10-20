package resolve

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/logs"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/utils"
)

type InstallArg struct {
	Admin    string     `json:"admin"`
	Password string     `json:"password"`
	WithDemo bool       `json:"withDemo"`
	Meta     utils.JSON `json:"meta"`
}

const INPUT = "input"

func InstallResolve(p graphql.ResolveParams, model *model.Model) (interface{}, error) {
	defer utils.PrintErrorStack()
	input := InstallArg{}
	mapstructure.Decode(p.Args[INPUT], &input)

	repos := repository.New(model)
	repos.MakeSupperVerifier()

	instance, err := addAndPublishMeta(input.Meta, model)
	if err != nil {
		logs.WriteBusinessLog(model, p, logs.INSTALL, logs.FAILURE, err.Error())
		return nil, err
	}

	model = repos.LoadModel(contexts.Values(p.Context).AppId)

	if input.Admin != "" {
		instance = data.NewInstance(
			adminInstance(input.Admin, input.Password),
			model.Graph.GetEntityByName(consts.META_USER),
		)
		_, err = repos.SaveOne(instance)
		if err != nil {
			logs.WriteBusinessLog(model, p, logs.INSTALL, logs.FAILURE, err.Error())
			return nil, err
		}
		if input.WithDemo {
			instance = data.NewInstance(
				demoInstance(),
				model.Graph.GetEntityByName(consts.META_USER),
			)
			_, err = repos.SaveOne(instance)
			if err != nil {
				logs.WriteBusinessLog(model, p, logs.INSTALL, logs.FAILURE, err.Error())
				return nil, err
			}
		}
	}
	isExist := repository.IsEntityExists(consts.META_USER)
	logs.WriteBusinessLog(model, p, logs.INSTALL, logs.SUCCESS, "")
	return isExist, nil
}

func addAndPublishMeta(meta utils.JSON, model *model.Model) (*data.Instance, error) {
	appUuid := meta[consts.APPUUID].(string)
	metaContent := meta["content"]
	repos := repository.New(model)
	repos.MakeSupperVerifier()
	nextMeta := repos.QueryNextMeta(appUuid)
	if nextMeta != nil {
		panic("Please pushish meta first then install new function ")
	}

	publishedMeta := repos.QueryPublishedMeta(appUuid)

	var createdAt interface{} = time.Now()
	if publishedMeta != nil {
		metaContent := *(publishedMeta.(map[string]interface{})[consts.META_CONTENT].(*utils.JSON))
		createdAt = publishedMeta.(map[string]interface{})[consts.META_CREATEDAT]
		clses := metaContent[consts.META_CLASSES].([]interface{})
		if metaContent["classes"] != nil {
			for i := range clses {
				metaContent["classes"] = append(metaContent["classes"].([]utils.JSON), clses[i].(map[string]interface{}))
			}
		} else {
			metaContent["classes"] = clses
		}

		relas := metaContent[consts.META_RELATIONS].([]interface{})
		if metaContent["relations"] != nil {
			for i := range relas {
				metaContent["relations"] = append(metaContent["relations"].([]utils.JSON), relas[i].(map[string]interface{}))
			}
		} else {
			metaContent["relations"] = relas
		}

	}

	predefined := map[string]interface{}{
		consts.APPUUID:        appUuid,
		"content":             metaContent,
		consts.META_CREATEDAT: createdAt,
		consts.META_UPDATEDAT: time.Now(),
	}

	//创建实体
	instance := data.NewInstance(predefined, model.Graph.GetMetaEntity())
	_, err := repos.SaveOne(instance)

	if err != nil {
		return nil, err
	}
	err = doPublish(repos, consts.SYSTEM_APP_UUID)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func adminInstance(name string, password string) map[string]interface{} {
	return map[string]interface{}{
		consts.NAME:           "Admin",
		consts.LOGIN_NAME:     name,
		consts.PASSWORD:       utils.BcryptEncode(password),
		consts.IS_SUPPER:      true,
		consts.META_CREATEDAT: time.Now(),
		consts.META_UPDATEDAT: time.Now(),
	}
}

func demoInstance() map[string]interface{} {
	return map[string]interface{}{
		consts.NAME:           "Demo",
		consts.LOGIN_NAME:     "demo",
		consts.PASSWORD:       utils.BcryptEncode("demo"),
		consts.IS_DEMO:        true,
		consts.META_CREATEDAT: time.Now(),
		consts.META_UPDATEDAT: time.Now(),
	}
}
