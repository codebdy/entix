package observer

import "sync"

type ModelObserver interface {
	Key() string
	ObjectCreated(object map[string]interface{})
	ObjectUpdated(object map[string]interface{})
	ObjectDeleted(object map[string]interface{})
}

var ModelObservers sync.Map

func AddObserver(obsr ModelObserver) {
	ModelObservers.Store(obsr.Key(), obsr)
}

func RemoveObserver(key string) {
	ModelObservers.Delete(key)
}
