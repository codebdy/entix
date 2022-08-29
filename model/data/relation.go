package data

import (
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/model/table"
)

//没有继承关系的关联
type Reference struct {
	Association *graph.Association
	Value       []interface{}
}

func doConvertToInstances(data interface{}, entity *graph.Entity) []*Instance {
	instances := []*Instance{}
	if data == nil {
		return nil
	}

	objects := data.([]interface{})
	for i := range objects {
		instances = append(instances, NewInstance(objects[i].(map[string]interface{}), entity))
	}

	return instances
}

func (r *Reference) convertToInstances(data interface{}) []*Instance {
	return doConvertToInstances(data, r.TypeEntity())
}

func (r *Reference) Associated() []*Instance {
	return r.convertToInstances(r.Value)
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
