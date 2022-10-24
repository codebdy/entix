package service

import "rxdrag.com/entify/model/data"

//把关联数据同步成输入参数的样子，只同步有组合关系的关联
//适用场景：恢复版本快照，导入应用
func SyncOne(instance *data.Instance) (interface{}, error) {
	return nil, nil
}
