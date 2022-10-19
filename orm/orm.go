package orm

import "rxdrag.com/entify/model/graph"

func Open() Session {
	return Session{}
}

func Migrage(newGraph *graph.Model) {
	//取旧数据，做diff
}
