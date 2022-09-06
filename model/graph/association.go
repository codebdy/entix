package graph

import "rxdrag.com/entify/model/meta"

type Association struct {
	Relation       *Relation
	OwnerClassUuid string
}

func NewAssociation(r *Relation, ownerUuid string) *Association {
	return &Association{
		Relation:       r,
		OwnerClassUuid: ownerUuid,
	}
}

func (a *Association) Name() string {
	if a.IsSource() {
		return a.Relation.RoleOfTarget
	} else {
		return a.Relation.RoleOfSource
	}
}

func (a *Association) Owner() *Entity {
	if a.IsSource() {
		return a.Relation.SourceEntity
	} else {
		return a.Relation.TargetEntity
	}
}

func (a *Association) TypeEntity() *Entity {
	if !a.IsSource() {
		return a.Relation.SourceEntity
	} else {
		return a.Relation.TargetEntity
	}
}

func (a *Association) Description() string {
	if a.IsSource() {
		return a.Relation.DescriptionOnTarget
	} else {
		return a.Relation.DescriptionOnSource
	}
}

func (a *Association) IsArray() bool {
	if a.IsSource() {
		return a.Relation.TargetMultiplicity == meta.ZERO_MANY
	} else {
		return a.Relation.SourceMutiplicity == meta.ZERO_MANY
	}
}

func (a *Association) IsSource() bool {
	return a.Relation.SourceEntity.Uuid() == a.OwnerClassUuid
}

func (a *Association) GetName() string {
	return a.Name()
}

func (a *Association) Path() string {
	return a.Owner().Domain.Name + "." + a.Name()
}
