package notification

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model/observer"
)

type NotificationObserver struct {
	channel chan (interface{})
	key     string
	p       graphql.ResolveParams
}

func newObserver(p graphql.ResolveParams) *NotificationObserver {
	ntObserver := &NotificationObserver{
		channel: make(chan interface{}),
		key:     uuid.New().String(),
		p:       p,
	}
	observer.AddObserver(ntObserver)

	return ntObserver
}
func (o *NotificationObserver) Key() string {
	return o.key
}

func (o *NotificationObserver) ObjectCreated(object map[string]interface{}) {
	o.calculateCounts(object)
}
func (o *NotificationObserver) ObjectUpdated(object map[string]interface{}) {
	o.calculateCounts(object)
}
func (o *NotificationObserver) ObjectDeleted(object map[string]interface{}) {
	o.calculateCounts(object)
}

func (o *NotificationObserver) calculateCounts(object map[string]interface{}) {
	//o.channel <- 0
}

func (o *NotificationObserver) destory() {
	close(o.channel)
	observer.RemoveObserver(o.key)
}
