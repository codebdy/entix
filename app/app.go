package app

import "rxdrag.com/entify/model/graph"

var Installed = false

//该变量要在启动时加载，查询UUID是SYSTEM的 app
var SYSTEM_APP_ID uint64 = 1

//用来从数据库读取数据的原始Model，最初只包含Meta类 会逐步进化并包含log、user等类
var META_MODEL *graph.Model

var appCache = map[uint64]*AppSchema{}

func Get(appId uint64) *AppSchema {
	if appCache[appId] == nil {
		appCache[appId] = NewAppSchema(appId)
	}

	return appCache[appId]
}
