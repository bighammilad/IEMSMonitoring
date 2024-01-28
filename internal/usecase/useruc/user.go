package useruc

import (
	"context"
	"monitoring/internal/repository/userrepo"
	hashPass "monitoring/pkg/hashPass"
)

type IRegisterUserUsecase interface {
	CheckUsername(username string) (ok bool, err error)
	checkExistionUserById(username string) error
}

type UserUsecase struct {
	DB *userrepo.UserRepo
}

func (ruu *UserUsecase) RegisterUser(ctx context.Context, username, password string, role int) (ok bool, err error) {

	hashPass, err := hashPass.HashPassword(password)
	if err != nil {
		return false, err
	}
	ok, err = ruu.DB.Register(ctx, username, hashPass, role)
	if err != nil {
		return false, err
	}
	return ok, err
}
