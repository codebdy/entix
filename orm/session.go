package orm

import (
	"database/sql"
	"log"

	"rxdrag.com/entify/config"
	"rxdrag.com/entify/db"
	"rxdrag.com/entify/db/dialect"
)

type Session struct {
	idSeed int //use for sql join table
	Dbx    *db.Dbx
}

func (con *Session) doCheckEntity(name string) bool {
	sqlBuilder := dialect.GetSQLBuilder()
	var count int
	err := con.Dbx.QueryRow(sqlBuilder.BuildTableCheckSQL(name, config.GetDbConfig().Database)).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Panic(err)
	}
	return count > 0
}
