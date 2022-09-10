package repository

import (
	"fmt"

	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/model"
	"rxdrag.com/entify/model/data"
	"rxdrag.com/entify/model/graph"
)

type Repository struct {
	Model *model.Model
	V     *AbilityVerifier
}

func New(model *model.Model) *Repository {
	return &Repository{
		Model: model,
	}
}

func (r *Repository) QueryInterface(intf *graph.Interface, args graph.QueryArg) map[string]interface{} {
	con, err := Open(r.V)
	if err != nil {
		panic(err.Error())
	}
	return con.doQueryInterface(intf, args)
}

func (r *Repository) QueryOneInterface(intf *graph.Interface, args graph.QueryArg) interface{} {
	con, err := Open(r.V)
	if err != nil {
		panic(err.Error())
	}
	return con.doQueryOneInterface(intf, args)
}

func (r *Repository) QueryEntity(entity *graph.Entity, args graph.QueryArg) map[string]interface{} {
	con, err := Open(r.V)
	if err != nil {
		panic(err.Error())
	}
	return con.doQueryEntity(entity, args)
}

func (r *Repository) QueryOneEntity(entity *graph.Entity, args graph.QueryArg) interface{} {
	con, err := Open(r.V)
	if err != nil {
		panic(err.Error())
	}
	return con.doQueryOneEntity(entity, args)
}

func (r *Repository) DeleteInstances(instances []*data.Instance) (interface{}, error) {
	con, err := Open(r.V)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.BeginTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer con.ClearTx()

	for i := range instances {
		con.doDeleteInstance(instances[i])
	}
	return instances, nil
}

func (r *Repository) DeleteInstance(instance *data.Instance) (interface{}, error) {
	con, err := Open(r.V)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.BeginTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer con.ClearTx()
	con.doDeleteInstance(instance)
	return instance.Id, nil
}

func (r *Repository) SaveOne(instance *data.Instance) (interface{}, error) {
	con, err := Open(r.V)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.BeginTx()
	defer con.ClearTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	obj, err := con.doSaveOne(instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return obj, nil
}

func (r *Repository) InsertOne(instance *data.Instance) (interface{}, error) {
	con, err := Open(r.V)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer con.ClearTx()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	obj, err := con.doInsertOne(instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = con.Commit()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return obj, nil
}

func (r *Repository) BatchQueryAssociations(
	association *graph.Association,
	ids []uint64,
	args graph.QueryArg,
) []map[string]interface{} {
	con, err := Open(r.V)
	if err != nil {
		panic(err.Error())
	}
	return con.doBatchRealAssociations(association, ids, args, r.V)
}

func IsEntityExists(name string) bool {
	con, err := Open(nil)
	if err != nil {
		panic(err.Error())
	}
	return con.doCheckEntity(name)
}

func InstallMeta() error {
	sqlBuilder := dialect.GetSQLBuilder()
	con, err := Open(nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = con.BeginTx()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = con.Dbx.Exec(sqlBuilder.BuildCreateMetaSQL())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = con.Dbx.Exec(sqlBuilder.BuildCreateAbilitySQL())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = con.Dbx.Exec(sqlBuilder.BuildCreateEntityAuthSettingsSQL())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = con.Commit()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
