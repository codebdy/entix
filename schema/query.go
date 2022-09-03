package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/resolve"
	"rxdrag.com/entify/utils"
)

func (a *AppSchema) rootQuery() *graphql.Object {
	rootQueryConfig := graphql.ObjectConfig{
		Name:   consts.ROOT_QUERY_NAME,
		Fields: a.queryFields(),
	}
	return graphql.NewObject(rootQueryConfig)
}

func (a *AppSchema) queryFields() graphql.Fields {
	queryFields := graphql.Fields{
		consts.INSTALLED: &graphql.Field{
			Type: graphql.Boolean,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				defer utils.PrintErrorStack()
				return Installed, nil
			},
		},
	}

	for _, intf := range a.model.Graph.RootInterfaces() {
		a.appendInterfaceToQueryFields(intf, queryFields)
	}

	for _, entity := range a.model.Graph.RootEnities() {
		a.appendEntityToQueryFields(entity, queryFields)
	}
	for _, service := range a.model.Graph.Services {
		a.appendServiceToQueryFields(service, queryFields)
	}
	// for _, service := range a.model.Graph.RootExternals() {
	// 	appendServiceQueryFields(service, queryFields)
	// }
	if repository.IsEntityExists(consts.META_USER) {
		a.appendAuthToQuery(queryFields)
	}
	return queryFields
}

func (a *AppSchema) QueryResponseType(class *graph.Class) graphql.Output {
	return a.modelParser.ClassListType(class)
}

func (a *AppSchema) appendInterfaceToQueryFields(intf *graph.Interface, fields graphql.Fields) {
	(fields)[intf.QueryName()] = &graphql.Field{
		Type:    a.QueryResponseType(&intf.Class),
		Args:    a.modelParser.QueryArgs(intf.Name()),
		Resolve: resolve.QueryInterfaceResolveFn(intf, a.Model()),
	}
	(fields)[intf.QueryOneName()] = &graphql.Field{
		Type:    a.modelParser.OutputType(intf.Name()),
		Args:    a.modelParser.QueryArgs(intf.Name()),
		Resolve: resolve.QueryOneInterfaceResolveFn(intf, a.Model()),
	}
}

func (a *AppSchema) appendEntityToQueryFields(entity *graph.Entity, fields graphql.Fields) {
	(fields)[entity.QueryName()] = &graphql.Field{
		Type:    a.QueryResponseType(&entity.Class),
		Args:    a.modelParser.QueryArgs(entity.Name()),
		Resolve: resolve.QueryEntityResolveFn(entity, a.Model()),
	}
	(fields)[entity.QueryOneName()] = &graphql.Field{
		Type:    a.modelParser.OutputType(entity.Name()),
		Args:    a.modelParser.QueryArgs(entity.Name()),
		Resolve: resolve.QueryOneEntityResolveFn(entity, a.Model()),
	}

	(fields)[entity.QueryAggregateName()] = &graphql.Field{
		Type:    a.modelParser.AggregateEntityType(entity),
		Args:    a.modelParser.QueryArgs(entity.Name()),
		Resolve: resolve.QueryEntityResolveFn(entity, a.Model()),
	}
}

func (a *AppSchema) appendServiceToQueryFields(partial *graph.Service, fields graphql.Fields) {

}

func (a *AppSchema) appendAuthToQuery(fields graphql.Fields) {
	fields[consts.ME] = &graphql.Field{
		Type:    a.modelParser.OutputType(consts.META_USER),
		Resolve: resolve.Me,
	}
}
