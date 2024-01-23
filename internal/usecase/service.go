package usecase

import (
	"context"
	"errors"
	"monitoring/internal/model"
	"monitoring/internal/repository"
)

type IServicesUsecase interface {
	Add(ctx context.Context, service model.Service) error
	Get(ctx context.Context, service model.Service) (model.Service, error)
	List(ctx context.Context) ([]model.Service, error)
	Update(ctx context.Context, service model.Service) error
	Delete(ctx context.Context, service model.Service) error
}

type ServicesUsecase struct {
	ServicesRepo repository.ServicesRepository
}

func (su *ServicesUsecase) Add(ctx context.Context, service model.Service) error {
	return su.ServicesRepo.Add(ctx, service)
}

func (su *ServicesUsecase) Get(ctx context.Context, service model.Service) (serviceRes model.Service, err error) {

	// check id or name has been passed
	name := service.Name
	id := service.ID

	switch {
	case name != "":
		serviceRes, err = su.ServicesRepo.GetServiceByName(ctx, service)
	case id != 0:
		serviceRes, err = su.ServicesRepo.GetServiceById(ctx, service)
	default:
		return model.Service{}, errors.New("id or name must be passed")
	}

	return serviceRes, err
}

func (su *ServicesUsecase) List(ctx context.Context) ([]model.Service, error) {
	return su.ServicesRepo.List(ctx)
}

func (su *ServicesUsecase) Update(ctx context.Context, service model.Service) error {
	return su.ServicesRepo.Update(ctx, service)
}

func (su *ServicesUsecase) Delete(ctx context.Context, service model.Service) error {
	return su.ServicesRepo.Delete(ctx, service)
}
