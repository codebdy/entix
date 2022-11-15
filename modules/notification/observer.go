package notification

import (
	"log"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/common/contexts"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/observer"
	"rxdrag.com/entify/orm"
)

const EntityNotificationName = "Notification"

type NotificationObserver struct {
	channel chan (interface{})
	key     string
	p       graphql.ResolveParams
	model   *model.Model
}

func newObserver(p graphql.ResolveParams, model *model.Model) *NotificationObserver {
	ntObserver := &NotificationObserver{
		channel: make(chan interface{}),
		key:     uuid.New().String(),
		p:       p,
		model:   model,
	}
	observer.AddObserver(ntObserver)

	return ntObserver
}
func (o *NotificationObserver) Key() string {
	return o.key
}

func (o *NotificationObserver) ObjectCreated(object map[string]interface{}, entityName string) {
	if entityName == EntityNotificationName {
		o.calculateCounts(object)
	}
}
func (o *NotificationObserver) ObjectUpdated(object map[string]interface{}, entityName string) {
	if entityName == EntityNotificationName {
		o.calculateCounts(object)
	}
}
func (o *NotificationObserver) ObjectDeleted(object map[string]interface{}, entityName string) {
	if entityName == EntityNotificationName {
		o.calculateCounts(object)
	}
}

func (o *NotificationObserver) calculateCounts(object map[string]interface{}) {
	entity := o.model.Graph.GetEntityByName(EntityNotificationName)
	if entity == nil {
		log.Panic("Can find entity Notification")
	}

	me := contexts.Values(o.p.Context).Me
	appId := contexts.Values(o.p.Context).AppId

	if me == nil || appId == 0 {
		log.Panic("User or app not set!")
	}
	session, err := orm.Open()
	if err != nil {
		log.Panic(err.Error())
	}

	if object["user"] == nil {
		log.Panic()
	}

	result := session.Query(entity, map[string]interface{}{}, []*graph.Attribute{})

	o.channel <- result.Total
}

func (o *NotificationObserver) destory() {
	close(o.channel)
	observer.RemoveObserver(o.key)
}
