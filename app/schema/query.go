package schema

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/app/resolve"
	"rxdrag.com/entify/app/schema/parser"
	"rxdrag.com/entify/model/graph"
)

func (a *AppProcessor) QueryFields() []*graphql.Field {
	queryFields := graphql.Fields{}

	for _, intf := range a.Model.Graph.RootInterfaces() {
		a.appendInterfaceToQueryFields(intf, queryFields)
	}

	for _, entity := range a.Model.Graph.RootEnities() {
		a.appendEntityToQueryFields(entity, queryFields)
	}
	for _, third := range a.Model.Graph.ThirdParties {
		a.appendThirdPartyToQueryFields(third, queryFields)
	}
	for _, service := range a.Model.Graph.Services {
		a.appendServiceToQueryFields(service, queryFields)
	}

	return convertFieldsArray(queryFields)
}

func (a *AppProcessor) EntityQueryResponseType(entity *graph.Entity) graphql.Output {
	return a.modelParser.EntityListType(entity)
}
func (a *AppProcessor) ClassQueryResponseType(cls *graph.Class) graphql.Output {
	return a.modelParser.ClassListType(cls)
}

func (a *AppProcessor) appendInterfaceToQueryFields(intf *graph.Interface, fields graphql.Fields) {
	// (fields)[intf.QueryName()] = &graphql.Field{
	// 	Type:    a.QueryResponseType(&intf.Class),
	// 	Args:    a.modelParser.QueryArgs(intf.Name()),
	// 	Resolve: resolve.QueryInterfaceResolveFn(intf, a.Model),
	// }
	(fields)[intf.QueryOneName()] = &graphql.Field{
		Type:    a.modelParser.OutputType(intf.Name()),
		Args:    a.modelParser.QueryArgs(intf.Name()),
		Resolve: resolve.QueryOneInterfaceResolveFn(intf, a.Model),
	}
}

func (a *AppProcessor) appendEntityToQueryFields(entity *graph.Entity, fields graphql.Fields) {
	(fields)[entity.QueryName()] = &graphql.Field{
		Type:    a.EntityQueryResponseType(entity),
		Args:    a.modelParser.QueryArgs(entity.Name()),
		Resolve: resolve.QueryEntityResolveFn(entity, a.Model),
	}
	(fields)[entity.QueryOneName()] = &graphql.Field{
		Type:    a.modelParser.OutputType(entity.Name()),
		Args:    a.modelParser.QueryArgs(entity.Name()),
		Resolve: resolve.QueryOneEntityResolveFn(entity, a.Model),
	}

}

func (a *AppProcessor) appendThirdPartyToQueryFields(third *graph.ThirdParty, fields graphql.Fields) {
	(fields)[third.QueryName()] = &graphql.Field{
		Type:    a.ClassQueryResponseType(&third.Class),
		Args:    a.modelParser.QueryArgs(third.Name()),
		Resolve: resolve.QueryThirdPartyResolveFn(third, a.Model),
	}
	(fields)[third.QueryOneName()] = &graphql.Field{
		Type:    a.modelParser.OutputType(third.Name()),
		Args:    a.modelParser.QueryArgs(third.Name()),
		Resolve: resolve.QueryOneThirdPartyResolveFn(third, a.Model),
	}

}

func (a *AppProcessor) appendServiceToQueryFields(service *graph.Service, fields graphql.Fields) {
	for _, method := range service.QueryMethods() {
		fields[service.Name()+"_"+method.GetName()] = &graphql.Field{
			Type:        parser.PropertyType(method.GetType()),
			Args:        a.modelParser.MethodArgs(method),
			Description: method.Method.Description,
			Resolve:     resolve.MethodResolveFn(method, a.Model),
		}
	}
}
