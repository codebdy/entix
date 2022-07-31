package authentication

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"rxdrag.com/entify/authentication/jwt"
	"rxdrag.com/entify/common"
	"rxdrag.com/entify/db/dialect"
	"rxdrag.com/entify/repository"
	"rxdrag.com/entify/utils"
)

var TokenCache = map[string]*common.User{}

type Authentication struct {
}

func New() *Authentication {
	return &Authentication{}
}

func (a *Authentication) loadUser(loginName string) *common.User {
	con, err := repository.Open(nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var user common.User
	var isSupper sql.NullBool
	var isDemo sql.NullBool

	sqlBuilder := dialect.GetSQLBuilder()
	err = con.Dbx.QueryRow(sqlBuilder.BuildMeSQL(), loginName).Scan(
		&user.Id,
		&user.Name,
		&user.LoginName,
		&isSupper,
		&isDemo,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err.Error())
	}

	user.IsSupper = isSupper.Bool
	user.IsDemo = isDemo.Bool

	rows, err := con.Dbx.Query(sqlBuilder.BuildRolesSQL(), user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var role common.Role
		err = rows.Scan(&role.Id, &role.Name)
		if err != nil {
			panic(err.Error())
		}
		user.Roles = append(user.Roles, role)
	}
	return &user
}

func (a *Authentication) CheckPassword(loginName, pwd string) (bool, error) {
	con, err := repository.Open(nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	sqlBuilder := dialect.GetSQLBuilder()
	var password string
	err = con.Dbx.QueryRow(sqlBuilder.BuildLoginSQL(), loginName).Scan(&password)
	if err != nil {
		fmt.Println(err)
		return false, errors.New("Login failed!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(pwd)) //验证（对比）
	if err != nil {
		fmt.Println(err, pwd, password)
		return false, errors.New("Password error!")
	}

	return true, nil
}

func (a *Authentication) GenerateToken(loginName string) (string, error) {
	token, err := jwt.GenerateToken(loginName)
	if err != nil {
		panic(err.Error())
	}

	user := a.loadUser(loginName)
	TokenCache[token] = user
	return token, err
}

func (a *Authentication) Login(loginName, pwd string) (string, error) {
	if _, err := a.CheckPassword(loginName, pwd); err != nil {
		return "", err
	}

	return a.GenerateToken(loginName)
}

func (a *Authentication) ChangePassword(loginName, oldPassword, newPassword string) (string, error) {
	if _, err := a.CheckPassword(loginName, oldPassword); err != nil {
		return "", err
	}

	con, err := repository.Open(nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	sqlBuilder := dialect.GetSQLBuilder()
	result, err := con.Dbx.Exec(
		sqlBuilder.BuildChangePasswordSQL(),
		utils.BcryptEncode(newPassword),
		loginName,
	)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("Login failed!")
	}

	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		return "", errors.New("Change password failed!")
	}

	return a.GenerateToken(loginName)
}

func Logout(token string) {
	TokenCache[token] = nil
}

func GetUserByToken(token string) (*common.User, error) {
	return TokenCache[token], nil
}
