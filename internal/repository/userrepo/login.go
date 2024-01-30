package userrepo

import (
	"context"
	"errors"
	"monitoring/internal/model"
	hashpass "monitoring/pkg/hashPass"
	"monitoring/pkg/postgres"
)

type ILoginRepo interface {
	Auth(ctx context.Context, username, password string) (user model.LoginResult, err error)
	getUserID(ctx context.Context, username string) (int, error)
	auth(ctx context.Context, userId int, password string) (user model.LoginResult, err error)
}

type LoginRepo struct {
	DB postgres.IPostgres
}

func (lr *LoginRepo) Auth(ctx context.Context, username, password string) (user model.LoginResult, err error) {
	userId, err := lr.getUserID(ctx, username)
	if err != nil {
		return model.LoginResult{}, err
	}
	if userId <= 0 {
		return model.LoginResult{}, errors.New("user not found")
	}
	user, err = lr.auth(ctx, userId, password)
	if err != nil {
		if err.Error() == "user/password is wrong" {
			return model.LoginResult{}, errors.New("incorrect password")
		} else {
			return model.LoginResult{}, err
		}
	}
	return user, nil
}

func (lr *LoginRepo) getUserID(ctx context.Context, username string) (int, error) {

	q := `SELECT id FROM users WHERE username = $1`
	row, err := lr.DB.QueryContext(ctx, q, username)
	if err != nil {
		return -1, err
	}
	defer row.Close()
	var userId int
	for row.Next() {
		err = row.Scan(&userId)
		if err != nil {
			return -1, err
		}
	}
	if userId <= 0 {
		return -1, errors.New("user not found")
	}
	return userId, nil
}

func (lr *LoginRepo) auth(ctx context.Context, userId int, password string) (user model.LoginResult, err error) {

	q := `SELECT username, password, role FROM users WHERE id = $1`
	rows, err := lr.DB.QueryContext(ctx, q, userId)
	if err != nil {
		return model.LoginResult{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.Username, &user.Password, &user.Role)
		if err != nil {
			return model.LoginResult{}, err
		}
	}
	ok := hashpass.CheckPasswordHash(password, user.Password)
	if !ok {
		return model.LoginResult{}, errors.New("user/password is wrong")
	}
	return user, nil
}
