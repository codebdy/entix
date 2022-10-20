package entry

import (
	"context"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

const (
	ROOT_QUERY_NAME        = "Query"
	ROOT_MUTATION_NAME     = "Mutation"
	ROOT_SUBSCRIPTION_NAME = "Subscription"
)

var modules []Moduler = []Moduler{}

type Moduler interface {
	QueryFields(ctx context.Context) []*graphql.Field
	MutationFields(ctx context.Context) []*graphql.Field
	SubscriptionFields(ctx context.Context) []*graphql.Field
	Directives(ctx context.Context) []*graphql.Directive
	Types(ctx context.Context) []graphql.Type
	Middlewares() []func(next http.Handler) http.Handler
}

func AddModuler(module Moduler) {
	modules = append(modules, module)
}

func GetSchema(ctx context.Context) graphql.Schema {
	rootQueryFields := graphql.Fields{}
	rootMutationFields := graphql.Fields{}
	rootSubscriptionFields := graphql.Fields{}
	directives := []*graphql.Directive{}
	types := []graphql.Type{}

	for _, module := range modules {
		queryFields := module.QueryFields(ctx)
		for i := range queryFields {
			field := queryFields[i]
			rootQueryFields[field.Name] = field
		}
		mutationFields := module.MutationFields(ctx)
		for i := range mutationFields {
			field := mutationFields[i]
			rootMutationFields[field.Name] = field
		}

		subscriptionFields := module.SubscriptionFields(ctx)
		for i := range subscriptionFields {
			field := subscriptionFields[i]
			rootSubscriptionFields[field.Name] = field
		}
		directives = append(directives, module.Directives(ctx)...)
		types = append(types, module.Types(ctx)...)
	}

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   ROOT_QUERY_NAME,
		Fields: rootQueryFields,
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   ROOT_MUTATION_NAME,
		Fields: rootMutationFields,
	})

	if len(rootMutationFields) == 0 {
		rootMutation = nil
	}

	rootSubscription := graphql.NewObject(graphql.ObjectConfig{
		Name:   ROOT_SUBSCRIPTION_NAME,
		Fields: rootSubscriptionFields,
	})

	if len(rootSubscriptionFields) == 0 {
		rootSubscription = nil
	}

	schemaConfig := graphql.SchemaConfig{
		Query:        rootQuery,
		Mutation:     rootMutation,
		Subscription: rootSubscription,
		Directives:   directives,
		Types:        types,
	}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		log.Panic(err)
	}

	return schema
}

func AppendMiddlewares(h http.Handler) http.Handler {
	for _, model := range modules {
		middlewares := model.Middlewares()
		for i := range middlewares {
			h = middlewares[i](h)
		}
	}
	return h
}
