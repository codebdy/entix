package imexport

import (
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
)

//处理要导入的实体对象，转化关联：
// a:xx=>a:{sync:xxx}
// 删掉关联的Id，保证所有数据都是新增
func convertInstanceValue(entity *graph.Entity, object map[string]interface{}) map[string]interface{} {
	object[consts.ID] = nil
	allAssociation := entity.Associations()
	for i := range allAssociation {
		asso := allAssociation[i]
		value := object[asso.Name()]

		if asso.IsCombination() {
			if value != nil {
				if asso.IsArray() {
					object[asso.Name()] = map[string]interface{}{
						consts.ARG_SYNC: convertInstanceValues(asso.TypeEntity(), value.([]map[string]interface{})),
					}
				} else {
					object[asso.Name()] = map[string]interface{}{
						consts.ARG_SYNC: convertInstanceValue(asso.TypeEntity(), value.(map[string]interface{})),
					}
				}

			} else {
				object[asso.Name()] = map[string]interface{}{
					consts.ARG_CLEAR: true,
				}
			}
		}
	}

	return object
}

func convertInstanceValues(entity *graph.Entity, objects []map[string]interface{}) []map[string]interface{} {
	for _, object := range objects {
		convertInstanceValue(entity, object)
	}
	return objects
}
