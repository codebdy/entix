package service

import (
	"fmt"

	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/orm"
)

func QueryInterface(intf *graph.Interface, args graph.QueryArg) map[string]interface{} {
	session, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return session.QueryInterface(intf, args)
}

func QueryOneInterface(intf *graph.Interface, args graph.QueryArg) interface{} {
	session, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return session.QueryOneInterface(intf, args)
}

func QueryEntity(entity *graph.Entity, args graph.QueryArg) map[string]interface{} {
	session, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return session.QueryEntity(entity, args)
}

func QueryOneEntity(entity *graph.Entity, args graph.QueryArg) interface{} {
	session, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return session.QueryOneEntity(entity, args)
}

func DeleteInstances(instances []*data.Instance) (interface{}, error) {
	session, err := orm.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = session.BeginTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer session.ClearTx()

	deletedIds := []interface{}{}

	for i := range instances {
		instance := instances[i]
		session.DeleteInstance(instance)
		deletedIds = append(deletedIds, instance.Id)
	}

	err = session.Dbx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return deletedIds, nil
}

func DeleteInstance(instance *data.Instance) (interface{}, error) {
	session, err := orm.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = session.BeginTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer session.ClearTx()
	session.DeleteInstance(instance)

	err = session.Dbx.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return instance.Id, nil
}

func Save(instances []*data.Instance) ([]interface{}, error) {
	session, err := orm.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = session.BeginTx()
	defer session.ClearTx()
	if err != nil {
		fmt.Println(err.Error())
		session.Dbx.Rollback()
		return nil, err
	}
	saved := []interface{}{}

	for i := range instances {
		obj, err := session.SaveOne(instances[i])
		if err != nil {
			fmt.Println(err.Error())
			session.Dbx.Rollback()
			return nil, err
		}

		saved = append(saved, obj)
	}

	err = session.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return saved, nil
}

func SaveOne(instance *data.Instance) (interface{}, error) {
	session, err := orm.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = session.BeginTx()
	defer session.ClearTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	obj, err := session.SaveOne(instance)
	if err != nil {
		fmt.Println(err.Error())
		session.Dbx.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return obj, nil
}

func InsertOne(instance *data.Instance) (interface{}, error) {
	session, err := orm.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	session.BeginTx()
	defer session.ClearTx()
	if err != nil {
		fmt.Println(err.Error())
		session.Dbx.Rollback()
		return nil, err
	}

	obj, err := session.InsertOne(instance)
	if err != nil {
		fmt.Println(err.Error())
		session.Dbx.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		fmt.Println(err.Error())
		session.Dbx.Rollback()
		return nil, err
	}
	return obj, nil
}

func BatchQueryAssociations(
	association *graph.Association,
	ids []uint64,
	args graph.QueryArg,
) []map[string]interface{} {
	con, err := orm.Open()
	if err != nil {
		panic(err.Error())
	}
	return con.BatchRealAssociations(association, ids, args)
}
