package repository

import (
	"rxdrag.com/entify/config"
	"rxdrag.com/entify/db"
)

type Session struct {
	idSeed int //use for sql join table
	Dbx    *db.Dbx
	v      *AbilityVerifier
	appId  uint64
}

func openWithConfig(cfg config.DbConfig, v *AbilityVerifier, appId uint64) (*Session, error) {
	dbx, err := db.Open(cfg.Driver, DbString(cfg))
	if err != nil {
		return nil, err
	}
	con := Session{
		idSeed: 1,
		Dbx:    dbx,
		v:      v,
		appId:  appId,
	}
	return &con, err
}

func Open(v *AbilityVerifier, appId uint64) (*Session, error) {
	cfg := config.GetDbConfig()
	return openWithConfig(cfg, v, appId)
}
