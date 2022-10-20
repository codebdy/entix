package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/resolve"
	"rxdrag.com/entify/scalars"
	"rxdrag.com/entify/utils"
)

var installInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "InstallInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"meta": &graphql.InputObjectFieldConfig{
				Type: &graphql.NonNull{
					OfType: scalars.JSONType,
				},
			},
			consts.ADMIN: &graphql.InputObjectFieldConfig{
				Type: &graphql.NonNull{
					OfType: graphql.String,
				},
			},
			consts.ADMINPASSWORD: &graphql.InputObjectFieldConfig{
				Type: &graphql.NonNull{
					OfType: graphql.String,
				},
			},
			consts.WITHDEMO: &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
		},
	},
)

func MakeInstallSchema() *graphql.Schema {
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: consts.ROOT_QUERY_NAME,
				Fields: graphql.Fields{
					consts.INSTALLED: &graphql.Field{
						Type: graphql.Boolean,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							defer utils.PrintErrorStack()
							return false, nil
						},
					},
				},
			},
		),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: consts.ROOT_MUTATION_NAME,
			Fields: graphql.Fields{
				"install": &graphql.Field{
					Type: graphql.Boolean,
					Args: graphql.FieldConfigArgument{
						INPUT: &graphql.ArgumentConfig{
							Type: &graphql.NonNull{
								OfType: installInputType,
							},
						},
					},
					Resolve: resolve.InstallResolve,
				},
			},
			Description: "Root mutation of entity engine. For install auth entix",
		}),
	}
	theSchema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		panic(err)
		//log.Fatalf("failed to create new schema, error: %v", err)
	}

	return &theSchema
}
