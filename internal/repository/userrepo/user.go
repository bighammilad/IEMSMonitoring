package userrepo

import (
	"context"
	"errors"
	"monitoring/internal/model"
	"monitoring/pkg/postgres"
)

type IUser interface {
	Register(ctx context.Context, username, password string, role int) (ok bool, err error)
	Update(ctx context.Context, username, password string, role int) (ok bool, err error)
}

type UserRepo struct {
	DB *postgres.Postgres
}

func (ru *UserRepo) Register(ctx context.Context, username, password string, role int) (ok bool, err error) {

	ok, err = ru.checkUsername(ctx, username)
	if err != nil {
		return ok, err
	}

	if ok {
		user := model.UserRegister{Username: username, Password: password, Role: role}
		err := ru.addUser(ctx, user)
		if err != nil {
			return false, err
		}
	}

	return
}

func (ru *UserRepo) Update(ctx context.Context, username, password string, role int) (ok bool, err error) {

	ok, err = ru.checkUsername(ctx, username)
	if err != nil {
		return ok, err
	}

	if ok {
		user := model.UserUpdate{Username: username, Password: password, Role: role}
		err := ru.updateUser(ctx, user)
		if err != nil {
			return false, err
		}
	}

	return
}

func (ru *UserRepo) checkUsername(ctx context.Context, username string) (ok bool, err error) {

	// q := `select count(*) from users where username=$1;`
	q := `select 1`
	rows, err := ru.DB.QueryContext(ctx, q)
	if err != nil {
		return false, err
	}

	users := make([]model.UserRes, 0)
	defer rows.Close()
	for rows.Next() {
		var entry model.UserRes
		err := rows.Scan(
			&entry.UserId,
			&entry.Username,
			&entry.Role,
		)
		if err != nil {
			return false, err
		}
		users = append(users, entry)
	}

	if len(users) > 0 {
		return false, errors.New("this username is already taken")
	}

	return true, nil
}

func (ruu *UserRepo) addUser(ctx context.Context, user model.UserRegister) (err error) {
	q := `INSERT INTO users ("username", "password", "role") VALUES (:username, :password, :role) returning *;`
	_, err = ruu.DB.NamedExec(ctx, user, q)
	if err != nil {
		return err
	}

	return
}

func (ruu *UserRepo) updateUser(ctx context.Context, user model.UserUpdate) (err error) {
	q := `UPDATE users 
	 	  SET username=$1, password=$2, role=$3)
		  WHERE username=$1`
	_, err = ruu.DB.NamedExec(ctx, user, q)
	if err != nil {
		return err
	}

	return
}
