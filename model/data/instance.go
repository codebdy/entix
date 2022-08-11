package data

import (
	"fmt"
	"strconv"
	"time"

	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/table"
)

type Field struct {
	Column *table.Column
	Value  interface{}
}

type Instance struct {
	Id           uint64
	Entity       *graph.Entity
	Fields       []*Field
	Associations []*Reference
}

func NewInstance(object map[string]interface{}, entity *graph.Entity) *Instance {
	instance := Instance{
		Entity: entity,
	}
	if object[consts.ID] != nil {
		instance.Id = parseId(object[consts.ID])
	}

	columns := entity.Table.Columns
	for i := range columns {
		column := columns[i]
		if object[column.Name] != nil {
			instance.Fields = append(instance.Fields, &Field{
				Column: column,
				Value:  object[column.Name],
			})
		} else if column.CreateDate || column.UpdateDate {
			instance.Fields = append(instance.Fields, &Field{
				Column: column,
				Value:  time.Now(),
			})
		}
	}
	allAssociation := entity.AllAssociations()
	for i := range allAssociation {
		asso := allAssociation[i]
		value := object[asso.Name()]
		if value != nil {
			ref := Reference{
				Association: asso,
				Value:       value.(map[string]interface{}),
			}
			instance.Associations = append(instance.Associations, &ref)
		}
	}
	return &instance
}

func (ins *Instance) IsInsert() bool {
	for i := range ins.Fields {
		field := ins.Fields[i]
		if field.Column.Name == consts.ID {
			if field.Value != nil {
				return false
			}
		}
	}
	return true
}

func (ins *Instance) Table() *table.Table {
	return ins.Entity.Table
}

func parseId(id interface{}) uint64 {
	switch v := id.(type) {
	default:
		panic(fmt.Sprintf("unexpected id type %T", v))
	case uint64:
		return id.(uint64)
	case string:
		if id.(string) == "" {
			return 0
		}
		u, err := strconv.ParseUint(id.(string), 0, 64)
		if err != nil {
			panic(err.Error())
		}
		return u
	}
}
