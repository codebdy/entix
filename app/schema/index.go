package schema

var Installed = false

var schemaCache = map[string]*AppSchema{}

func Get(appUuid string) *AppSchema {
	if schemaCache[appUuid] == nil {
		schemaCache[appUuid] = NewAppSchema(appUuid)
	}

	return schemaCache[appUuid]
}
