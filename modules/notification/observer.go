package notification

import (
	"log"
	"sync"

	"rxdrag.com/entify/model/observer"
	"rxdrag.com/entify/modules/app"
)

const EntityNotificationName = "Notification"

var NoticeModelObserver *NotificationObserver

type NotificationObserver struct {
	key         string
	subscribers sync.Map
}

func init() {
	//创建模型监听器
	NoticeModelObserver = &NotificationObserver{
		key: "NotificationObserver",
	}
	observer.AddObserver(NoticeModelObserver)
}

func (o *NotificationObserver) Key() string {
	return o.key
}

func (o *NotificationObserver) ObjectPosted(object map[string]interface{}, entityName string, userId, appId uint64) {
	if entityName == EntityNotificationName {
		o.distributeChanged(object)
	}
}
func (o *NotificationObserver) ObjectMultiPosted(objects []map[string]interface{}, entityName string, userId, appId uint64) {
	if entityName == EntityNotificationName {
		for _, object := range objects {
			o.distributeChanged(object)
		}
	}
}
func (o *NotificationObserver) ObjectDeleted(object map[string]interface{}, entityName string, userId, appId uint64) {
	if entityName == EntityNotificationName {
		o.distributeDeleted(userId, appId)
	}
}

func (o *NotificationObserver) ObjectMultiDeleted(objects []map[string]interface{}, entityName string, userId, appId uint64) {
	if entityName == EntityNotificationName {
		o.distributeDeleted(userId, appId)
	}
}

func (o *NotificationObserver) isEmperty() bool {
	emperty := true
	o.subscribers.Range(func(key interface{}, value interface{}) bool {
		emperty = false
		return true
	})
	return emperty
}

//分发详细信息到各订阅者
func (o *NotificationObserver) distributeChanged(object map[string]interface{}) {
	if o.isEmperty() {
		return
	}
	model := app.GetSystemApp().Model
	entity := model.Graph.GetEntityByName(EntityNotificationName)
	if entity == nil {
		log.Panic("Can find entity Notification")
	}

	//补全信息
	newObject := object

	//分发
	o.subscribers.Range(func(key interface{}, value interface{}) bool {
		value.(*Subscriber).notificationChanged(newObject)
		return true
	})

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

func (o *NotificationObserver) distributeDeleted(userId, appId uint64) {
	o.subscribers.Range(func(key interface{}, value interface{}) bool {
		value.(*Subscriber).notificationDeleted(userId, appId)
		return true
	})
}

func (o *NotificationObserver) addSubscriber(s *Subscriber) {
	o.subscribers.Store(s.key, s)
}

func (o *NotificationObserver) delteSubscriber(key string) {
	o.subscribers.Delete(key)
}
