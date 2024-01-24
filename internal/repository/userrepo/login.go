package userrepo

import (
	"context"
	"errors"
	"monitoring/internal/model"
	"monitoring/pkg/postgres"
)

type ILoginRepo interface {
	Get(ctx context.Context, username string) (user model.LoginResult, err error)
	getUserID(ctx context.Context, username string) (int, error)
}

type LoginRepo struct {
	DB postgres.IPostgres
}

func (lr *LoginRepo) Get(ctx context.Context, username, password string) (user model.LoginResult, err error) {

	userId, err := lr.getUserID(ctx, username)
	if err != nil {
		return model.LoginResult{}, err
	}

	user, err = lr.checkUserPass(ctx, userId, password)
	if err != nil {
		return model.LoginResult{}, err
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

	if userId > -1 {
		return userId, nil
	}

	return -1, errors.New("user not found")
}

func (lr *LoginRepo) checkUserPass(ctx context.Context, userId int, password string) (user model.LoginResult, err error) {

	q := `SELECT id, username, password, role FROM users WHERE id = $1`

	rows, err := lr.DB.QueryContext(ctx, q, userId)
	if err != nil {
		return model.LoginResult{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
		if err != nil {
			return model.LoginResult{}, err
		}
	}

	if userId != user.ID || user.Password != password {
		return model.LoginResult{}, errors.New("user/password is wrong")
	}

	return user, nil
}

// sundhar
