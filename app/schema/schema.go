package schema

import (
	"rxdrag.com/entify/app/schema/parser"
	"rxdrag.com/entify/model"
)

type AppSchema struct {
	Model       *model.Model
	modelParser parser.ModelParser
}

func New(model *model.Model) *AppSchema {
	appSchema := &AppSchema{
		Model: model,
	}

	appSchema.modelParser.ParseModel(model)
	return appSchema
}
