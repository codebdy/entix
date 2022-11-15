package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model"
)

type Subscriber struct {
	key     string
	channel <-chan (interface{})
	p       graphql.ResolveParams
	model   *model.Model
}

func newSubscriber(p graphql.ResolveParams, model *model.Model) *Subscriber {
	s := &Subscriber{
		key:     uuid.New().String(),
		channel: make(<-chan interface{}),
		p:       p,
		model:   model,
	}
	NoticeModelObserver.addSubscriber(s)
	return s
}

func (s *Subscriber) notificationChanged(notification map[string]interface{}) {
	// me := contexts.Values(o.p.Context).Me
	// appId := contexts.Values(o.p.Context).AppId

	// if me == nil || appId == 0 {
	// 	log.Panic("User or app not set!")
	// }
	// session, err := orm.Open()
	// if err != nil {
	// 	log.Panic(err.Error())
	// }

	// if object["user"] == nil {
	// 	log.Panic()
	// }

	//result := session.Query(entity, map[string]interface{}{}, []*graph.Attribute{})

}

func (s *Subscriber) notificationDeleted(ctx context.Context) {

}

func (s *Subscriber) destory() {
	NoticeModelObserver.delteSubscriber(s.key)
}
