package parser

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
)

func (p *ModelParser) makeInputs() {
	for i := range p.model.Graph.Entities {
		entity := p.model.Graph.Entities[i]
		p.setInputMap[entity.Name()] = p.makeEntitySetInput(entity)
		p.saveInputMap[entity.Name()] = p.makeEntitySaveInput(entity)
		p.mutationResponseMap[entity.Name()] = p.makeEntityMutationResponseType(entity)
	}
	p.makeEntityInputRelations()
}

func (p *ModelParser) makeEntityInputRelations() {
	for i := range p.model.Graph.Entities {
		entity := p.model.Graph.Entities[i]

		input := p.setInputMap[entity.Name()]
		update := p.saveInputMap[entity.Name()]

		associas := entity.AllAssociations()

		for i := range associas {
			assoc := associas[i]

			typeInput := p.SaveInput(assoc.Owner().Name())
			if typeInput == nil {
				panic("can not find save input:" + assoc.Owner().Name())
			}
			if len(typeInput.Fields()) == 0 {
				fmt.Println("Fields == 0")
				continue
			}

			arrayType := p.getAssociationType(assoc)

			if arrayType == nil {
				panic("Can not get association type:" + assoc.Owner().Name() + "." + assoc.Name())
			}
			input.AddFieldConfig(assoc.Name(), &graphql.InputObjectFieldConfig{
				Type:        arrayType,
				Description: assoc.Description(),
			})
			update.AddFieldConfig(assoc.Name(), &graphql.InputObjectFieldConfig{
				Type:        arrayType,
				Description: assoc.Description(),
			})
		}
	}
}

func (p *ModelParser) getAssociationType(association *graph.Association) *graphql.List {
	//if association.IsArray() {
	typeInput := p.SaveInput(association.TypeEntity().Name())
	return &graphql.List{
		OfType: typeInput,
	}
	//return p.HasManyInput(association.TypeEntity().Name())
	//}
	// } else {
	// 	return p.HasOneInput(association.TypeEntity().Name())
	// }
}

func (p *ModelParser) inputFields(entity *graph.Entity, withId bool) graphql.InputObjectConfigFieldMap {
	fields := graphql.InputObjectConfigFieldMap{}
	for _, column := range entity.AllAttributes() {
		if (column.Name != consts.ID || withId) && !column.DeleteDate && !column.CreateDate && !column.UpdateDate {
			fields[column.Name] = &graphql.InputObjectFieldConfig{
				Type:        p.InputPropertyType(column),
				Description: column.Description,
			}
		}
	}
	return fields
}

func (p *ModelParser) makeEntitySaveInput(entity *graph.Entity) *graphql.InputObject {
	name := entity.Name() + consts.INPUT
	return graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:   name,
			Fields: p.inputFields(entity, true),
		},
	)
}

func (p *ModelParser) makeEntitySetInput(entity *graph.Entity) *graphql.InputObject {
	return graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name:   entity.Name() + consts.SET_INPUT,
			Fields: p.inputFields(entity, false),
		},
	)
}

func (p *ModelParser) makeEntityMutationResponseType(entity *graph.Entity) *graphql.Object {
	var returnValue *graphql.Object

	returnValue = graphql.NewObject(
		graphql.ObjectConfig{
			Name: entity.Name() + consts.MUTATION_RESPONSE,
			Fields: graphql.Fields{
				consts.RESPONSE_AFFECTEDROWS: &graphql.Field{
					Type: graphql.Int,
				},
				consts.RESPONSE_RETURNING: &graphql.Field{
					Type: &graphql.NonNull{
						OfType: &graphql.List{
							OfType: p.OutputType(entity.Name()),
						},
					},
				},
			},
		},
	)

	return returnValue
}
