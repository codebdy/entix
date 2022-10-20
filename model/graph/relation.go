package graph

import (
	"rxdrag.com/entify/model/domain"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/model/table"
)

type Relation struct {
	AppId                  uint64
	Uuid                   string
	InnerId                uint64
	RelationType           string
	SourceEntity           *Entity
	TargetEntity           *Entity
	RoleOfTarget           string
	RoleOfSource           string
	DescriptionOnSource    string
	DescriptionOnTarget    string
	SourceMutiplicity      string
	TargetMultiplicity     string
	EnableAssociaitonClass bool
	AssociationClass       meta.AssociationClass
	Table                  *table.Table
}

func NewRelation(
	r *domain.Relation,
	sourceEntity *Entity,
	targetEntity *Entity,
) *Relation {
	relation := &Relation{
		Uuid:                   r.Uuid,
		InnerId:                r.InnerId,
		RelationType:           r.RelationType,
		SourceEntity:           sourceEntity,
		TargetEntity:           targetEntity,
		RoleOfTarget:           r.RoleOfTarget,
		RoleOfSource:           r.RoleOfSource,
		DescriptionOnSource:    r.DescriptionOnSource,
		DescriptionOnTarget:    r.DescriptionOnTarget,
		SourceMutiplicity:      r.SourceMutiplicity,
		TargetMultiplicity:     r.TargetMultiplicity,
		EnableAssociaitonClass: r.EnableAssociaitonClass,
		AssociationClass:       r.AssociationClass,
		AppId:                  r.AppId,
	}

	return relation
}
