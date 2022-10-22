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
		log.Panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		log.Panic(err)
	}

	instance.Inserted(uint64(id))
}

func (s *Session) update(instance *data.Instance) (interface{}, error) {

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

	for _, ref := range instance.Associations {
		s.doSaveAssociation(ref, instance.Id)
	}

	savedObject := s.QueryOneEntityById(instance.Entity, instance.Id)

	return savedObject, nil
}

func newAssociationPovit(r *data.Reference, ownerId uint64, tarId uint64) *data.AssociationPovit {
	sourceId := ownerId
	targetId := tarId

	if !r.IsSource() {
		sourceId = targetId
		targetId = ownerId
	}

	return data.NewAssociationPovit(r, sourceId, targetId)

}

func (s *Session) saveAssociationInstance(ins *data.Instance) (interface{}, error) {
	targetData := InsanceData{consts.ID: ins.Id}

	saved, err := s.SaveOne(ins)
	if err != nil {
		return nil, err
	}
	targetData = saved.(InsanceData)

	return targetData, nil
}

func (s *Session) doSaveAssociation(r *data.AssociationRef, ownerId uint64) error {

	for _, ins := range r.Deleted() {
		if r.Cascade() {
			s.DeleteInstance(ins)
		} else {
			povit := newAssociationPovit(r, ownerId, ins.Id)
			s.DeleteAssociationPovit(povit)
		}
	}

	for _, ins := range r.Added() {
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
