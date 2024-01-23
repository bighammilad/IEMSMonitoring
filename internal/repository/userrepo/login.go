package userrepo

import (
	"context"
	"errors"
	"monitoring/internal/model"
)

type ILoginRepo interface {
	Get(ctx context.Context, username string) (user model.LoginResult, err error)
	getUserID(ctx context.Context, username string) (int, error)
}

type LoginRepo struct {
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

	if username == "jon" {
		return 12, nil
	}

	return -1, errors.New("user not found")
}

func (lr *LoginRepo) checkUserPass(ctx context.Context, userId int, password string) (user model.LoginResult, err error) {

	// get user from db
	var testUser model.LoginResult
	testUser.ID = userId
	testUser.Username = "jon"
	testUser.Password = "123"
	testUser.Role = 0

	if userId != testUser.ID && testUser.Password != password {
		return model.LoginResult{}, errors.New("user/password is wrong")
	}

	return testUser, nil
}
