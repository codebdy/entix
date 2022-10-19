package orm

import "rxdrag.com/entify/db"

type Session struct {
	idSeed int //use for sql join table
	Dbx    *db.Dbx
}
