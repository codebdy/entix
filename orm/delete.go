package orm

import (
	"log"

	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/model/data"
)

type InsanceData = map[string]interface{}

func (con *Session) clearAssociation(r *data.AssociationRef, ownerId uint64) {
	con.deleteAssociationPovit(r, ownerId)

	if r.Association.IsCombination() {
		con.deleteAssociatedInstances(r, ownerId)
	}
}

func (s *Session) deleteAssociationPovit(r *data.AssociationRef, ownerId uint64) {
	sqlBuilder := dialect.GetSQLBuilder()
	//先检查是否有数据，如果有再删除，避免死锁
	sql := sqlBuilder.BuildCheckAssociationSQL(ownerId, r.Table().Name, r.OwnerColumn().Name)
	rows, err := s.Dbx.Query(sql)
	if err != nil {
		log.Println(err.Error())
		return
	} else {
		var count int64
		for rows.Next() {
			rows.Scan(&count)
		}
		if count > 0 {
			sql = sqlBuilder.BuildClearAssociationSQL(ownerId, r.Table().Name, r.OwnerColumn().Name)
			_, err := s.Dbx.Exec(sql)
			log.Println("deleteAssociationPovit SQL:" + sql)
			if err != nil {
				panic(err.Error())
			}
		}
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
