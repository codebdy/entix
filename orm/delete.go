package orm

import (
	"fmt"

	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/model/data"
)

type InsanceData = map[string]interface{}

func (con *Session) clearAssociation(r *data.Reference, ownerId uint64) {
	con.deleteAssociationPovit(r, ownerId)

	if r.IsCombination() {
		con.deleteAssociatedInstances(r, ownerId)
	}
}

func (con *Session) deleteAssociationPovit(r *data.Reference, ownerId uint64) {
	sqlBuilder := dialect.GetSQLBuilder()
	sql := sqlBuilder.BuildClearAssociationSQL(ownerId, r.Table().Name, r.OwnerColumn().Name)
	_, err := con.Dbx.Exec(sql)
	fmt.Println("deleteAssociationPovit SQL:" + sql)
	if err != nil {
		panic(err.Error())
	}
}

func (con *Session) deleteAssociatedInstances(r *data.Reference, ownerId uint64) {
	typeEntity := r.TypeEntity()
	associatedInstances := con.QueryAssociatedInstances(r, ownerId)
	for i := range associatedInstances {
		ins := data.NewInstance(associatedInstances[i], typeEntity)
		con.DeleteInstance(ins)
	}
}

func (con *Session) DeleteAssociationPovit(povit *data.AssociationPovit) {
	sqlBuilder := dialect.GetSQLBuilder()
	sql := sqlBuilder.BuildDeletePovitSQL(povit)
	_, err := con.Dbx.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
}

func (con *Session) DeleteInstance(instance *data.Instance) {
	var sql string
	sqlBuilder := dialect.GetSQLBuilder()
	tableName := instance.Table().Name
	if instance.Entity.IsSoftDelete() {
		sql = sqlBuilder.BuildSoftDeleteSQL(instance.Id, tableName)
	} else {
		sql = sqlBuilder.BuildDeleteSQL(instance.Id, tableName)
	}

	_, err := con.Dbx.Exec(sql)
	if err != nil {
		panic(err.Error())
	}

	associstions := instance.Associations
	for i := range associstions {
		asso := associstions[i]
		if asso.IsCombination() {
			if !asso.TypeEntity().IsSoftDelete() {
				con.deleteAssociationPovit(asso, instance.Id)
			}
			con.deleteAssociatedInstances(asso, instance.Id)
		}
	}
}
