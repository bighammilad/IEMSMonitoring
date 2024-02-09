package usecase

import (
	"context"
	"monitoring/internal/model"
	"monitoring/internal/repository"
)

type IServicesUsecase interface {
	Add(ctx context.Context, service model.Service, userIds []int) error
	GetUserService(ctx context.Context, serviceName string, roleID, userId int) (service model.Service, err error)
	GetUserServices(ctx context.Context, roleID, userId int) (serviceRes []model.Service, err error)
	List(ctx context.Context) ([]model.Service, error)
	Update(ctx context.Context, service model.Service) error
	Delete(ctx context.Context, service model.Service) error
}

type ServicesUsecase struct {
	IServicesRepo repository.IServicesRepository
}

func (su *ServicesUsecase) Add(ctx context.Context, service model.Service, userIds []int) error {
	return su.IServicesRepo.Add(ctx, service, userIds)
}

func (su *ServicesUsecase) GetUserService(ctx context.Context, serviceName string, roleID, userId int) (service model.Service, err error) {

	service, err = su.IServicesRepo.GetUserService(ctx, serviceName, userId, roleID)
	return
}

func (su *ServicesUsecase) GetUserServices(ctx context.Context, roleID, userId int) (serviceRes []model.Service, err error) {

	serviceRes, err = su.IServicesRepo.GetUserServices(ctx, roleID, userId)
	return serviceRes, err
}

func (su *ServicesUsecase) List(ctx context.Context) ([]model.Service, error) {
	return su.IServicesRepo.List(ctx)
}

func (su *ServicesUsecase) Update(ctx context.Context, service model.Service) error {
	return su.IServicesRepo.Update(ctx, service)
}

func (su *ServicesUsecase) Delete(ctx context.Context, service model.Service) error {
	return su.IServicesRepo.Delete(ctx, service)
}
