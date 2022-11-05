package orm

import (
	"log"

	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/model/data"
)

type InsanceData = map[string]interface{}

func (con *Session) clearSyncedAssociation(r *data.AssociationRef, ownerId uint64, synced []*data.Instance) {
	//不能暴力删除，这个地方需要这么处理：
	// 1、统计sync的ids
	// 2、清除跟这些ids不重合的数据（先检查是否存在，再删除）
	syncedIds := []uint64{}
	for _, ins := range synced {
		if ins.Id != 0 {
			syncedIds = append(syncedIds, ins.Id)
		}
	}

	//instancesToDeleted :=

	if r.Association.IsCombination() {

		con.deleteAssociatedInstances(r, ownerId)
	}
	con.deleteAssociationPovit(r, ownerId)
}
func (con *Session) clearAssociation(r *data.AssociationRef, ownerId uint64) {
	if r.Association.IsCombination() {
		con.deleteAssociatedInstances(r, ownerId)
	}
	con.deleteAssociationPovit(r, ownerId)
}

func (s *Session) deleteAssociationPovit(r *data.AssociationRef, ownerId uint64) {
	sqlBuilder := dialect.GetSQLBuilder()
	//先检查是否有数据，如果有再删除，避免死锁
	sql := sqlBuilder.BuildCheckAssociationSQL(ownerId, r.Table().Name, r.OwnerColumn().Name)
	count := s.queryCount(sql)
	if count > 0 {
		sql = sqlBuilder.BuildClearAssociationSQL(ownerId, r.Table().Name, r.OwnerColumn().Name)
		_, err := s.Dbx.Exec(sql)
		log.Println("deleteAssociationPovit SQL:" + sql)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (s *Session) queryCount(countSql string) int64 {
	rows, err := s.Dbx.Query(countSql)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return 0
	} else {
		var count int64
		for rows.Next() {
			rows.Scan(&count)
		}
		return count
	}
}

func (s *Session) deleteAssociatedInstances(r *data.AssociationRef, ownerId uint64) {
	typeEntity := r.TypeEntity()
	associatedInstances := s.QueryAssociatedInstances(r, ownerId)
	for i := range associatedInstances {
		ins := data.NewInstance(associatedInstances[i], typeEntity)
		s.DeleteInstance(ins)
	}
}

func (s *Session) DeleteAssociationPovit(povit *data.AssociationPovit) {
	sqlBuilder := dialect.GetSQLBuilder()
	sql := sqlBuilder.BuildDeletePovitSQL(povit)
	_, err := s.Dbx.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
}

func (s *Session) DeleteInstance(instance *data.Instance) {
	var sql string
	sqlBuilder := dialect.GetSQLBuilder()
	tableName := instance.Table().Name
	if instance.Entity.IsSoftDelete() {
		sql = sqlBuilder.BuildSoftDeleteSQL(instance.Id, tableName)
	} else {
		sql = sqlBuilder.BuildDeleteSQL(instance.Id, tableName)
	}
	_, err := s.Dbx.Exec(sql)
	if err != nil {
		panic(err.Error())
	}

	associstions := instance.Associations
	for i := range associstions {
		asso := associstions[i]
		if asso.Association.IsCombination() {
			if !asso.TypeEntity().IsSoftDelete() {
				s.deleteAssociationPovit(asso, instance.Id)
			}
			s.deleteAssociatedInstances(asso, instance.Id)
		}
	}
}
