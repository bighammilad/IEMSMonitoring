package usecase

import (
	"context"
	"monitoring/internal/model"
	"monitoring/internal/repository"
)

type IUserService interface {
	Add(ctx context.Context, userservice model.UserService) error
	GetUserService(ctx context.Context, userservice model.UserService) (model.UserService, error)
	DeleteUserService(ctx context.Context, userservice model.UserService) error
	GetServiceUsers(ctx context.Context, userservice model.UserService) ([]model.UserService, error)
	GetUserServices(ctx context.Context, userservice model.UserService) ([]model.UserService, error)
	List(ctx context.Context) ([]model.UserService, error)
	Update(ctx context.Context, newUserService model.UserService) error
	DeleteUserServices(ctx context.Context, usrservices []model.UserService) error
	AddUserServices(ctx context.Context, usrservices []model.UserService) error
}

type UserService struct {
	UserServiceRepo repository.IUserServiceRepo
}

func (us *UserService) Add(ctx context.Context, userservice model.UserService) error {
	return us.UserServiceRepo.Add(ctx, userservice)
}

func (us *UserService) GetUserService(ctx context.Context, userservice model.UserService) (model.UserService, error) {
	return us.UserServiceRepo.GetUserService(ctx, userservice)
}

func (us *UserService) DeleteUserService(ctx context.Context, userservice model.UserService) error {
	return us.UserServiceRepo.DeleteUserService(ctx, userservice)
}

func (us *UserService) GetServiceUsers(ctx context.Context, userservice model.UserService) ([]model.UserService, error) {
	return us.UserServiceRepo.GetServiceUsers(ctx, userservice)
}

func (us *UserService) GetUserServices(ctx context.Context, userservice model.UserService) ([]model.UserService, error) {
	return us.UserServiceRepo.GetUserServices(ctx, userservice)
}

func (us *UserService) List(ctx context.Context) ([]model.UserService, error) {
	return us.UserServiceRepo.List(ctx)
}

func (us *UserService) Update(ctx context.Context, newUserService model.UserService) error {
	return us.UserServiceRepo.Update(ctx, newUserService)
}

func (us *UserService) DeleteUserServices(ctx context.Context, usrservices []model.UserService) error {
	return us.UserServiceRepo.DeleteUserServices(ctx, usrservices)
}

func (us *UserService) AddUserServices(ctx context.Context, usrservices []model.UserService) error {
	return us.UserServiceRepo.AddUserServices(ctx, usrservices)
}
