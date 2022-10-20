package orm

import (
	"fmt"
	"log"

	"rxdrag.com/entify/config"
	"rxdrag.com/entify/db"
	"rxdrag.com/entify/model"
)

func DbString(cfg config.DbConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

func Open(model *model.Model) (*Session, error) {
	cfg := config.GetDbConfig()
	dbx, err := db.Open(cfg.Driver, DbString(cfg))
	if err != nil {
		return nil, err
	}
	session := Session{
		idSeed: 1,
		Dbx:    dbx,
		model:  model,
	}
	return &session, nil
}

func IsEntityExists(name string) bool {
	session, err := Open(nil)
	if err != nil {
		log.Panic(err)
	}
	return session.doCheckEntity(name)
}
