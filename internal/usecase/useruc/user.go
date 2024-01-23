package useruc

import (
	"context"
	"monitoring/internal/repository/userrepo"
)

type IRegisterUserUsecase interface {
	CheckUsername(username string) (ok bool, err error)
	checkExistionUserById(username string) error
}

type UserUsecase struct {
	Repo *userrepo.UserRepo
}

func (ruu *UserUsecase) RegisterUser(ctx context.Context, username, password string, role int) (ok bool, err error) {

	ok, err = ruu.Repo.Register(ctx, username, password, role)
	if err != nil {
		return
	}
	if !ok {
		return
	}

	// Add user
	ok, err = ruu.Repo.Register(ctx, username, password, role)

	return
}

// func (ruu *RegisterUserUsecase) CheckUsername(ctx context.Context, username string) (ok bool, err error) {

// 	ok, err = ruu.Repo.CheckUsername(ctx, username)
// 	if err != nil {
// 		return ok, err
// 	}

// 	return ok, nil
// }
