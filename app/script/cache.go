package script

var cache = map[string]interface{}{}

func WriteToCache(key string, value interface{}) {
	cache[key] = value
}

func ReadFromCache(key string) interface{} {
	return cache[key]
}
