package repository

import (
	"rxdrag.com/entify/config"
	"rxdrag.com/entify/db"
)

type Connection struct {
	idSeed int //use for sql join table
	Dbx    *db.Dbx
	v      *AbilityVerifier
	appId  uint64
}

func openWithConfig(cfg config.DbConfig, v *AbilityVerifier, appId uint64) (*Connection, error) {
	dbx, err := db.Open(cfg.Driver, DbString(cfg))
	if err != nil {
		return nil, err
	}
	con := Connection{
		idSeed: 1,
		Dbx:    dbx,
		v:      v,
		appId:  appId,
	}
	return &con, err
}

func Open(v *AbilityVerifier, appId uint64) (*Connection, error) {
	cfg := config.GetDbConfig()
	return openWithConfig(cfg, v, appId)
}

func (c *Connection) BeginTx() error {
	return c.Dbx.BeginTx()
}

func (c *Connection) Commit() error {
	return c.Dbx.Commit()
}

func (c *Connection) ClearTx() {
	c.Dbx.ClearTx()
}

//use for sql join table
func (c *Connection) CreateId() int {
	c.idSeed++
	return c.idSeed
}
