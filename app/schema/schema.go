package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/schema/parser"
)

type AppSchema struct {
	appUuid     string
	model       *model.Model
	modelParser parser.ModelParser
	schema      *graphql.Schema
}

func NewAppSchema(appUuid string) *AppSchema {
	appSchema := &AppSchema{
		appUuid: appUuid,
	}

	if !Installed && appUuid != consts.SYSTEM_APP_UUID {
		panic("Server is not installed, please install first")
	}

	if !Installed {
		appSchema.schema = MakeInstallSchema()
	} else {
		appSchema.Make()
	}

	return appSchema
}
func (s *AppSchema) Make() {
	if s.model == nil {
		//第一步初始值，用于取meta信息，取完后，换掉该部分内容
		initMeta := meta.MetaContent{
			Classes: []meta.ClassMeta{
				meta.MetaStatusEnum,
				meta.MetaClass,
				meta.EntityAuthSettingsClass,
				meta.AbilityTypeEnum,
				meta.AbilityClass,
			},
		}
		s.model = model.New(s.appUuid, &initMeta)
	}
	repos := repository.New(s.Model())
	repos.MakeSupperVerifier()

	//第二步， 取系统应用，以此为基础，进一步取数据。
	//   保证取APPID时，能拿到App实体
	s.model = repos.LoadModel(consts.SYSTEM_APP_UUID)
	repos = repository.New(s.Model())
	repos.MakeSupperVerifier()

	//第三步，加载真正的Model
	s.model = repos.LoadModel(s.appUuid)
	s.schema = s.doMake()
}
func (s *AppSchema) Model() *model.Model {
	return s.model
}

func (s *AppSchema) Schema() *graphql.Schema {
	return s.schema
}

func (s *AppSchema) doMake() *graphql.Schema {
	s.modelParser.ParseModel(s.model)

	schemaConfig := graphql.SchemaConfig{
		Query:    s.rootQuery(),
		Mutation: s.rootMutation(),
		Directives: []*graphql.Directive{
			graphql.NewDirective(graphql.DirectiveConfig{
				Name:      "forEdit",
				Locations: []string{graphql.DirectiveLocationField},
			}),
		},
		Types: append(s.modelParser.EntityTypes()),
	}
	theSchema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		panic(err)
		//log.Fatalf("failed to create new schema, error: %v", err)
	}
	return &theSchema
}
