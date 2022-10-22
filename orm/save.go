package orm

import (
	"fmt"
	"log"

	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/model/data"
)

// 两阶段操作，先插入、再更新
func (s *Session) SaveOne(instance *data.Instance) (interface{}, error) {
	//第一阶段插入
	s.preInsertAll(instance)
	//第二阶段更新
	s.update(instance)

	savedObject := s.QueryOneEntityById(instance.Entity, instance.Id)
	return savedObject, nil
}

func (s *Session) preInsertAll(instance *data.Instance) {
	if instance.IsInsert() {
		s.insert(instance)
	}

	for i := range instance.Associations {
		assoc := instance.Associations[i]
		for j := range assoc.Added {
			s.preInsertAll(assoc.Added[j])
		}
		for j := range assoc.Updated {
			s.preInsertAll(assoc.Updated[j])
		}
		for j := range assoc.Synced {
			s.preInsertAll(assoc.Synced[j])
		}
	}
}

func (s *Session) insert(instance *data.Instance) {
	sqlBuilder := dialect.GetSQLBuilder()
	saveStr := sqlBuilder.BuildInsertSQL(instance.Fields, instance.Table())
	values := makeFieldValues(instance.Fields)
	result, err := s.Dbx.Exec(saveStr, values...)
	if err != nil {
		log.Println(err.Error())
		log.Panic(err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		log.Panic(err.Error())
	}

	instance.Inserted(uint64(id))
}

func (s *Session) update(instance *data.Instance) (interface{}, error) {
	oldInstanceData := s.QueryOneEntityById(instance.Entity, instance.Id)
	sqlBuilder := dialect.GetSQLBuilder()
	//在本方存的关联
	columnAssocs := instance.ColumnAssociations()
	saveStr := sqlBuilder.BuildUpdateSQL(instance.Id, instance.Fields, columnAssocs, instance.Table())
	values := makeFieldValues(instance.Fields)
	values = append(values, makeAssociationValues(columnAssocs))
	fmt.Println(saveStr)
	_, err := s.Dbx.Exec(saveStr, values...)
	if err != nil {
		log.Panic("Update data failed:", err.Error())
	}

	//本方关联类进一步处理其下级关联
	for _, assocRef := range instance.ColumnAssociations() {
		//如果级联，删除旧数据
		if assocRef.Cascade() && oldInstanceData != nil {
			//旧ID
			oldId := oldInstanceData.(map[string]interface{})[assocRef.Association.Name()]
			//新ID
			newId := assocRef.AssociatedId()

			if newId != nil && oldId != newId {
				s.DeleteInstance(data.NewInstance(map[string]interface{}{"id": oldId}, assocRef.TypeEntity()))
			}
		}

		for _, toAdd := range assocRef.Added {
			if toAdd.IsEmperty() {
				s.update(toAdd)
			}

		}

		for _, toSync := range assocRef.Synced {
			s.update(toSync)
		}

		for _, toUpdated := range assocRef.Updated {
			s.update(toUpdated)
		}
	}

	//对方存的关联
	for _, assocRef := range instance.TargetColumnAssociations() {
		s.saveTargetAssociation(assocRef, instance.Id)
	}

	//中间表存的关联
	for _, assocRef := range instance.PovitAssociations() {
		s.savePovitAssociation(assocRef, instance.Id)
	}

	savedObject := s.QueryOneEntityById(instance.Entity, instance.Id)

	return savedObject, nil
}

func newAssociationPovit(r *data.AssociationRef, ownerId uint64, tarId uint64) *data.AssociationPovit {
	sourceId := ownerId
	targetId := tarId

	if !r.IsSource() {
		sourceId = targetId
		targetId = ownerId
	}

	return data.NewAssociationPovit(r, sourceId, targetId)

}

func (s *Session) saveAssociationInstance(ins *data.Instance) interface{} {
	targetData := InsanceData{consts.ID: ins.Id}

	saved, err := s.SaveOne(ins)
	if err != nil {
		log.Panic(err.Error())
	}
	targetData = saved.(InsanceData)

	return targetData
}

func (s *Session) saveTargetAssociation(r *data.AssociationRef, ownerId uint64) {
}

func (s *Session) savePovitAssociation(r *data.AssociationRef, ownerId uint64) {

	for _, ins := range r.Deleted() {
		if r.Cascade() {
			s.DeleteInstance(ins)
		} else {
			povit := newAssociationPovit(r, ownerId, ins.Id)
			s.DeleteAssociationPovit(povit)
		}
	}

	for _, ins := range r.Added() {
		targetData := s.saveAssociationInstance(ins)

		if savedIns, ok := targetData.(InsanceData); ok {
			tarId := savedIns[consts.ID].(uint64)
			relationInstance := newAssociationPovit(r, ownerId, tarId)
			s.SaveAssociationPovit(relationInstance)
		} else {
			panic("Save Association error")
		}

	}

	for _, ins := range r.Updated() {
		// if ins.Id == 0 {
		// 	panic("Can not add new instance when update")
		// }
		targetData, err := s.saveAssociationInstance(ins)
		if err != nil {
			panic("Save Association error:" + err.Error())
		} else {
			if savedIns, ok := targetData.(InsanceData); ok {
				tarId := savedIns[consts.ID].(uint64)
				relationInstance := newAssociationPovit(r, ownerId, tarId)

				s.SaveAssociationPovit(relationInstance)
			} else {
				panic("Save Association error")
			}
		}
	}

	synced := r.Synced()
	if len(synced) == 0 {
		return nil
	}

	//有死锁bug，暂时不解决
	s.clearAssociation(r, ownerId)

	for _, ins := range synced {
		targetId := ins.Id
		if !ins.IsEmperty {
			targetData, err := s.saveAssociationInstance(ins)
			if err != nil {
				panic("Save Association error:" + err.Error())
			} else {
				if savedIns, ok := targetData.(InsanceData); ok {
					targetId = savedIns[consts.ID].(uint64)
				} else {
					panic("Save Association error")
				}
			}
		}
		relationInstance := newAssociationPovit(r, ownerId, targetId)
		s.SaveAssociationPovit(relationInstance)
	}

	return nil
}

func (s *Session) SaveAssociationPovit(povit *data.AssociationPovit) {
	sqlBuilder := dialect.GetSQLBuilder()
	sql := sqlBuilder.BuildQueryPovitSQL(povit)
	rows, err := s.Dbx.Query(sql)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	if !rows.Next() {
		sql = sqlBuilder.BuildInsertPovitSQL(povit)
		_, err := s.Dbx.Exec(sql)
		if err != nil {
			panic(err.Error())
		}
	}
}
