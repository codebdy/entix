package data

import (
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/model/table"
)

//没有继承关系的关联
type Reference struct {
	Association *graph.Association
	Value       map[string]interface{}
}

func doConvertToInstances(data interface{}, isArray bool, entity *graph.Entity) []*Instance {
	instances := []*Instance{}
	if data == nil {
		return []*Instance{}
	}
	if isArray {
		objects := data.([]interface{})
		for i := range objects {
			instances = append(instances, NewInstance(objects[i].(map[string]interface{}), entity))
		}
	} else {
		instances = append(instances, NewInstance(data.(map[string]interface{}), entity))
	}

	return instances
}

func (r *Reference) convertToInstances(data interface{}) []*Instance {
	return doConvertToInstances(data, r.Association.IsArray(), r.TypeEntity())
}

func (r *Reference) Deleted() []*Instance {
	return r.convertToInstances(r.Value[consts.ARG_DELETE])
}

func (r *Reference) Added() []*Instance {
	return r.convertToInstances(r.Value[consts.ARG_ADD])
}

func (r *Reference) Updated() []*Instance {
	return r.convertToInstances(r.Value[consts.ARG_UPDATE])
}

func (r *Reference) Synced() []*Instance {
	return r.convertToInstances(r.Value[consts.ARG_SYNC])
}

func (r *Reference) Cascade() bool {
	if r.Value[consts.ARG_CASCADE] != nil {
		return r.Value[consts.ARG_CASCADE].(bool)
	}
	return false
}

func (r *Reference) SourceColumn() *table.Column {
	for i := range r.Association.Relation.Table.Columns {
		column := r.Association.Relation.Table.Columns[i]
		if column.Name == r.Association.Relation.SourceEntity.TableName() {
			return column
		}
	}
	return nil
}

func (r *Reference) TargetColumn() *table.Column {
	for i := range r.Association.Relation.Table.Columns {
		column := r.Association.Relation.Table.Columns[i]
		if column.Name == r.Association.Relation.TargetEntity.TableName() {
			return column
		}
	}
	return nil
}

func (r *Reference) Table() *table.Table {
	return r.Association.Relation.Table
}

func (r *Reference) IsSource() bool {
	return r.Association.IsSource()
}

func (r *Reference) OwnerColumn() *table.Column {
	if r.IsSource() {
		return r.SourceColumn()
	} else {
		return r.TargetColumn()
	}
}
func (r *Reference) TypeColumn() *table.Column {
	if !r.IsSource() {
		return r.SourceColumn()
	} else {
		return r.TargetColumn()
	}
}

func (r *Reference) TypeEntity() *graph.Entity {
	entity := r.Association.TypeEntity()
	if entity != nil {
		return entity
	}

	panic("Can not find reference entity")
}

func (r *Reference) IsCombination() bool {
	return r.IsSource() &&
		(r.Association.Relation.RelationType == meta.TWO_WAY_COMBINATION ||
			r.Association.Relation.RelationType == meta.ONE_WAY_COMBINATION)
}
