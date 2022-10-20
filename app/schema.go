package app

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/schema"
)

//本函数要重写
func NewAppSchema(appId uint64) *App {
	appSchema := &App{
		AppId: appId,
	}

	// if !Installed && appId != SYSTEM_APP_ID {
	// 	panic("Server is not installed, please install first")
	// }

	if !Installed {
		appSchema.schema = schema.MakeInstallSchema()
	} else {
		appSchema.Make()
	}

	return appSchema
}
func (a *App) Make() {
	if a.model == nil {
		//第一步初始值，用于取meta信息，取完后，换掉该部分内容
		// initMeta := meta.MetaContent{
		// 	Classes: []meta.ClassMeta{
		// 		meta.MetaStatusEnum,
		// 		meta.MetaClass,
		// 		meta.EntityAuthSettingsClass,
		// 		meta.AbilityTypeEnum,
		// 		meta.AbilityClass,
		// 	},
		// }
		//s.model = model.New(s.appId, &initMeta)
	}
	repos := repository.New(a.Model())
	repos.MakeSupperVerifier()

	//第二步， 取系统应用，以此为基础，进一步取数据。
	//   保证取APPID时，能拿到App实体
	//s.model = repos.LoadModel(consts.SYSTEM_APP_UUID)
	repos = repository.New(a.Model())
	repos.MakeSupperVerifier()

	//第三步，加载真正的Model
	//s.model = repos.LoadModel(s.appId)
	a.schema = a.doMake()
}
func (a *App) Model() *model.Model {
	return a.model
}

func (a *App) Schema() *graphql.Schema {
	return a.schema
}

func (a *App) doMake() *graphql.Schema {
	//s.modelParser.ParseModel(s.model)

	schemaConfig := graphql.SchemaConfig{
		//Query:    s.rootQuery(),
		//Mutation: s.rootMutation(),
		Directives: []*graphql.Directive{
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "forEdit",
				Locations: []string{graphql.DirectiveLocationField},
			}),
		},
		//Types: append(s.modelParser.EntityTypes()),
	}
	theSchema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		panic(err)
		//log.Fatalf("failed to create new schema, error: %v", err)
	}
	return &theSchema
}
