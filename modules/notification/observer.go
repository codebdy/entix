package notification

import (
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model/observer"
)

type NotificationObserver struct {
	c   chan (interface{})
	key string
	p   graphql.ResolveParams
}

func newObserver(p graphql.ResolveParams) *NotificationObserver {
	ntObserver := &NotificationObserver{
		c:   make(chan interface{}),
		key: uuid.New().String(),
		p:   p,
	}
	observer.AddObserver(ntObserver)

	return ntObserver
}
func (o *NotificationObserver) Key() string {
	return o.key
}

func (o *NotificationObserver) ObjectCreated(object map[string]interface{}) {

}
func (o *NotificationObserver) ObjectUpdated(object map[string]interface{}) {

}
func (o *NotificationObserver) ObjectDeleted(object map[string]interface{}) {

}

func (o *NotificationObserver) destory() {
	close(o.c)
	observer.RemoveObserver(o.key)
}
