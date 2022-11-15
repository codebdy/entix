package observer

import "sync"

type ModelObserver interface {
	Key() string
	ObjectPosted(object map[string]interface{}, entityName string, userId, appId uint64)
	ObjectMultiPosted(objects []map[string]interface{}, entityName string, userId, appId uint64)
	ObjectDeleted(object map[string]interface{}, entityName string, userId, appId uint64)
	ObjectMultiDeleted(objects []map[string]interface{}, entityName string, userId, appId uint64)
}

var ModelObservers sync.Map

func AddObserver(obsr ModelObserver) {
	ModelObservers.Store(obsr.Key(), obsr)
}

func RemoveObserver(key string) {
	ModelObservers.Delete(key)
}

func EmitObjectPosted(object map[string]interface{}, entityName string, userId, appId uint64) {
	go func() {
		ModelObservers.Range(func(key interface{}, value interface{}) bool {
			value.(ModelObserver).ObjectPosted(object, entityName, userId, appId)
			return true
		})
	}()
}

func EmitObjectMultiPosted(objects []map[string]interface{}, entityName string, userId, appId uint64) {
	go func() {
		ModelObservers.Range(func(key interface{}, value interface{}) bool {
			value.(ModelObserver).ObjectMultiPosted(objects, entityName, userId, appId)
			return true
		})
	}()
}

func EmitObjectDeleted(object map[string]interface{}, entityName string, userId, appId uint64) {
	go func() {
		ModelObservers.Range(func(key interface{}, value interface{}) bool {
			value.(ModelObserver).ObjectDeleted(object, entityName, userId, appId)
			return true
		})
	}()
}

func EmitObjectMultiDeleted(objects []map[string]interface{}, entityName string, userId, appId uint64) {
	go func() {
		ModelObservers.Range(func(key interface{}, value interface{}) bool {
			value.(ModelObserver).ObjectMultiDeleted(objects, entityName, userId, appId)
			return true
		})
	}()
}
