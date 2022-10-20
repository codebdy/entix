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
			Type:      ID,
			TypeLabel: ID,
			Uuid:      "APP_COLUMN_ID_UUID",
			System:    true,
		},
		{
			Name:      "uuid",
			Type:      String,
			TypeLabel: String,
			Uuid:      "APP_COLUMN_UUID_UUID",
			System:    true,
		},
		{
			Name:      "name",
			Type:      String,
			TypeLabel: String,
			Uuid:      "APP_COLUMN_NAME_UUID",
			System:    true,
		},
		{
			Name:      "meta",
			Type:      JSON,
			TypeLabel: JSON,
			Uuid:      "APP_COLUMN_META_UUID",
			System:    true,
		},
		{
			Name:      "publishedMeta",
			Type:      JSON,
			TypeLabel: JSON,
			Uuid:      "APP_COLUMN_PUBLISH_META_UUID",
			System:    true,
		},
		{
			Name:       "createdAt",
			Type:       Date,
			TypeLabel:  Date,
			CreateDate: true,
			Uuid:       "APP_COLUMN_CREATED_AT_UUID",
			System:     true,
		},
		{
			Name:       "saveMetaAt",
			Type:       Date,
			TypeLabel:  Date,
			CreateDate: true,
			Uuid:       "APP_COLUMN_SAVE_META_AT_UUID",
			System:     true,
		},
		{
			Name:       "publishMetaAt",
			Type:       Date,
			TypeLabel:  Date,
			CreateDate: true,
			Uuid:       "APP_COLUMN_PUBLISH_META_AT_UUID",
			System:     true,
		},
	},
	PackageUuid: PACKAGE_SYSTEM_UUID,
}

var UserClass = ClassMeta{
	PackageUuid: PACKAGE_SYSTEM_UUID,
	InnerId: USER_INNER_ID,
	Name: USER_ENTITY_NAME,
	Root: true,
	StereoType: CLASSS_ENTITY,
	Uuid: USER_ENTITY_UUID,
	System: true,
	Attributes: AttributeMeta{
		{
			System: true,
			Name: consts.ID,
			Primary: true,
			Type: ID,
			TypeLabel: ID,
			Uuid: "RX_USER_ID_UUID",
		},
		{
			System: true,
			Name: "name",
			Nullable: true,
			Type: String,
			TypeLabel: String,
			uuid: "RX_USER_NAME_UUID",
		},
		{
			System: true,
			Length: 128,
			Name: "loginName",
			Type: String,
			TypeLabel: String,
			uuid: "RX_USER_LOGINNAME_UUID"
		},
		{
			"System": true,
			"length": 256,
			"name": "password",
			"type": "String",
			"typeLabel": "String",
			"uuid": "RX_USER_PASSWORD_UUID"
		},
		{
			"System": true,
			"name": "isSupper",
			"nullable": true,
			"type": "Boolean",
			"typeLabel": "Boolean",
			"uuid": "RX_USER_ISSUPPER_UUID"
		},
		{
			"System": true,
			"name": "isDemo",
			"nullable": true,
			"type": "Boolean",
			"typeLabel": "Boolean",
			"uuid": "RX_USER_ISDEMO_UUID"
		},
		{
			"System": true,
			"createDate": true,
			"name": "createdAt",
			"type": "Date",
			"typeLabel": "Date",
			"uuid": "RX_USER_CREATEDAT_UUID"
		},
		{
			"System": true,
			"name": "updatedAt",
			"type": "Date",
			"typeLabel": "Date",
			"updateDate": true,
			"uuid": "RX_USER_UPDATEDAT_UUID"
		}
	},
}

var RoleClass = ClassMeta{
	PackageUuid: PACKAGE_SYSTEM_UUID,
}

var Relations = []RelationMeta{}


{
	"System": true,
	"attributes": [
		{
			"System": true,
			"name": "id",
			"primary": true,
			"type": "ID",
			"typeLabel": "ID",
			"uuid": "RX_ROLE_ID_UUID"
		},
		{
			"System": true,
			"name": "name",
			"type": "String",
			"typeLabel": "String",
			"uuid": "RX_ROLE_NAME_UUID"
		},
		{
			"System": true,
			"name": "description",
			"nullable": true,
			"type": "String",
			"typeLabel": "String",
			"uuid": "RX_ROLE_DESCRIPTION_UUID"
		},
		{
			"System": true,
			"createDate": true,
			"name": "createdAt",
			"type": "Date",
			"typeLabel": "Date",
			"uuid": "RX_ROLE_CREATEDAT_UUID"
		},
		{
			"System": true,
			"name": "updatedAt",
			"type": "Date",
			"typeLabel": "Date",
			"updateDate": true,
			"uuid": "RX_ROLE_META_UPDATEDAT_UUID"
		}
	],
	"innerId": 5,
	"name": "Role",
	"packageUuid": "PACKAGE_SYSTEM_UUID",
	"root": true,
	"stereoType": "Entity",
	"uuid": "META_ROLE_UUID"
},