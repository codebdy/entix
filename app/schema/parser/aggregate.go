package parser

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
)

func (p *ModelParser) avgFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, column := range attrs {
		if column.Type == meta.INT || column.Type == meta.FLOAT {
			fields[column.Name] = &graphql.Field{
				Type: p.PropertyType(column.Type),
			}
		}

	}
	return fields
}

func (p *ModelParser) maxFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
			}
		}

	}
	return fields
}

func (p *ModelParser) minFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}

	}
	return fields
}

func (p *ModelParser) selectFields(attrs []*graph.Attribute) graphql.InputObjectConfigFieldMap {
	fields := graphql.InputObjectConfigFieldMap{}
	for _, attr := range attrs {
		if attr.Type != meta.FILE {
			fields[attr.Name] = &graphql.InputObjectFieldConfig{
				Type: p.InputPropertyType(attr),
			}
		}
	}

	return fields
}

func (p *ModelParser) stddevFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}

	}
	return fields
}

func (p *ModelParser) stddevPopFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}
	}
	return fields
}

func (p *ModelParser) stddevSampFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}

	}
	return fields
}

func (p *ModelParser) sumFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}

	}
	return fields
}

func (p *ModelParser) varPopFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}

	}
	return fields
}

func (p *ModelParser) varSampFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}

	}
	return fields
}

func (p *ModelParser) varianceFields(attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	for _, attr := range attrs {
		if attr.Type == meta.INT || attr.Type == meta.FLOAT {
			fields[attr.Name] = &graphql.Field{
				Type: p.PropertyType(attr.Type),
				// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// 	fmt.Println(p.Context.Value("data"))
				// 	return "world", nil
				// },
			}
		}

	}
	return fields
}

func (p *ModelParser) AggregateFields(name string, attrs []*graph.Attribute) graphql.Fields {
	fields := graphql.Fields{}
	avgFields := p.avgFields(attrs)
	if len(avgFields) > 0 {
		fields["avg"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "AvgFields",
					Fields: avgFields,
				},
			),
		}
	}

	maxFields := p.maxFields(attrs)
	if len(maxFields) > 0 {
		fields["max"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "MaxFields",
					Fields: maxFields,
				},
			),
		}
	}

	minFields := p.minFields(attrs)
	if len(minFields) > 0 {
		fields["min"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "MinFields",
					Fields: minFields,
				},
			),
		}
	}

	countFields := p.selectFields(attrs)
	if len(countFields) > 0 {
		selectColumnName := name + "SelectColumn"
		selectColumn := graphql.NewInputObject(
			graphql.InputObjectConfig{
				Name:   selectColumnName,
				Fields: countFields,
			},
		)
		p.selectColumnsMap[selectColumnName] = selectColumn
		fields[consts.ARG_COUNT] = &graphql.Field{
			Args: graphql.FieldConfigArgument{
				consts.ARG_COLUMNS: &graphql.ArgumentConfig{
					Type: selectColumn,
				},
				consts.ARG_DISTINCT: &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
			Type: graphql.Int,
		}
	}

	stddevFields := p.stddevFields(attrs)
	if len(stddevFields) > 0 {
		fields["stddev"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "StddevFields",
					Fields: stddevFields,
				},
			),
		}
	}

	stddevPopFields := p.stddevPopFields(attrs)
	if len(stddevPopFields) > 0 {
		fields["stddevPop"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "StddevPopFields",
					Fields: stddevPopFields,
				},
			),
		}
	}

	stddevSampFields := p.stddevSampFields(attrs)
	if len(stddevSampFields) > 0 {
		fields["stddevSamp"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "StddevSampFields",
					Fields: stddevSampFields,
				},
			),
		}
	}

	sumFields := p.sumFields(attrs)
	if len(sumFields) > 0 {
		fields["sum"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "SumFields",
					Fields: sumFields,
				},
			),
		}
	}
	varPopFields := p.varPopFields(attrs)
	if len(varPopFields) > 0 {
		fields["varPop"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "VarPopFields",
					Fields: varPopFields,
				},
			),
		}
	}
	varSampFields := p.varSampFields(attrs)
	if len(varSampFields) > 0 {
		fields["varSamp"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "VarSampFields",
					Fields: varSampFields,
				},
			),
		}
	}
	varianceFields := p.varianceFields(attrs)
	if len(varianceFields) > 0 {
		fields["variance"] = &graphql.Field{
			Type: graphql.NewObject(
				graphql.ObjectConfig{
					Name:   name + "VarianceFields",
					Fields: varianceFields,
				},
			),
		}
	}
	return fields
}

func (p *ModelParser) aggregateType(entity *graph.Entity) *graphql.Object {
	aggregateName := entity.AggregateName()
	if p.aggregateMap[aggregateName] != nil {
		return p.aggregateMap[aggregateName]
	}
	aggregateFields := p.AggregateFields(entity.Name(), entity.AllAttributes())

	obj := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   aggregateName + consts.FIELDS,
			Fields: aggregateFields,
		},
	)
	p.aggregateMap[aggregateName] = obj
	return obj
}
