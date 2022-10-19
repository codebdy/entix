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

	model = repos.LoadModel(contexts.Values(p.Context).AppUuid)

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

// var mediaClasses = []map[string]interface{}{
// 	{
// 		consts.NAME:    consts.MEDIA_ENTITY_NAME,
// 		consts.UUID:    consts.MEDIA_UUID,
// 		consts.INNERID: consts.MEDIA_INNER_ID,
// 		consts.ROOT:    true,
// 		consts.SYSTEM:  true,
// 		"attributes": []map[string]interface{}{
// 			{
// 				consts.NAME:   "id",
// 				consts.TYPE:   "ID",
// 				consts.UUID:   "RX_MEDIA_ID_UUID",
// 				"primary":     true,
// 				"typeLabel":   "ID",
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "name",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_MEDIA_NAME_UUID",
// 				"typeLabel":   "String",
// 				"nullable":    true,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "mimetype",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_MEDIA_MIMETYPE_UUID",
// 				"typeLabel":   "String",
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "fileName",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_MEDIA_FILENAME_UUID",
// 				"typeLabel":   "String",
// 				"length":      128,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "path",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_MEDIA_PATH_UUID",
// 				"typeLabel":   "String",
// 				"length":      256,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "size",
// 				consts.TYPE:   "Int",
// 				consts.UUID:   "RX_MEDIA_SIZE_UUID",
// 				"typeLabel":   "Int",
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "mediaType",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_MEDIA_MEDIATYPE_UUID",
// 				"typeLabel":   "String",
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:       consts.META_CREATEDAT,
// 				consts.TYPE:       "Date",
// 				consts.UUID:       "RX_MEDIA_CREATEDAT_UUID",
// 				"typeLabel":       "Date",
// 				consts.CREATEDATE: true,
// 				consts.SYSTEM:     true,
// 			},
// 			{
// 				consts.NAME:       consts.META_UPDATEDAT,
// 				consts.TYPE:       "Date",
// 				consts.UUID:       "RX_MEDIA_UPDATEDAT_UUID",
// 				"typeLabel":       "Date",
// 				consts.UPDATEDATE: true,
// 				consts.SYSTEM:     true,
// 			},
// 		},
// 		"stereoType": "Entity",
// 	},
// }

// var authClasses = []map[string]interface{}{
// 	{
// 		consts.NAME:    consts.META_USER,
// 		consts.UUID:    consts.USER_UUID,
// 		consts.INNERID: consts.USER_INNER_ID,
// 		consts.ROOT:    true,
// 		consts.SYSTEM:  true,
// 		"attributes": []map[string]interface{}{
// 			{
// 				consts.NAME:   "id",
// 				consts.TYPE:   "ID",
// 				consts.UUID:   "RX_USER_ID_UUID",
// 				"primary":     true,
// 				"typeLabel":   "ID",
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "name",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_USER_NAME_UUID",
// 				"typeLabel":   "String",
// 				"nullable":    true,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "loginName",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_USER_LOGINNAME_UUID",
// 				"typeLabel":   "String",
// 				"length":      128,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "password",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_USER_PASSWORD_UUID",
// 				"typeLabel":   "String",
// 				"length":      256,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   consts.IS_SUPPER,
// 				consts.TYPE:   "Boolean",
// 				consts.UUID:   "RX_USER_ISSUPPER_UUID",
// 				"typeLabel":   "Boolean",
// 				"nullable":    true,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   consts.IS_DEMO,
// 				consts.TYPE:   "Boolean",
// 				consts.UUID:   "RX_USER_ISDEMO_UUID",
// 				"typeLabel":   "Boolean",
// 				"nullable":    true,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:       consts.META_CREATEDAT,
// 				consts.TYPE:       "Date",
// 				consts.UUID:       "RX_USER_CREATEDAT_UUID",
// 				"typeLabel":       "Date",
// 				consts.CREATEDATE: true,
// 				consts.SYSTEM:     true,
// 			},
// 			{
// 				consts.NAME:       consts.META_UPDATEDAT,
// 				consts.TYPE:       "Date",
// 				consts.UUID:       "RX_USER_UPDATEDAT_UUID",
// 				"typeLabel":       "Date",
// 				consts.UPDATEDATE: true,
// 				consts.SYSTEM:     true,
// 			},
// 		},
// 		"stereoType": "Entity",
// 	},
// 	{
// 		consts.NAME:    consts.META_ROLE,
// 		consts.UUID:    consts.ROLE_UUID,
// 		consts.INNERID: consts.ROLE_INNER_ID,
// 		consts.ROOT:    true,
// 		consts.SYSTEM:  true,
// 		"attributes": []map[string]interface{}{
// 			{
// 				consts.NAME:   "id",
// 				consts.TYPE:   "ID",
// 				consts.UUID:   "RX_ROLE_ID_UUID",
// 				"primary":     true,
// 				"typeLabel":   "ID",
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "name",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_ROLE_NAME_UUID",
// 				"typeLabel":   "String",
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:   "description",
// 				consts.TYPE:   "String",
// 				consts.UUID:   "RX_ROLE_DESCRIPTION_UUID",
// 				"typeLabel":   "String",
// 				"nullable":    true,
// 				consts.SYSTEM: true,
// 			},
// 			{
// 				consts.NAME:       consts.META_CREATEDAT,
// 				consts.TYPE:       "Date",
// 				consts.UUID:       "RX_ROLE_CREATEDAT_UUID",
// 				"typeLabel":       "Date",
// 				consts.CREATEDATE: true,
// 				consts.SYSTEM:     true,
// 			},
// 			{
// 				consts.NAME:       consts.META_UPDATEDAT,
// 				consts.TYPE:       "Date",
// 				consts.UUID:       "RX_ROLE_META_UPDATEDAT_UUID",
// 				"typeLabel":       "Date",
// 				consts.UPDATEDATE: true,
// 				consts.SYSTEM:     true,
// 			},
// 		},
// 		"stereoType": "Entity",
// 	},
// }

// var authRelations = []map[string]interface{}{
// 	{
// 		consts.UUID:          "META_RELATION_USER_ROLE_UUID",
// 		consts.INNERID:       consts.ROLE_USER_RELATION_INNER_ID,
// 		"sourceId":           consts.ROLE_UUID,
// 		"targetId":           consts.USER_UUID,
// 		"relationType":       "twoWayAssociation",
// 		"roleOfSource":       "roles",
// 		"roleOfTarget":       "users",
// 		"sourceMutiplicity":  "0..*",
// 		"targetMultiplicity": "0..*",
// 	},
// }

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
