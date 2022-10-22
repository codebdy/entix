package service

import (
	"log"

	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/orm"
)

func Save(instances []*data.Instance) ([]interface{}, error) {
	session, err := orm.Open()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = session.BeginTx()
	defer session.ClearTx()
	if err != nil {
		log.Println(err.Error())
		session.Dbx.Rollback()
		return nil, err
	}
	saved := []interface{}{}

	for i := range instances {
		obj, err := session.SaveOne(instances[i])
		if err != nil {
			log.Println(err.Error())
			session.Dbx.Rollback()
			return nil, err
		}

		saved = append(saved, obj)
	}

	err = session.Commit()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return saved, nil
}

func SaveOne(instance *data.Instance) (interface{}, error) {
	session, err := orm.Open()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = session.BeginTx()
	defer session.ClearTx()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	obj, err := session.SaveOne(instance)
	if err != nil {
		log.Println(err.Error())
		session.Dbx.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return obj, nil
}

func InsertOne(instance *data.Instance) (interface{}, error) {
	instance.AsInsert()
	return SaveOne(instance)
}
