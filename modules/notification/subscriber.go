package notification

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/orm"
)

type Subscriber struct {
	key     string
	channel chan (interface{})
	p       graphql.ResolveParams
	model   *model.Model
}

func newSubscriber(p graphql.ResolveParams, model *model.Model) *Subscriber {
	s := &Subscriber{
		key:     uuid.New().String(),
		channel: make(chan interface{}),
		p:       p,
		model:   model,
	}
	NoticeModelObserver.addSubscriber(s)
	return s
}

func (s *Subscriber) notificationChanged(notification map[string]interface{}, ctx context.Context) {
	me := contexts.Values(ctx).Me
	appId := contexts.Values(ctx).AppId

	if me == nil || appId == 0 {
		log.Panic("User or app not set!")
	}
	session, err := orm.Open()
	if err != nil {
		log.Panic(err.Error())
	}

	if notification["user"] == nil {
		log.Panic("Notification no user")
	}

	result := session.Query(
		s.model.Graph.GetEntityByName(EntityNotificationName),
		map[string]interface{}{},
		[]*graph.Attribute{},
	)
	s.channel <- result.Total
}

func (s *Subscriber) notificationDeleted(ctx context.Context) {

}

func (s *Subscriber) destory() {
	close(s.channel)
	NoticeModelObserver.deleteSubscriber(s.key)
}
