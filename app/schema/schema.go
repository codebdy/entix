package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app/schema/parser"
	"rxdrag.com/entify/model"
)

type AppGraphqlSchema struct {
	QueryFields        []*graphql.Field
	MutationFields     []*graphql.Field
	SubscriptionFields []*graphql.Field
	Directives         []*graphql.Directive
	Types              []graphql.Type
}

type AppProcessor struct {
	Model       *model.Model
	modelParser parser.ModelParser
}

func New(model *model.Model) AppGraphqlSchema {
	processor := &AppProcessor{
		Model: model,
	}

	processor.modelParser.ParseModel(model)
	return AppGraphqlSchema{
		QueryFields:    processor.QueryFields(),
		MutationFields: processor.mutationFields(),
		Types:          processor.modelParser.EntityTypes(),
	}
}
