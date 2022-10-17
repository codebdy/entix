package schema

import (
	"errors"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/authentication"
	"rxdrag.com/entify/common/errorx"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/resolve"
	"rxdrag.com/entify/scalars"
	"rxdrag.com/entify/utils"
)

const INPUT = "input"

func (a *AppSchema) appendAuthMutation(fields graphql.Fields) {
	fields[consts.LOGIN] = &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			consts.LOGIN_NAME: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.String},
			},
			consts.PASSWORD: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.String},
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			defer utils.PrintErrorStack()
			return nil, errorx.New("001", "Can not login")
			//auth := authentication.New()
			//return auth.Login(p.Args[consts.LOGIN_NAME].(string), p.Args[consts.PASSWORD].(string))
		},
	}

	fields[consts.LOGOUT] = &graphql.Field{
		Type:    graphql.Boolean,
		Resolve: resolve.Logout,
	}
	fields[consts.CHANGE_PASSWORD] = &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			consts.LOGIN_NAME: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.String},
			},
			consts.OLD_PASSWORD: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.String},
			},
			consts.New_PASSWORD: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.String},
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			defer utils.PrintErrorStack()
			if p.Args[consts.LOGIN_NAME] == nil ||
				p.Args["oldPassword"] == nil ||
				p.Args["newPassword"] == nil {
				return "", errors.New("loginName, oldPassword or newPassword is emperty!")
			}
			auth := authentication.New()

			return auth.ChangePassword(p.Args[consts.LOGIN_NAME].(string),
				p.Args["oldPassword"].(string),
				p.Args["newPassword"].(string))
		},
	}
}

func (a *AppSchema) rootMutation() *graphql.Object {
	metaEntity := a.model.Graph.GetMetaEntity()
	mutationFields := graphql.Fields{}

	mutationFields[consts.UPLOAD] = &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			consts.ARG_FILE: &graphql.ArgumentConfig{
				Type: scalars.UploadType,
			},
		},
		Resolve: resolve.UploadResolveResolveFn(a.model.AppId),
	}

	mutationFields[consts.UPLOAD_PLUGIN] = &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			consts.ARG_FILE: &graphql.ArgumentConfig{
				Type: scalars.UploadType,
			},
		},
		Resolve: resolve.UploadPluginResolveResolveFn(a.model.AppId),
	}

	//if a.appUuid == consts.SYSTEM_APP_UUID {
	mutationFields[consts.PUBLISH] = &graphql.Field{
		Type: a.modelParser.OutputType(metaEntity.Name()),
		Args: graphql.FieldConfigArgument{
			consts.APPUUID: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.String},
			},
		},
		Resolve: publishResolve,
	}

	mutationFields[consts.DEPLOY_RPOCESS] = &graphql.Field{
		Type: graphql.ID,
		Args: graphql.FieldConfigArgument{
			consts.ID: &graphql.ArgumentConfig{
				Type: &graphql.NonNull{OfType: graphql.ID},
			},
		},
		Resolve: resolve.DeployProcessResolveFn(a.model),
	}
	//}

	for _, entity := range a.model.Graph.RootEnities() {
		if entity.Domain.Root {
			a.appendEntityMutationToFields(entity, mutationFields)
		}
	}

	for _, service := range a.model.Graph.Services {
		a.appendServiceMutationToFields(service, mutationFields)
	}

	a.appendAuthMutation(mutationFields)

	rootMutation := graphql.ObjectConfig{
		Name:        consts.ROOT_MUTATION_NAME,
		Fields:      mutationFields,
		Description: "Root mutation of entity engine.",
	}

	return graphql.NewObject(rootMutation)
}

func (a *AppSchema) deleteArgs(entity *graph.Entity) graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		consts.ARG_WHERE: &graphql.ArgumentConfig{
			Type: a.modelParser.WhereExp(entity.Name()),
		},
	}
}

func deleteByIdArgs() graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		consts.ID: &graphql.ArgumentConfig{
			Type: graphql.ID,
		},
	}
}

func (a *AppSchema) upsertArgs(entity *graph.Entity) graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		consts.ARG_OBJECTS: &graphql.ArgumentConfig{
			Type: &graphql.NonNull{
				OfType: &graphql.List{
					OfType: &graphql.NonNull{
						OfType: a.modelParser.SaveInput(entity.Name()),
					},
				},
			},
		},
	}
}

func (a *AppSchema) upsertOneArgs(entity *graph.Entity) graphql.FieldConfigArgument {
	return graphql.FieldConfigArgument{
		consts.ARG_OBJECT: &graphql.ArgumentConfig{
			Type: &graphql.NonNull{
				OfType: a.modelParser.SaveInput(entity.Name()),
			},
		},
	}
}

func (a *AppSchema) setArgs(entity *graph.Entity) graphql.FieldConfigArgument {
	updateInput := a.modelParser.SetInput(entity.Name())
	return graphql.FieldConfigArgument{
		consts.ARG_SET: &graphql.ArgumentConfig{
			Type: &graphql.NonNull{
				OfType: updateInput,
			},
		},
		consts.ARG_WHERE: &graphql.ArgumentConfig{
			Type: a.modelParser.WhereExp(entity.Name()),
		},
	}
}

func (a *AppSchema) appendEntityMutationToFields(entity *graph.Entity, feilds graphql.Fields) {
	(feilds)[entity.DeleteName()] = &graphql.Field{
		Type:    a.modelParser.MutationResponse(entity.Name()),
		Args:    a.deleteArgs(entity),
		Resolve: resolve.DeleteResolveFn(entity, a.Model()),
	}
	(feilds)[entity.DeleteByIdName()] = &graphql.Field{
		Type:    a.modelParser.OutputType(entity.Name()),
		Args:    deleteByIdArgs(),
		Resolve: resolve.DeleteByIdResolveFn(entity, a.Model()),
	}
	(feilds)[entity.UpsertName()] = &graphql.Field{
		Type:    &graphql.List{OfType: a.modelParser.OutputType(entity.Name())},
		Args:    a.upsertArgs(entity),
		Resolve: resolve.PostResolveFn(entity, a.Model()),
	}
	(feilds)[entity.UpsertOneName()] = &graphql.Field{
		Type:    a.modelParser.OutputType(entity.Name()),
		Args:    a.upsertOneArgs(entity),
		Resolve: resolve.PostOneResolveFn(entity, a.Model()),
	}

	updateInput := a.modelParser.SetInput(entity.Name())
	if len(updateInput.Fields()) > 0 {
		(feilds)[entity.SetName()] = &graphql.Field{
			Type:    a.modelParser.MutationResponse(entity.Name()),
			Args:    a.setArgs(entity),
			Resolve: resolve.SetResolveFn(entity, a.Model()),
		}
	}
}

func (a *AppSchema) appendServiceMutationToFields(service *graph.Service, feilds graphql.Fields) {

	// (feilds)[service.DeleteName()] = &graphql.Field{
	// 	Type: a.modelParser.MutationResponse(service.Name()),
	// 	Args: a.deleteArgs(&service.Entity),
	// 	//Resolve: entity.QueryResolve(),
	// }
	// (feilds)[service.DeleteByIdName()] = &graphql.Field{
	// 	Type: a.modelParser.OutputType(service.Name()),
	// 	Args: deleteByIdArgs(),
	// 	//Resolve: entity.QueryResolve(),
	// }
}
