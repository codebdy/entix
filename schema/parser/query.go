package parser

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
)

func (p *ModelParser) ClassListType(entity *graph.Class) *graphql.Object {
	return p.listType(entity.Name(), entity.ListName())
}

func (p *ModelParser) listType(name string, listName string) *graphql.Object {
	if p.listMap[listName] != nil {
		return p.listMap[listName]
	}

	returnValue := graphql.NewObject(
		graphql.ObjectConfig{
			Name: listName,
			Fields: graphql.Fields{
				consts.NODES: &graphql.Field{
					Type: &graphql.List{
						OfType: p.OutputType(name),
					},
				},
				consts.TOTAL: &graphql.Field{
					Type: graphql.Int,
				},
			},
		},
	)

	p.listMap[listName] = returnValue
	return returnValue
}

func (p *ModelParser) makeEntityOutputObjects(entities []*graph.Entity) {
	for i := range entities {
		p.makeEntityObject(entities[i])
	}
}

func (p *ModelParser) makeEntityObject(entity *graph.Entity) {
	objType := p.ObjectType(entity)
	p.objectTypeMap[entity.Name()] = objType
	p.objectMapById[entity.InnerId()] = objType
}

func (p *ModelParser) ObjectType(entity *graph.Entity) *graphql.Object {
	name := entity.Name()
	interfaces := p.mapInterfaces(entity.Interfaces)

	if len(interfaces) > 0 {
		return graphql.NewObject(
			graphql.ObjectConfig{
				Name:        name,
				Fields:      p.outputFields(entity.AllAttributes(), entity.AllMethods()),
				Description: entity.Description(),
				Interfaces:  interfaces,
			},
		)
	} else {
		return graphql.NewObject(
			graphql.ObjectConfig{
				Name:        name,
				Fields:      p.outputFields(entity.AllAttributes(), entity.AllMethods()),
				Description: entity.Description(),
			},
		)
	}

}

func (p *ModelParser) outputFields(attrs []*graph.Attribute, methods []*graph.Method) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		fields[attr.Name] = &graphql.Field{
			Type:        p.PropertyType(attr),
			Description: attr.Description,
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	fmt.Println(p.Context.Value("data"))
			// 	return "world", nil
			// },
		}
	}

	for _, method := range methods {
		fields[method.GetName()] = &graphql.Field{
			Type:        p.PropertyType(method),
			Description: method.Method.Description,
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	fmt.Println(p.Context.Value("data"))
			// 	return "world", nil
			// },
		}
	}
	return fields
}
