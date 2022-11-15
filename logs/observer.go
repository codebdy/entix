package logs

import (
	"context"

	"rxdrag.com/entify/model/observer"
)

type ModelObserver struct {
	key string
}

func init() {
	//创建模型监听器
	modelObserver := &ModelObserver{
		key: "ModelObserverForLogs",
	}
	observer.AddObserver(modelObserver)
}

func (o *ModelObserver) Key() string {
	return o.key
}

func (o *ModelObserver) ObjectPosted(object map[string]interface{}, entityName string, ctx context.Context) {

}

func (o *ModelObserver) ObjectMultiPosted(objects []map[string]interface{}, entityName string, ctx context.Context) {

}
func (o *ModelObserver) ObjectDeleted(object map[string]interface{}, entityName string, ctx context.Context) {

}

func (o *ModelObserver) ObjectMultiDeleted(objects []map[string]interface{}, entityName string, ctx context.Context) {

}
