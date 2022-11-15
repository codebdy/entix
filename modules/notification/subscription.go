package notification

import (
	"github.com/graphql-go/graphql"
)

func (m *SubscriptionModule) SubscriptionFields() []*graphql.Field {
	if m.app != nil {
		return []*graphql.Field{
			{
				Name: "unreadNoticationCounts",
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source, nil
				},
				Subscribe: func(p graphql.ResolveParams) (interface{}, error) {
					subscrber := newSubscriber(p, m.app.Model)
					go func() {
						<-p.Context.Done()
						subscrber.destory()
						return
					}()

					return subscrber.channel, nil
				},
			},
		}
	} else {
		return []*graphql.Field{}
	}
}
