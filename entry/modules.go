package entry

import (
	"context"

	"github.com/graphql-go/graphql"
)

var modules []Moduler = []Moduler{}

type Moduler interface {
	QueryFields(ctx context.Context) []*graphql.Field
	MutationFields(ctx context.Context) []*graphql.Field
	SubscriptionFields(ctx context.Context) []*graphql.Field
}

func GetSchema() *graphql.Schema {

}
