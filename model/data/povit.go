package data

import (
	"rxdrag.com/entify/model/table"
)

type AssociationPovit struct {
	Source *Field
	Target *Field
	//Fields      []*Field
	Association *Reference
}

func NewAssociationPovit(association *Reference, sourceId uint64, targetId uint64) *AssociationPovit {
	sourceColumn := association.SourceColumn()
	targetColumn := association.TargetColumn()
	povit := AssociationPovit{
		Association: association,
		Source: &Field{
			Column: sourceColumn,
			Value:  sourceId,
		},
		Target: &Field{
			Column: targetColumn,
			Value:  targetId,
		},
	}

	return &povit
}

func (a *AssociationPovit) Table() *table.Table {
	return a.Association.Table()
}

// func NewDerivedAssociationPovit(association Associationer, sourceId uint64, targetId uint64) *AssociationPovit {
// 	sourceColumn := association.SourceColumn()
// 	targetColumn := association.TargetColumn()
// 	povit := DerivedAssociationPovit{
// 		Association: association,
// 		source: &Field{
// 			Column: sourceColumn,
// 			Value:  sourceId,
// 		},
// 		target: &Field{
// 			Column: targetColumn,
// 			Value:  targetId,
// 		},
// 	}

// 	return &povit
// }
