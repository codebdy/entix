package meta

import (
	"rxdrag.com/entify/consts"
)

var AppClass = ClassMeta{
	Uuid:       APP_ENTITY_UUID,
	Name:       APP_ENTITY_NAME,
	InnerId:    APP_INNER_ID,
	StereoType: CLASSS_ENTITY,
	IdNoShift:  true,
	Root:       true,
	System:     true,
	Attributes: []AttributeMeta{
		{
			Name:      consts.ID,
			Primary:   true,
			Type:      "ID",
			TypeLabel: "ID",
			Uuid:      "APP_COLUMN_ID_UUID",
			System:    true,
		},
		{
			Name:      "uuid",
			Type:      "String",
			TypeLabel: "String",
			Uuid:      "APP_COLUMN_UUID_UUID",
			System:    true,
		},
		{
			Name:      "name",
			Type:      "String",
			TypeLabel: "String",
			Uuid:      "APP_COLUMN_NAME_UUID",
			System:    true,
		},
		{
			Name:      "meta",
			Type:      "JSON",
			TypeLabel: "JSON",
			Uuid:      "APP_COLUMN_META_UUID",
			System:    true,
		},
		{
			Name:      "publishedMeta",
			Type:      "JSON",
			TypeLabel: "JSON",
			Uuid:      "APP_COLUMN_PUBLISH_META_UUID",
			System:    true,
		},
		{
			Name:       "createdAt",
			Type:       "Date",
			TypeLabel:  "Date",
			CreateDate: true,
			Uuid:       "APP_COLUMN_CREATED_AT_UUID",
			System:     true,
		},
		{
			Name:       "saveMetaAt",
			Type:       "Date",
			TypeLabel:  "Date",
			CreateDate: true,
			Uuid:       "APP_COLUMN_SAVE_META_AT_UUID",
			System:     true,
		},
		{
			Name:       "publishMetaAt",
			Type:       "Date",
			TypeLabel:  "Date",
			CreateDate: true,
			Uuid:       "APP_COLUMN_PUBLISH_META_AT_UUID",
			System:     true,
		},
	},
	PackageUuid: PACKAGE_SYSTEM_UUID,
}

var EntityAuthSettingsClass = ClassMeta{
	Name:    "EntityAuthSettings",
	Uuid:    "META_ENTITY_AUTH_SETTINGS_UUID",
	InnerId: ENTITY_AUTH_SETTINGS_INNER_ID,
	Root:    true,
	System:  true,
	Attributes: []AttributeMeta{
		{
			Name:    consts.ID,
			Type:    ID,
			Uuid:    "RX_ENTITY_AUTH_SETTINGS_ID_UUID",
			Primary: true,
			System:  true,
		},
		{
			Name:   "entityUuid",
			Type:   "String",
			Uuid:   "RX_ENTITY_AUTH_SETTINGS_ENTITY_UUID_UUID",
			System: true,
			Unique: true,
		},
		{
			Name:   "expand",
			Type:   "Boolean",
			Uuid:   "RX_ENTITY_AUTH_SETTINGS_EXPAND_UUID",
			System: true,
		},
	},
	StereoType: "Entity",
}

var AbilityTypeEnum = ClassMeta{
	Uuid:       META_ABILITY_TYPE_ENUM_UUID,
	Name:       "AbilityType",
	StereoType: ENUM,
	Attributes: []AttributeMeta{
		{
			Name: META_ABILITY_TYPE_CREATE,
		},
		{
			Name: META_ABILITY_TYPE_READ,
		},
		{
			Name: META_ABILITY_TYPE_UPDATE,
		},
		{
			Name: META_ABILITY_TYPE_DELETE,
		},
	},
}

var AbilityClass = ClassMeta{
	Name:    "Ability",
	Uuid:    ABILITY_UUID,
	InnerId: Ability_INNER_ID,
	Root:    true,
	System:  true,
	Attributes: []AttributeMeta{
		{
			Name:    consts.ID,
			Type:    ID,
			Uuid:    "RX_ABILITY_ID_UUID",
			Primary: true,
			System:  true,
		},
		{
			Name:   "entityUuid",
			Type:   "String",
			Uuid:   "RX_ABILITY_ENTITY_UUID_UUID",
			System: true,
		},
		{
			Name:   "columnUuid",
			Type:   "String",
			Uuid:   "RX_ABILITY_COLUMN_UUID_UUID",
			System: true,
		},
		{
			Name:   "can",
			Type:   "Boolean",
			Uuid:   "RX_ABILITY_CAN_UUID",
			System: true,
		},
		{
			Name:   "expression",
			Type:   "String",
			Uuid:   "RX_ABILITY_EXPRESSION_UUID",
			Length: 2000,
			System: true,
		},
		{
			Name:     "abilityType",
			Type:     ENUM,
			Uuid:     "RX_ABILITY_ABILITYTYPE_UUID",
			System:   true,
			TypeUuid: META_ABILITY_TYPE_ENUM_UUID,
		},
		{
			Name:   "roleId",
			Type:   ID,
			Uuid:   "RX_ABILITY_ROLE_ID_UUID",
			System: true,
		},
	},
	StereoType: "Entity",
}
