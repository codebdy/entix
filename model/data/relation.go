package data

import (
	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/model/table"
)

type AssociationRef struct {
	Association *graph.Association
	Added       []*Instance
	Deleted     []*Instance
	Updated     []*Instance
	Synced      []*Instance
	Cascade     bool
}

func NewAssociation(value map[string]interface{}, assoc *graph.Association) *AssociationRef {
	AssociationRef := AssociationRef{
		Association: assoc,
	}

	AssociationRef.init(value)
	return &AssociationRef
}

func (r *AssociationRef) init(value map[string]interface{}) {
	r.Deleted = r.convertToInstances(value[consts.ARG_DELETE])
	r.Added = r.convertToInstances(value[consts.ARG_ADD])
	r.Updated = r.convertToInstances(value[consts.ARG_UPDATE])
	r.Synced = r.convertToInstances(value[consts.ARG_SYNC])
	if value[consts.ARG_CASCADE] != nil {
		r.Cascade = value[consts.ARG_CASCADE].(bool)
	}
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

func (r *AssociationRef) convertToInstances(data interface{}) []*Instance {
	return doConvertToInstances(data, r.Association.IsArray(), r.TypeEntity())
}

func (r *AssociationRef) SourceColumn() *table.Column {
	for i := range r.Association.Relation.Table.Columns {
		column := r.Association.Relation.Table.Columns[i]
		if column.Name == r.Association.Relation.SourceEntity.TableName() {
			return column
		}
	}
	return nil
}

func (r *AssociationRef) TargetColumn() *table.Column {
	for i := range r.Association.Relation.Table.Columns {
		column := r.Association.Relation.Table.Columns[i]
		if column.Name == r.Association.Relation.TargetEntity.TableName() {
			return column
		}
	}
	return nil
}

func (r *AssociationRef) Table() *table.Table {
	return r.Association.Relation.Table
}

func (r *AssociationRef) IsSource() bool {
	return r.Association.IsSource()
}

func (r *AssociationRef) OwnerColumn() *table.Column {
	if r.IsSource() {
		return r.SourceColumn()
	} else {
		return r.TargetColumn()
	}
}
func (r *AssociationRef) TypeColumn() *table.Column {
	if !r.IsSource() {
		return r.SourceColumn()
	} else {
		return r.TargetColumn()
	}
}

func (r *AssociationRef) TypeEntity() *graph.Entity {
	entity := r.Association.TypeEntity()
	if entity != nil {
		return entity
	}

	panic("Can not find reference entity")
}

func (r *AssociationRef) IsCombination() bool {
	return r.IsSource() &&
		(r.Association.Relation.RelationType == meta.TWO_WAY_COMBINATION ||
			r.Association.Relation.RelationType == meta.ONE_WAY_COMBINATION)
}
