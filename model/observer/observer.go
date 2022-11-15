package observer

import "sync"

type ModelObserver interface {
	Key() string
	ObjectCreated(object map[string]interface{}, entityName string)
	ObjectUpdated(object map[string]interface{}, entityName string)
	ObjectDeleted(object map[string]interface{}, entityName string)
}

var ModelObservers sync.Map

func AddObserver(obsr ModelObserver) {
	ModelObservers.Store(obsr.Key(), obsr)
}

func RemoveObserver(key string) {
	ModelObservers.Delete(key)
}
