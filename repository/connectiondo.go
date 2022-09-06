package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"rxdrag.com/entify/config"
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/db"
	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
)

type InsanceData = map[string]interface{}

func (con *Connection) buildQueryInterfaceSQL(intf *graph.Interface, args map[string]interface{}) (string, []interface{}) {
	var (
		sqls       []string
		paramsList []interface{}
	)
	builder := dialect.GetSQLBuilder()
	for i := range intf.Children {
		entity := intf.Children[i]
		whereArgs := con.v.WeaveAuthInArgs(entity.Uuid(), args[consts.ARG_WHERE])
		argEntity := graph.BuildArgEntity(
			entity,
			whereArgs,
			con,
		)
		queryStr := builder.BuildQuerySQLBody(argEntity, intf.AllAttributes())
		if where, ok := whereArgs.(graph.QueryArg); ok {
			whereSQL, params := builder.BuildWhereSQL(argEntity, intf.AllAttributes(), where)
			if whereSQL != "" {
				queryStr = queryStr + " WHERE " + whereSQL
			}

			paramsList = append(paramsList, params...)
		}
		queryStr = queryStr + builder.BuildOrderBySQL(argEntity, args[consts.ARG_ORDERBY])

		sqls = append(sqls, queryStr)
	}

	return strings.Join(sqls, " UNION "), paramsList
}

func (con *Connection) buildQueryEntitySQL(
	entity *graph.Entity,
	args map[string]interface{},
	whereArgs interface{},
	argEntity *graph.ArgEntity,
	queryBody string,
) (string, []interface{}) {
	var paramsList []interface{}
	//whereArgs := con.v.WeaveAuthInArgs(entity.Uuid(), args[consts.ARG_WHERE])
	// argEntity := graph.BuildArgEntity(
	// 	entity,
	// 	whereArgs,
	// 	con,
	// )
	builder := dialect.GetSQLBuilder()
	queryStr := queryBody
	if where, ok := whereArgs.(graph.QueryArg); ok {
		whereSQL, params := builder.BuildWhereSQL(argEntity, entity.AllAttributes(), where)
		if whereSQL != "" {
			queryStr = queryStr + " WHERE " + whereSQL
		}
		paramsList = append(paramsList, params...)
	}

	queryStr = queryStr + builder.BuildOrderBySQL(argEntity, args[consts.ARG_ORDERBY])
	return queryStr, paramsList
}

func (con *Connection) buildQueryEntityRecordsSQL(entity *graph.Entity, args map[string]interface{}) (string, []interface{}) {
	whereArgs := con.v.WeaveAuthInArgs(entity.Uuid(), args[consts.ARG_WHERE])
	argEntity := graph.BuildArgEntity(
		entity,
		whereArgs,
		con,
	)
	builder := dialect.GetSQLBuilder()
	queryStr := builder.BuildQuerySQLBody(argEntity, entity.AllAttributes())
	sqlStr, params := con.buildQueryEntitySQL(entity, args, whereArgs, argEntity, queryStr)

	if args[consts.ARG_LIMIT] != nil {
		sqlStr = sqlStr + fmt.Sprintf(" LIMIT %d ", args[consts.ARG_LIMIT])
	}
	if args[consts.ARG_OFFSET] != nil {
		sqlStr = sqlStr + fmt.Sprintf(" OFFSET %d ", args[consts.ARG_OFFSET])
	}

	return sqlStr, params
}

func (con *Connection) buildQueryEntityCountSQL(entity *graph.Entity, args map[string]interface{}) (string, []interface{}) {
	whereArgs := con.v.WeaveAuthInArgs(entity.Uuid(), args[consts.ARG_WHERE])
	argEntity := graph.BuildArgEntity(
		entity,
		whereArgs,
		con,
	)
	builder := dialect.GetSQLBuilder()
	queryStr := builder.BuildQueryCountSQLBody(argEntity)
	return con.buildQueryEntitySQL(
		entity,
		args,
		whereArgs,
		argEntity,
		queryStr,
	)
}

func (con *Connection) doQueryInterface(intf *graph.Interface, args map[string]interface{}) map[string]interface{} {
	sql, params := con.buildQueryInterfaceSQL(intf, args)

	rows, err := con.Dbx.Query(sql, params...)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	var instances []InsanceData
	for rows.Next() {
		values := makeInterfaceQueryValues(intf)
		err = rows.Scan(values...)
		if err != nil {
			panic(err.Error())
		}
		instances = append(instances, convertValuesToInterface(values, intf))
	}

	instancesIds := make([]interface{}, len(instances))
	for i := range instances {
		instancesIds[i] = instances[i][consts.ID]
	}

	for i := range intf.Children {
		child := intf.Children[i]
		oneEntityInstances := con.doQueryByIds(child, instancesIds)
		merageInstances(instances, oneEntityInstances)
	}

	return map[string]interface{}{
		consts.NODES: instances,
		consts.TOTAL: 0,
	}
}

func (con *Connection) doQueryEntity(entity *graph.Entity, args map[string]interface{}) map[string]interface{} {
	if !con.v.CanReadEntity(entity.Uuid()) {
		panic(consts.NO_PERMISSION)
	}
	sqlStr, params := con.buildQueryEntityRecordsSQL(entity, args)
	fmt.Println("doQueryEntity SQL:", sqlStr, params)
	rows, err := con.Dbx.Query(sqlStr, params...)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	var instances []InsanceData
	for rows.Next() {
		values := makeEntityQueryValues(entity)
		err = rows.Scan(values...)
		if err != nil {
			panic(err.Error())
		}
		instances = append(instances, convertValuesToEntity(values, entity))
	}

	sqlStr, params = con.buildQueryEntityCountSQL(entity, args)
	fmt.Println("doQueryEntity count SQL:", sqlStr, params)
	count := 0
	err = con.Dbx.QueryRow(sqlStr, params...).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		count = 0
	case err != nil:
		panic(err.Error())
	}

	defer rows.Close()

	return map[string]interface{}{
		consts.NODES: instances,
		consts.TOTAL: count,
	}
}

func (con *Connection) QueryOneEntityById(entity *graph.Entity, id interface{}) interface{} {
	return con.doQueryOneEntity(entity, graph.QueryArg{
		consts.ARG_WHERE: graph.QueryArg{
			consts.ID: graph.QueryArg{
				consts.ARG_EQ: id,
			},
		},
	})
}
func (con *Connection) doQueryOneInterface(intf *graph.Interface, args map[string]interface{}) interface{} {
	querySql, params := con.buildQueryInterfaceSQL(intf, args)

	values := makeInterfaceQueryValues(intf)
	err := con.Dbx.QueryRow(querySql, params...).Scan(values...)

	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err.Error())
	}

	instance := convertValuesToInterface(values, intf)
	for i := range intf.Children {
		child := intf.Children[i]
		oneEntityInstances := con.doQueryByIds(child, []interface{}{instance[consts.ID]})
		if len(oneEntityInstances) > 0 {
			return oneEntityInstances[0]
		}
	}
	return nil
}

func (con *Connection) doQueryOneEntity(entity *graph.Entity, args map[string]interface{}) interface{} {
	queryStr, params := con.buildQueryEntityRecordsSQL(entity, args)

	values := makeEntityQueryValues(entity)
	fmt.Println("doQueryOneEntity SQL:", queryStr)
	err := con.Dbx.QueryRow(queryStr, params...).Scan(values...)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err.Error())
	}

	return convertValuesToEntity(values, entity)
}

func (con *Connection) doInsertOne(instance *data.Instance) (interface{}, error) {
	sqlBuilder := dialect.GetSQLBuilder()
	saveStr := sqlBuilder.BuildInsertSQL(instance.Fields, instance.Table())
	values := makeSaveValues(instance.Fields)
	result, err := con.Dbx.Exec(saveStr, values...)
	if err != nil {
		fmt.Println("Insert data failed:", err.Error())
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("LastInsertId failed:", err.Error())
		return nil, err
	}
	for _, asso := range instance.Associations {
		err = con.doSaveAssociation(asso, uint64(id))
		if err != nil {
			fmt.Println("Save reference failed:", err.Error())
			return nil, err
		}
	}

	savedObject := con.QueryOneEntityById(instance.Entity, id)

	//affectedRows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("RowsAffected failed:", err.Error())
		return nil, err
	}

	return savedObject, nil
}

func (con *Connection) doQueryAssociatedInstances(r *data.Reference, ownerId uint64) []InsanceData {
	var instances []InsanceData
	builder := dialect.GetSQLBuilder()
	entity := r.TypeEntity()
	queryStr := builder.BuildQueryAssociatedInstancesSQL(entity, ownerId, r.Table().Name, r.OwnerColumn().Name, r.TypeColumn().Name)
	rows, err := con.Dbx.Query(queryStr)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		values := makeEntityQueryValues(entity)
		err = rows.Scan(values...)
		if err != nil {
			panic(err.Error())
		}
		instances = append(instances, convertValuesToEntity(values, entity))
	}

	return instances
}

func (con *Connection) doQueryByIds(entity *graph.Entity, ids []interface{}) []InsanceData {
	var instances []map[string]interface{}
	builder := dialect.GetSQLBuilder()
	sql := builder.BuildQueryByIdsSQL(entity, len(ids))
	rows, err := con.Dbx.Query(sql, ids...)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		values := makeEntityQueryValues(entity)
		err = rows.Scan(values...)
		if err != nil {
			panic(err.Error())
		}
		instances = append(instances, convertValuesToEntity(values, entity))
	}

	return instances
}

func (con *Connection) doBatchRealAssociations(
	association *graph.Association,
	ids []uint64,
	args graph.QueryArg,
	v *AbilityVerifier,
) []InsanceData {
	var instances []map[string]interface{}
	var paramsList []interface{}

	builder := dialect.GetSQLBuilder()
	typeEntity := association.TypeEntity()
	whereArgs := con.v.WeaveAuthInArgs(typeEntity.Uuid(), args[consts.ARG_WHERE])
	argEntity := graph.BuildArgEntity(
		typeEntity,
		whereArgs,
		con,
	)

	queryStr := builder.BuildBatchAssociationBodySQL(argEntity,
		typeEntity.AllAttributes(),
		association.Relation.Table.Name,
		association.Owner().TableName(),
		association.TypeEntity().TableName(),
		ids,
	)

	if where, ok := whereArgs.(graph.QueryArg); ok {
		whereSQL, params := builder.BuildWhereSQL(argEntity, typeEntity.AllAttributes(), where)
		if whereSQL != "" {
			queryStr = queryStr + " AND " + whereSQL
		}
		paramsList = append(paramsList, params...)
	}

	queryStr = queryStr + builder.BuildOrderBySQL(argEntity, args[consts.ARG_ORDERBY])
	fmt.Println("doBatchRealAssociations SQL:	", queryStr)
	rows, err := con.Dbx.Query(queryStr, paramsList...)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		values := makeEntityQueryValues(typeEntity)
		var idValue db.NullUint64
		values = append(values, &idValue)
		err = rows.Scan(values...)
		if err != nil {
			panic(err.Error())
		}
		instance := convertValuesToEntity(values, typeEntity)
		instance[consts.ASSOCIATION_OWNER_ID] = values[len(values)-1].(*db.NullUint64).Uint64
		instances = append(instances, instance)
	}

	return instances
}

func merageInstances(source []InsanceData, target []InsanceData) {
	for i := range source {
		souceObj := source[i]
		for j := range target {
			targetObj := target[j]
			if souceObj[consts.ID] == targetObj[consts.ID] {
				targetObj[consts.ASSOCIATION_OWNER_ID] = souceObj[consts.ASSOCIATION_OWNER_ID]
				source[i] = targetObj
			}
		}
	}
}

func (con *Connection) doUpdateOne(instance *data.Instance) (interface{}, error) {

	sqlBuilder := dialect.GetSQLBuilder()

	saveStr := sqlBuilder.BuildUpdateSQL(instance.Id, instance.Fields, instance.Table())
	values := makeSaveValues(instance.Fields)
	fmt.Println(saveStr)
	_, err := con.Dbx.Exec(saveStr, values...)
	if err != nil {
		fmt.Println("Update data failed:", err.Error())
		return nil, err
	}

	for _, ref := range instance.Associations {
		con.doSaveAssociation(ref, instance.Id)
	}

	savedObject := con.QueryOneEntityById(instance.Entity, instance.Id)

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

func (con *Connection) saveAssociationInstance(ins *data.Instance) interface{} {
	targetData := InsanceData{consts.ID: ins.Id}

	saved, err := con.doSaveOne(ins)
	if err != nil {
		return err
	}
	targetData = saved.(InsanceData)

	return targetData
}

func (con *Connection) doSaveAssociation(r *data.Reference, ownerId uint64) error {

	for _, ins := range r.Deleted() {
		if r.Cascade() {
			con.doDeleteInstance(ins)
		} else {
			povit := newAssociationPovit(r, ownerId, ins.Id)
			con.doDeleteAssociationPovit(povit)
		}
	}

	for _, ins := range r.Added() {
		targetData := con.saveAssociationInstance(ins)

		if savedIns, ok := targetData.(InsanceData); ok {
			tarId := savedIns[consts.ID].(uint64)
			relationInstance := newAssociationPovit(r, ownerId, tarId)
			con.doSaveAssociationPovit(relationInstance)
		} else {
			panic("Save Association error")
		}

	}

	for _, ins := range r.Updated() {
		// if ins.Id == 0 {
		// 	panic("Can not add new instance when update")
		// }
		targetData := con.saveAssociationInstance(ins)
		if savedIns, ok := targetData.(InsanceData); ok {
			tarId := savedIns[consts.ID].(uint64)
			relationInstance := newAssociationPovit(r, ownerId, tarId)

			con.doSaveAssociationPovit(relationInstance)
		} else {
			panic("Save Association error")
		}
	}

	synced := r.Synced()
	if len(synced) == 0 {
		return nil
	}

	con.clearAssociation(r, ownerId)

	for _, ins := range synced {
		// if ins.Id == 0 {
		// 	panic("Can not add new instance when update")
		// }
		targetData := con.saveAssociationInstance(ins)

		if savedIns, ok := targetData.(InsanceData); ok {
			tarId := savedIns[consts.ID].(uint64)
			relationInstance := newAssociationPovit(r, ownerId, tarId)

			con.doSaveAssociationPovit(relationInstance)
		} else {
			panic("Save Association error")
		}
	}

	return nil
}

func (con *Connection) clearAssociation(r *data.Reference, ownerId uint64) {
	con.deleteAssociationPovit(r, ownerId)

	if r.IsCombination() {
		con.deleteAssociatedInstances(r, ownerId)
	}
}

func (con *Connection) deleteAssociationPovit(r *data.Reference, ownerId uint64) {
	sqlBuilder := dialect.GetSQLBuilder()
	sql := sqlBuilder.BuildClearAssociationSQL(ownerId, r.Table().Name, r.OwnerColumn().Name)
	_, err := con.Dbx.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
}

func (con *Connection) deleteAssociatedInstances(r *data.Reference, ownerId uint64) {
	typeEntity := r.TypeEntity()
	associatedInstances := con.doQueryAssociatedInstances(r, ownerId)
	for i := range associatedInstances {
		ins := data.NewInstance(associatedInstances[i], typeEntity)
		con.doDeleteInstance(ins)
	}
}

func (con *Connection) doSaveAssociationPovit(povit *data.AssociationPovit) {
	sqlBuilder := dialect.GetSQLBuilder()
	sql := sqlBuilder.BuildQueryPovitSQL(povit)
	rows, err := con.Dbx.Query(sql)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}
	if !rows.Next() {
		sql = sqlBuilder.BuildInsertPovitSQL(povit)
		_, err := con.Dbx.Exec(sql)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (con *Connection) doDeleteAssociationPovit(povit *data.AssociationPovit) {
	sqlBuilder := dialect.GetSQLBuilder()
	sql := sqlBuilder.BuildDeletePovitSQL(povit)
	_, err := con.Dbx.Exec(sql)
	if err != nil {
		panic(err.Error())
	}
}

func (con *Connection) doSaveOne(instance *data.Instance) (interface{}, error) {
	if instance.IsInsert() {
		return con.doInsertOne(instance)
	} else {
		return con.doUpdateOne(instance)
	}
}

func (con *Connection) doDeleteInstance(instance *data.Instance) {
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

func (con *Connection) doCheckEntity(name string) bool {
	sqlBuilder := dialect.GetSQLBuilder()
	var count int
	err := con.Dbx.QueryRow(sqlBuilder.BuildTableCheckSQL(name, config.GetDbConfig().Database)).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		panic(err.Error())
	}
	return count > 0
}
