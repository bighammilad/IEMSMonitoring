package userrepo

import (
	"context"
	"errors"
	"monitoring/internal/model"
	"monitoring/pkg/postgres"
)

type IUser interface {
	Create(ctx context.Context, username, hashPass string, role int) (ok bool, err error)
	Read(ctx context.Context, username string) (user model.UserRes, err error)
	ReadAll(ctx context.Context) (users []model.UserRes, err error)
	Update(ctx context.Context, username, hashPass string, role int) (ok bool, err error)
	Delete(ctx context.Context, username string) error
	GetUsrId(ctx context.Context, username string) (userid int, err error)
}

type UserRepo struct {
	DB postgres.IPostgres
}

func (ur *UserRepo) Create(ctx context.Context, username, hashPass string, role int) (ok bool, err error) {

	q := `INSERT INTO users ("username", "password", "role") VALUES (:username, :password, :role) returning *;`
	ok, err = ur.checkUsername(ctx, username)
	if err != nil {
		return ok, err
	}
	if ok {
		user := model.UserRegister{Username: username, Password: hashPass, Role: role}
		_, err = ur.DB.NamedExec(ctx, user, q)
		if err != nil {
			return false, err
		}
	}

	return
}

func (ur *UserRepo) Read(ctx context.Context, username string) (user model.UserRes, err error) {
	q := `SELECT id,username,role FROM users WHERE username=$1;`
	rows, err := ur.DB.QueryContext(ctx, q, username)
	if err != nil {
		return user, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry model.UserRes
		err := rows.Scan(
			&entry.UserId,
			&entry.Username,
			&entry.Role,
		)
		if err != nil {
			return user, err
		}
		user = entry
	}
	return user, nil
}

func (ur *UserRepo) ReadAll(ctx context.Context) (users []model.UserRes, err error) {
	q := `SELECT id,username,role FROM users;`
	rows, err := ur.DB.QueryContext(ctx, q)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry model.UserRes
		err := rows.Scan(
			&entry.UserId,
			&entry.Username,
			&entry.Role,
		)
		if err != nil {
			return users, err
		}
		users = append(users, entry)
	}
	return users, nil
}

func (ur *UserRepo) Update(ctx context.Context, username, hashPass string, role int) (ok bool, err error) {

	first_q := `UPDATE users
		SET username=$1, password=$2, role=$3
		WHERE username=$1`

	second_q := `UPDATE users
		SET username=$1, role=$2
		WHERE username=$1`

	third_q := `UPDATE users
		SET username=$1, password=$2
		WHERE username=$1`

	fourth_q := `UPDATE users
		SET username=$1
		WHERE username=$1`

	switch {
	case username != "" && hashPass != "" && role != 0:
		_, err = ur.DB.ExecContext(ctx, first_q, username, hashPass, role)
	case username != "" && role != 0 && hashPass == "":
		_, err = ur.DB.ExecContext(ctx, second_q, username, role)
	case username != "" && hashPass != "" && role == 0:
		_, err = ur.DB.ExecContext(ctx, third_q, username, hashPass)
	case username != "" && hashPass == "" && role == 0:
		_, err = ur.DB.ExecContext(ctx, fourth_q, username)
	default:
		return false, errors.New("no fields to update")
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (ur *UserRepo) Delete(ctx context.Context, username string) error {
	_, err := ur.DB.ExecContext(ctx, "DELETE FROM users WHERE username=$1", username)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) checkUsername(ctx context.Context, username string) (ok bool, err error) {

	q := `select count(*) as cnt from users where username=$1;`
	rows, err := ur.DB.QueryContext(ctx, q, username)
	if err != nil {
		return false, err
	}

	var usercnt int
	defer rows.Close()
	for rows.Next() {
		type entry struct {
			cnt int
		}
		var ntry entry
		err := rows.Scan(
			&ntry.cnt,
		)
		if err != nil {
			return false, err
		}
		usercnt = ntry.cnt
	}

	if usercnt > 0 {
		return false, errors.New("this username is already taken")
	}

	return true, nil
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
