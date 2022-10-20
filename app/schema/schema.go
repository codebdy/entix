package schema

import (
	"github.com/graphql-go/graphql"
)

func MakeSchema() *graphql.Schema {
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
