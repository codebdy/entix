package orm

import (
	"log"

	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/model/data"
)

func (s *Session) SaveOne(instance *data.Instance) (interface{}, error) {
	if instance.IsInsert() {
		return s.InsertOne(instance)
	} else {
		return s.UpdateOne(instance)
	}
}

func (s *Session) InsertOne(instance *data.Instance) (interface{}, error) {
	id := s.InsertOneBody(instance)

	for _, asso := range instance.Associations {
		err := s.saveAssociation(asso, uint64(id))
		if err != nil {
			log.Println("Save reference failed:", err.Error())
			return nil, err
		}
	}

	savedObject := s.QueryOneEntityById(instance.Entity, id)

	return savedObject, nil
}

//只保存属性，不保存关联
func (s *Session) InsertOneBody(instance *data.Instance) int64 {
	sqlBuilder := dialect.GetSQLBuilder()
	saveStr := sqlBuilder.BuildInsertSQL(instance.Fields, instance.Table())
	values := makeFieldValues(instance.Fields)
	result, err := s.Dbx.Exec(saveStr, values...)
	if err != nil {
		log.Panic("Insert data failed:", err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Panic("LastInsertId failed:", err.Error())
	}
	return id
}

func (s *Session) UpdateOne(instance *data.Instance) (interface{}, error) {

	s.UpdateOneBody(instance)

	for _, ref := range instance.Associations {
		s.saveAssociation(ref, instance.Id)
	}

	savedObject := s.QueryOneEntityById(instance.Entity, instance.Id)

	return savedObject, nil
}

//只保存属性，不保存关联
func (s *Session) UpdateOneBody(instance *data.Instance) {
	sqlBuilder := dialect.GetSQLBuilder()

	saveStr := sqlBuilder.BuildUpdateSQL(instance.Id, instance.Fields, instance.Table())
	values := makeFieldValues(instance.Fields)
	log.Println(saveStr)
	_, err := s.Dbx.Exec(saveStr, values...)
	if err != nil {
		log.Panic("Update data failed:", err.Error())
	}
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

func (s *Session) saveAssociationInstance(ins *data.Instance) (interface{}, error) {
	targetData := InsanceData{consts.ID: ins.Id}

	saved, err := s.SaveOne(ins)
	if err != nil {
		return nil, err
	}
	targetData = saved.(InsanceData)

	return targetData, nil
}
func (s *Session) saveAssociation(r *data.AssociationRef, ownerId uint64) error {

	//这块逻辑还需要优化
	if r.Clear {
		s.clearAssociation(r, ownerId)
	}

	for _, ins := range r.Deleted {
		if r.Cascade() {
			s.DeleteInstance(ins)
		} else {
			povit := newAssociationPovit(r, ownerId, ins.Id)
			s.DeleteAssociationPovit(povit)
		}
	}

	for _, ins := range r.Added {
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

	for _, ins := range r.Updated {
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

	synced := r.Synced
	if len(synced) == 0 {
		return nil
	}

	s.clearAssociation(r, ownerId)

	for _, ins := range synced {
		targetId := ins.Id
		if !ins.IsEmperty() {
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
