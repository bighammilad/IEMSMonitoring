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

	// // Check if a user has access to a service
	// hasAccess := func(user model.User, service model.Service) bool {
	// 	if user.AccessLevel == 0 {
	// 		return true
	// 	}
	// 	for _, allowedUser := range service.AllowedUsers {
	// 		if allowedUser == user.Username {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }
	switch {
	case name != "":
		serviceRes, err = su.ServicesRepo.GetServiceByName(ctx, service)
	case id != 0:
		serviceRes, err = su.ServicesRepo.GetServiceById(ctx, service)
	default:
		return model.Service{}, errors.New("id or name must be passed")
	}

	// "sql: Scan error on column index 3, name \"header\": unsupported Scan, storing driver.Value type []uint8 into type *map[string]string"

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
