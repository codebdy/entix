package mysql

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/table"
)

func (b *MySQLBuilder) BuildUpdateSQL(id uint64, fields []*data.Field, assocs []*data.AssociationRef, table *table.Table) string {
	sql := fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE ID = %d",
		table.Name,
		updateSetFields(fields, assocs),
		id,
	)

	return sql
}

func updateSetFields(fields []*data.Field, assocs []*data.AssociationRef) string {
	if len(fields) == 0 && len(assocs) == 0 {
		log.Panic(errors.New("No update fields"))
	}
	fieldLen := len(fields)
	columns := make([]string, fieldLen+len(assocs))
	for i, field := range fields {
		columns[i] = field.Column.Name + "=?"
	}

	for i, assoc := range assocs {
		columns[fieldLen+i] = assoc.OwnerColumn().Name + "=?"
	}
	return strings.Join(columns, ",")
}
