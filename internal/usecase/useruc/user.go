package useruc

import (
	"context"
	"monitoring/internal/model"
	"monitoring/internal/repository/userrepo"
	hashPass "monitoring/pkg/hashPass"
)

type IUserUsecase interface {
	Create(ctx context.Context, username, password string, role int) (ok bool, err error)
	Read(ctx context.Context, username string) (user model.UserRes, err error)
	ReadAll(ctx context.Context) (users []model.UserRes, err error)
	Update(ctx context.Context, username, password string, role int) (ok bool, err error)
	Delete(ctx context.Context, username string) error
	GetUsrId(ctx context.Context, username string) (userid int, err error)
}

type UserUsecase struct {
	IUserRepo userrepo.IUserRepo
}

func (ruu *UserUsecase) Create(ctx context.Context, username, password string, role int) (ok bool, err error) {

	hashPass, err := hashPass.HashPassword(password)
	if err != nil {
		return false, err
	}
	ok, err = ruu.IUserRepo.Create(ctx, username, hashPass, role)
	if err != nil {
		return false, err
	}
	return ok, err
}

func (ruu *UserUsecase) Update(ctx context.Context, username, password string, role int) (ok bool, err error) {

	var hashpass string
	if password != "" {
		hashpass, err = hashPass.HashPassword(password)
		if err != nil {
			return false, err
		}
	}
	ok, err = ruu.IUserRepo.Update(ctx, username, hashpass, role)
	if err != nil {
		return false, err
	}
	return ok, err
}

func (ruu *UserUsecase) Delete(ctx context.Context, username string) error {
	err := ruu.IUserRepo.Delete(ctx, username)
	if err != nil {
		return err
	}
	return nil
}

func (ruu *UserUsecase) Read(ctx context.Context, username string) (user model.UserRes, err error) {
	user, err = ruu.IUserRepo.Read(ctx, username)
	if err != nil {
		return user, err
	}
	return user, err
}

func (ruu *UserUsecase) ReadAll(ctx context.Context) (users []model.UserRes, err error) {
	users, err = ruu.IUserRepo.ReadAll(ctx)
	if err != nil {
		return users, err
	}
	return users, err
}

func (ruu *UserUsecase) GetUsrId(ctx context.Context, username string) (userid int, err error) {
	userid, err = ruu.IUserRepo.GetUsrId(ctx, username)
	if err != nil {
		return userid, err
	}
	return userid, err
}
