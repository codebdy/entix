package service

import (
	"log"

	"rxdrag.com/entify/consts"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/orm"
)

func ImportApp(app *data.Instance) error {
	session, err := orm.Open()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	//根据uuid查出App，如果App存在，更新；不存在，新建

	oldApp := session.QueryOneEntity(
		app.Entity,
		map[string]interface{}{
			consts.ARG_WHERE: map[string]interface{}{
				"uuid": map[string]interface{}{
					consts.ARG_EQ: app.ValueMap["uuid"],
				},
			},
		},
	)

	err = session.BeginTx()
	defer session.ClearTx()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if oldApp != nil {
		_, err = session.SaveOne(app)
	} else {

	}

	if err != nil {
		log.Println(err.Error())
		session.Dbx.Rollback()
		return err
	}
	err = session.Commit()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
