package meta

var SystemApp = map[string]interface{}{
	"uuid": SYSTEM_APP_UUID,
	"name": "Apper",
	"meta": MetaContent{
		Packages: []PackageMeta{
			{
				Name:   "System",
				System: true,
				Uuid:   PACKAGE_SYSTEM_UUID,
			},
		},
		Classes: []ClassMeta{
			AppClass,
			UserClass,
			RoleClass,
		},
		Relations: Relations,
	},
}
