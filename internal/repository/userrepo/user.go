package userrepo

import (
	"context"
	"errors"
	"monitoring/internal/model"
	hashPass "monitoring/pkg/hashPass"
	"monitoring/pkg/postgres"
)

type IUser interface {
	Register(ctx context.Context, username, hashPass string, role int) (ok bool, err error)
	Update(ctx context.Context, username, hashPass string, role int) (ok bool, err error)
}

type UserRepo struct {
	DB postgres.IPostgres
}

func (ru *UserRepo) Register(ctx context.Context, username, hashPass string, role int) (ok bool, err error) {

	ok, err = ru.checkUsername(ctx, username)
	if err != nil {
		return ok, err
	}

	if ok {
		user := model.UserRegister{Username: username, Password: hashPass, Role: role}
		err := ru.addUser(ctx, user)
		if err != nil {
			return false, err
		}
	}

	return
}

func (ru *UserRepo) Update(ctx context.Context, username, hashPass string, role int) (ok bool, err error) {

	ok, err = ru.checkUserPass(ctx, username, hashPass)
	if err != nil {
		return ok, err
	}

	if ok {
		user := model.UserUpdate{Username: username, Password: hashPass, Role: role}
		err := ru.updateUser(ctx, user)
		if err != nil {
			return false, err
		}
	}

	return
}

func (ru *UserRepo) checkUsername(ctx context.Context, username string) (ok bool, err error) {

	q := `select count(*) from users where username=$1;`
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

func (ru *UserRepo) checkUserPass(ctx context.Context, username, password string) (ok bool, err error) {

	q := `select count(*) from users where username=$1;`
	rows, err := ru.DB.QueryContext(ctx, q)
	if err != nil {
		return false, err
	}

	var user model.UserUpdate
	defer rows.Close()
	for rows.Next() {
		var entry model.UserUpdate
		err := rows.Scan(
			&entry.Username,
			&entry.Password,
			&entry.Role,
		)
		if err != nil {
			return false, err
		}
		user = entry
	}

	if user.Username != username || hashPass.CheckPasswordHash(password, user.Password) {
		return false, errors.New("username or password isn't match")
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

func (ruu *UserRepo) GetUsrId(ctx context.Context, username string) (userid int, err error) {
	q := `SELECT id FROM users WHERE username=$1;`
	rows, err := ruu.DB.QueryContext(ctx, q, username)
	if err != nil {
		return 0, err
	}
	var user model.UserRes
	defer rows.Close()
	for rows.Next() {
		var entry model.UserRes
		err := rows.Scan(
			&entry.UserId,
		)
		if err != nil {
			return 0, err
		}
		user = entry
	}
	return user.UserId, nil
}
